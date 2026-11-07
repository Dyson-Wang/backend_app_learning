package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Task 是要返回的 JSON 资源
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// ErrCode 用 iota 做错误码枚举（13 课的写法用回服务）
type ErrCode int

const (
	ErrBadRequest ErrCode = iota + 1
	ErrNotFound
	ErrConflict
)

func (e ErrCode) String() string {
	switch e {
	case ErrBadRequest:
		return "BAD_REQUEST"
	case ErrNotFound:
		return "NOT_FOUND"
	case ErrConflict:
		return "CONFLICT"
	default:
		return "UNKNOWN"
	}
}

func (e ErrCode) HTTPStatus() int {
	switch e {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeJSON 是泛型辅助：任何载荷都用同一套编码流程（14 课）
func writeJSON[T any](w http.ResponseWriter, status int, body T) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func fail(w http.ResponseWriter, code ErrCode, msg string) {
	writeJSON(w, code.HTTPStatus(), errorBody{Code: code.String(), Message: msg})
}

// Store 泛型内存存储，用 RWMutex 保证并发安全（14 + 19 课）
type Store[T any] struct {
	mu    sync.RWMutex
	items map[int]T
	next  int
	setID func(*T, int) // 闭包：把自增 id 写回具体类型（16 课）
}

func NewStore[T any](setID func(*T, int)) *Store[T] {
	return &Store[T]{items: make(map[int]T), next: 1, setID: setID}
}

func (s *Store[T]) Add(v T) T {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setID(&v, s.next)
	s.items[s.next] = v
	s.next++
	return v
}

func (s *Store[T]) Get(id int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.items[id]
	return v, ok
}

func (s *Store[T]) List() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]T, 0, len(s.items))
	for id := 1; id < s.next; id++ {
		if v, ok := s.items[id]; ok {
			out = append(out, v)
		}
	}
	return out
}

// statusRecorder 记录真实写出的状态码，供中间件记日志
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func tasksHandler(store *Store[Task]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, store.List())
		case http.MethodPost:
			var in struct {
				Title string `json:"title"`
			}
			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				fail(w, ErrBadRequest, "JSON 解析失败: "+err.Error())
				return
			}
			in.Title = strings.TrimSpace(in.Title)
			if in.Title == "" {
				fail(w, ErrBadRequest, "title 不能为空")
				return
			}
			created := store.Add(Task{Title: in.Title})
			writeJSON(w, http.StatusCreated, created)
		default:
			fail(w, ErrBadRequest, "不支持的方法 "+r.Method)
		}
	}
}

func taskByIDHandler(store *Store[Task]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/tasks/"))
		if err != nil {
			fail(w, ErrBadRequest, "id 必须是整数")
			return
		}
		task, ok := store.Get(id)
		if !ok {
			fail(w, ErrNotFound, "任务不存在")
			return
		}
		writeJSON(w, http.StatusOK, task)
	}
}

// do 发一个请求并返回状态码与响应体（自请求自关闭）
func do(ctx context.Context, client *http.Client, method, url, body string) (int, string) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return 0, err.Error()
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, strings.TrimSpace(string(data))
}

func main() {
	built := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	fmt.Printf("服务: 构建时间 %s（固定值，只演示 time 格式化，不打印当前时间）\n", built)

	store := NewStore[Task](func(t *Task, id int) { t.ID = id })
	store.Add(Task{Title: "读规格", Done: true})
	store.Add(Task{Title: "搭骨架", Done: false})

	var reqID atomic.Int64
	logCh := make(chan string, 32)
	// 日志中间件：闭包包装 handler（16 课），用 atomic 计数（19 课）
	withLogging := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := reqID.Add(1)
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			logCh <- fmt.Sprintf("#%d %s %s -> %d", id, r.Method, r.URL.Path, rec.status)
		})
	}

	mux := http.NewServeMux() // 路由（Go 1.21 不支持 "GET /x" 模式，方法在 handler 内判）
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/tasks", tasksHandler(store))
	mux.HandleFunc("/tasks/", taskByIDHandler(store))
	srv := httptest.NewServer(withLogging(mux))
	defer func() {
		srv.Close()
		fmt.Println("服务: 所有请求完成，httptest 服务已自行关闭（3 秒内退出）")
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // 20 课：取消/超时
	defer cancel()
	client := srv.Client()
	base := srv.URL // 随机回环端口，不打印具体地址以保持输出可复现

	requests := []struct {
		label, method, path, body string
	}{
		{"健康检查", http.MethodGet, "/health", ""},
		{"任务列表", http.MethodGet, "/tasks", ""},
		{"创建(合法)", http.MethodPost, "/tasks", `{"title":"写代码"}`},
		{"创建(空标题)", http.MethodPost, "/tasks", `{"title":"  "}`},
		{"查询不存在", http.MethodGet, "/tasks/999", ""},
		{"查询 id=1", http.MethodGet, "/tasks/1", ""},
	}
	for _, rq := range requests {
		status, body := do(ctx, client, rq.method, base+rq.path, rq.body)
		fmt.Printf("请求[%s] %s %s -> %d %s\n", rq.label, rq.method, rq.path, status, body)
		fmt.Printf("  日志: %s\n", <-logCh) // 等本请求日志，保证顺序稳定
	}
	fmt.Println("说明: 泛型 writeJSON/Store、iota 错误码枚举、闭包中间件、atomic/Mutex 并发保护、context+time 取消一并使用")
}
