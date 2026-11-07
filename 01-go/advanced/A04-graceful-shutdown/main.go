package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

func runSelfTest() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	addr := ln.Addr().String()

	var completed int32
	mux := http.NewServeMux()
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond) // 模拟在途的慢请求
		atomic.AddInt32(&completed, 1)
		fmt.Fprintln(w, "done")
	})
	srv := &http.Server{Handler: mux}

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()

	type result struct {
		status int
		err    error
	}
	resCh := make(chan result, 1)
	client := &http.Client{Timeout: 3 * time.Second}
	go func() {
		resp, err := client.Get("http://" + addr + "/slow")
		if err != nil {
			resCh <- result{0, err}
			return
		}
		defer resp.Body.Close()
		io.Copy(io.Discard, resp.Body)
		resCh <- result{resp.StatusCode, nil}
	}()

	// 等请求真正进入 handler 再关闭，否则测的是「还没开始」而不是「正在处理」
	time.Sleep(50 * time.Millisecond)

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	fmt.Printf("Shutdown returned after %s (it waited for the in-flight request)\n",
		time.Since(start).Round(10*time.Millisecond))

	res := <-resCh
	fmt.Printf("in-flight request -> status=%d err=%v handler_completed=%d\n",
		res.status, res.err, atomic.LoadInt32(&completed))
	if res.err != nil || res.status != http.StatusOK || atomic.LoadInt32(&completed) != 1 {
		log.Fatal("FAIL: in-flight request was dropped")
	}
	if err := <-serveErr; err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve returned unexpected error: %v", err)
	}

	// Shutdown 后不再接受新连接
	client.CloseIdleConnections()
	resp, err := client.Get("http://" + addr + "/slow")
	if err == nil {
		resp.Body.Close()
		log.Fatal("FAIL: server still accepting new connections after Shutdown")
	}
	fmt.Println("new connections rejected after Shutdown")
	fmt.Println("selftest PASS")
}

func main() {
	runSelfTest()
}
