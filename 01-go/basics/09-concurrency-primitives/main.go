package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// 1. goroutine + WaitGroup
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("worker %d done\n", id)
		}(i)
	}
	wg.Wait()

	// 2. channel：关闭后仍可 range 出缓冲里剩余的值
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	for v := range ch {
		fmt.Println("buffered channel:", v)
	}

	// 3. 生产者 / 消费者
	jobs := make(chan int)
	go func() {
		defer close(jobs)
		for i := 0; i < 3; i++ {
			jobs <- i
		}
	}()
	for j := range jobs {
		fmt.Println("job:", j)
	}

	// 4. select + 超时（没有超时的等待是生产事故之源）
	slow := make(chan string, 1)
	go func() {
		time.Sleep(30 * time.Millisecond)
		slow <- "result"
	}()
	select {
	case res := <-slow:
		fmt.Println("got:", res)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timeout")
	}

	// 5. 共享状态必须加锁
	var mu sync.Mutex
	counter := 0
	var cwg sync.WaitGroup
	for i := 0; i < 100; i++ {
		cwg.Add(1)
		go func() {
			defer cwg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	cwg.Wait()
	fmt.Println("counter with mutex:", counter)

	// 6. sync.Once：只初始化一次
	var once sync.Once
	for i := 0; i < 3; i++ {
		once.Do(func() { fmt.Println("sync.Once: initialized exactly once") })
	}

	// 7. context 取消传播
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done:", ctx.Err())
				return
			case <-time.After(20 * time.Millisecond):
				fmt.Println("tick")
			}
		}
	}()
	<-done
}
