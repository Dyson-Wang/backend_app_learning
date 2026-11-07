package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

// merge 扇入：把多个只读 channel 汇成一个
func merge(chans ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	for _, c := range chans {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	// 1. 有界队列 + 生产消费 + 关闭信号传播
	jobs := make(chan int, 3) // 容量 3，满则发送阻塞，天然背压
	results := make(chan int, 5)
	var workers sync.WaitGroup
	for w := 0; w < 2; w++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for j := range jobs { // close(jobs) 后 range 结束，goroutine 退出
				results <- j * 2
			}
		}()
	}
	for j := 1; j <= 5; j++ {
		jobs <- j // 第 4 个会阻塞直到 worker 取走，演示有界队列背压
	}
	close(jobs) // 关闭信号沿 channel 传播给所有 worker
	workers.Wait()
	close(results)
	var doubled []int
	for r := range results {
		doubled = append(doubled, r)
	}
	sort.Ints(doubled) // 排序让并发结果稳定可复现
	fmt.Printf("生产消费/有界队列/关闭传播: %v\n", doubled)

	// 2. 扇出扇入：3 个 worker 竞争同一源，结果各自 channel，再 merge
	src := make(chan int)
	outs := make([]<-chan int, 3)
	for i := 0; i < 3; i++ {
		ch := make(chan int)
		outs[i] = ch
		go func() {
			defer close(ch)
			for v := range src { // 扇出：多个 worker 竞争消费同一输入
				ch <- v * 10
			}
		}()
	}
	go func() {
		for v := 1; v <= 6; v++ {
			src <- v
		}
		close(src)
	}()
	var fanned []int
	for v := range merge(outs...) { // 扇入：合并所有 worker 输出
		fanned = append(fanned, v)
	}
	sort.Ints(fanned)
	fmt.Printf("扇出扇入(3 worker, 输入1..6各×10): %v\n", fanned)

	// 3. 等待多个任务 + close 广播关闭信号
	done := make(chan struct{})
	var tasks sync.WaitGroup
	for i := 0; i < 3; i++ {
		tasks.Add(1)
		go func() {
			defer tasks.Done()
			<-done // 收到广播信号才结束
		}()
	}
	close(done) // 一次 close 即向所有等待者广播
	tasks.Wait()
	fmt.Println("关闭信号: close(done) 一次广播，3 个任务全部被唤醒并结束")

	// 4. 共享状态保护：atomic 与 Mutex 两种写法
	var atomicSum int64
	var mu sync.Mutex
	mutexSum := 0
	var sw sync.WaitGroup
	for i := 1; i <= 100; i++ {
		sw.Add(1)
		go func(n int) {
			defer sw.Done()
			atomic.AddInt64(&atomicSum, int64(n))
			mu.Lock()
			mutexSum += n
			mu.Unlock()
		}(i)
	}
	sw.Wait()
	fmt.Printf("共享状态: atomic = %d, mutex = %d（都应为 5050）\n", atomicSum, mutexSum)
}
