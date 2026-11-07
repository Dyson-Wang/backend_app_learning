package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Pool struct {
	slots chan struct{}
}

func NewPool(size int) *Pool {
	return &Pool{slots: make(chan struct{}, size)}
}

// Acquire 带超时：拿不到槽位就返回错误，而不是无限排队
func (p *Pool) Acquire(timeout time.Duration) bool {
	if timeout <= 0 {
		p.slots <- struct{}{}
		return true
	}
	select {
	case p.slots <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (p *Pool) Release() { <-p.slots }

func run(pool *Pool, n int, work time.Duration, acquireTimeout time.Duration) (ok, timeout int, elapsed time.Duration) {
	var wg sync.WaitGroup
	var okCount, timeoutCount int32
	start := time.Now()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !pool.Acquire(acquireTimeout) {
				atomic.AddInt32(&timeoutCount, 1)
				return
			}
			defer pool.Release()
			time.Sleep(work)
			atomic.AddInt32(&okCount, 1)
		}()
	}
	wg.Wait()
	return int(okCount), int(timeoutCount), time.Since(start)
}

func main() {
	const size = 3
	const requests = 10
	const work = 100 * time.Millisecond

	fmt.Println("== 1. 无超时：全部排队，慢但都成功 ==")
	ok, to, elapsed := run(NewPool(size), requests, work, 0)
	fmt.Printf("ok=%d timeout=%d 耗时=%v\n", ok, to, elapsed.Round(10*time.Millisecond))

	fmt.Println("== 2. 等待超过 50ms 就快速失败 ==")
	ok, to, elapsed = run(NewPool(size), requests, work, 50*time.Millisecond)
	fmt.Printf("ok=%d timeout=%d 耗时=%v\n", ok, to, elapsed.Round(10*time.Millisecond))

	fmt.Println("== 3. 结论 ==")
	fmt.Println("连接池必须配等待超时：没有超时，故障时请求会无限排队，最终拖垮整个服务")
	fmt.Println("池大小要按下游容量设，不是越大越好；过大会把数据库连接打满")
}
