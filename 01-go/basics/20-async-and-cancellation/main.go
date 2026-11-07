package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "requestID"

// worker 是一个「异步任务」：goroutine + channel，靠 select 监听 ctx.Done() 实现可取消
func worker(ctx context.Context, id int, out chan<- string) {
	for {
		select {
		case <-ctx.Done():
			out <- fmt.Sprintf("worker%d 取消(%v) → 已清理", id, ctx.Err())
			return
		case <-time.After(5 * time.Millisecond):
			// 模拟持续工作；每次循环重新检查 ctx
		}
	}
}

func main() {
	// 1. 超时：WithTimeout 到点自动取消，worker 经 select 收到并清理
	tctx, tcancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer tcancel()
	ch := make(chan string, 1)
	go worker(tctx, 0, ch)
	fmt.Printf("超时演示: %s\n", <-ch)
	fmt.Printf("超时后 ctx.Err() = %v\n", tctx.Err())

	// 2. WithCancel 手动取消 + 取消沿父子链传播 + 取消后清理
	parent, pcancel := context.WithCancel(context.Background())
	child, ccancel := context.WithCancel(parent)
	defer ccancel()
	var wg sync.WaitGroup
	results := make(chan string, 3)
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-child.Done()
			results <- fmt.Sprintf("任务%d 收到父级取消(%v)", id, child.Err())
		}(i)
	}
	pcancel() // 一次取消，parent 与 child 同时收到
	wg.Wait()
	close(results)
	for m := range results {
		fmt.Printf("取消传播: %s\n", m)
	}

	// 3. 并发等待多个异步任务：启动 N 个 goroutine，select 收齐全部
	const n = 3
	waitCh := make(chan int, n)
	for i := 0; i < n; i++ {
		go func(v int) { waitCh <- v * v }(i)
	}
	sum := 0
	for i := 0; i < n; i++ {
		select {
		case v := <-waitCh:
			sum += v
		case <-time.After(500 * time.Millisecond): // 兜底，本例不会触发
			fmt.Println("等待超时")
		}
	}
	fmt.Printf("并发等待: 收齐 %d 个任务, 平方和 = %d\n", n, sum)

	// 4. 上下文传值：请求作用域的数据随 ctx 一路传递
	vctx := context.WithValue(context.Background(), requestIDKey, "req-42")
	done := make(chan string, 1)
	go func(ctx context.Context) {
		done <- fmt.Sprintf("取到 %v = %v", requestIDKey, ctx.Value(requestIDKey))
	}(vctx)
	fmt.Printf("上下文传值: %s\n", <-done)
	fmt.Printf("上下文传值: 取不存在的键 = %v\n", vctx.Value(ctxKey("missing")))
}
