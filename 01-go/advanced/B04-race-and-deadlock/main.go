package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 无锁「读-改-写」：用 WaitGroup 强制两个 goroutine 都先读完再写，
// 保证交错必然发生，输出确定。
var racy int
var readBarrier sync.WaitGroup

func racyInc() {
	v := racy
	readBarrier.Done()
	readBarrier.Wait()
	racy = v + 1
}

func main() {
	fmt.Println("== 1. 无锁读改写：两个 goroutine 各加一次 ==")
	readBarrier.Add(2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); racyInc() }()
	go func() { defer wg.Done(); racyInc() }()
	wg.Wait()
	fmt.Printf("期望 2，实际 %d（丢更新）\n", racy)

	fmt.Println("== 2. 加锁：4 goroutine × 1000 次 ==")
	var mu sync.Mutex
	counter := 0
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("期望 4000，实际 %d\n", counter)

	fmt.Println("== 3. 原子操作：同样 4000 次 ==")
	var atomicCounter int64
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&atomicCounter, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("期望 4000，实际 %d\n", atomicCounter)

	fmt.Println("== 4. 反序加锁导致死锁 ==")
	var a, b sync.Mutex
	done := make(chan string, 2)
	go func() {
		a.Lock()
		time.Sleep(10 * time.Millisecond)
		b.Lock()
		b.Unlock()
		a.Unlock()
		done <- "goroutine-1 完成"
	}()
	go func() {
		b.Lock()
		time.Sleep(10 * time.Millisecond)
		a.Lock()
		a.Unlock()
		b.Unlock()
		done <- "goroutine-2 完成"
	}()
	select {
	case msg := <-done:
		fmt.Println("未死锁：", msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("检测到死锁：100ms 内无人完成（各持一把锁等对方）")
	}

	fmt.Println("== 5. 固定加锁顺序 ==")
	var c, d sync.Mutex
	done2 := make(chan string, 2)
	lockBoth := func(name string) {
		c.Lock()
		time.Sleep(5 * time.Millisecond)
		d.Lock()
		d.Unlock()
		c.Unlock()
		done2 <- name + " 完成"
	}
	go lockBoth("goroutine-1")
	go lockBoth("goroutine-2")
	for i := 0; i < 2; i++ {
		select {
		case msg := <-done2:
			fmt.Println("  ", msg)
		case <-time.After(200 * time.Millisecond):
			fmt.Println("  意外超时")
		}
	}
	fmt.Println("结论：死锁靠「统一加锁顺序」或「带超时的 TryLock」避免，不能靠运气")
}
