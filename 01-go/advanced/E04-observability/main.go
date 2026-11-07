package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// 结构化日志：一行一个 JSON，字段固定，方便采集
type LogEntry struct {
	Level   string `json:"level"`
	TraceID string `json:"trace_id"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Status  int    `json:"status"`
	Latency int64  `json:"latency_ms"`
	Msg     string `json:"msg,omitempty"`
}

// RED 指标：Rate / Errors / Duration
type Metrics struct {
	mu       sync.Mutex
	requests int
	errors   int
	totalMs  int64
	maxMs    int64
}

func (m *Metrics) Observe(status int, ms int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests++
	if status >= 500 {
		m.errors++
	}
	m.totalMs += ms
	if ms > m.maxMs {
		m.maxMs = ms
	}
}

func (m *Metrics) Snapshot() (int, int, float64, int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	avg := 0.0
	if m.requests > 0 {
		avg = float64(m.totalMs) / float64(m.requests)
	}
	return m.requests, m.errors, avg, m.maxMs
}

var traceSeq int
var traceMu sync.Mutex

func nextTraceID() string {
	traceMu.Lock()
	defer traceMu.Unlock()
	traceSeq++
	return fmt.Sprintf("trace-%04d", traceSeq)
}

// 模拟一次请求：trace id 在入口生成，沿途透传
func handle(m *Metrics, method, path string, status int, ms int64) {
	traceID := nextTraceID()
	m.Observe(status, ms)

	entry := LogEntry{
		Level:   "info",
		TraceID: traceID,
		Method:  method,
		Path:    path,
		Status:  status,
		Latency: ms,
	}
	if status >= 500 {
		entry.Level = "error"
		entry.Msg = "upstream failed"
	}
	line, _ := json.Marshal(entry)
	fmt.Println(string(line))
}

func main() {
	m := &Metrics{}

	// 固定的一组请求，保证输出可复现
	reqs := []struct {
		method, path string
		status       int
		ms           int64
	}{
		{"GET", "/healthz", 200, 1},
		{"GET", "/users", 200, 12},
		{"POST", "/users", 201, 18},
		{"GET", "/users/999", 404, 3},
		{"POST", "/users", 400, 2},
		{"GET", "/orders", 500, 240},
		{"GET", "/users", 200, 9},
		{"GET", "/users/1", 200, 7},
		{"GET", "/orders", 500, 260},
		{"DELETE", "/users/1", 204, 5},
	}

	for _, r := range reqs {
		handle(m, r.method, r.path, r.status, r.ms)
	}

	// 模拟一点真实耗时，让“Duration”不只是数字
	start := time.Now()
	time.Sleep(5 * time.Millisecond)
	elapsed := time.Since(start)

	requests, errors, avg, max := m.Snapshot()
	fmt.Println("---- RED 指标 ----")
	fmt.Printf("requests=%d errors=%d error_rate=%.2f%% avg_ms=%.1f max_ms=%d\n",
		requests, errors, float64(errors)/float64(requests)*100, avg, max)
	fmt.Printf("采集自身耗时=%v（仅演示计时）\n", elapsed.Round(time.Millisecond))
	fmt.Println("排查顺序：先看 error_rate 与 max_ms，再用 trace_id 串起同一次请求的全部日志")
}
