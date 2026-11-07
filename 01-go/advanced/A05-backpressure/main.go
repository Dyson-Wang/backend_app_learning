package main

import (
	"fmt"
)

// 有界队列 + 满则拒绝，是最朴素的背压。
// 消费者先被 gate 挡住，保证提交阶段队列一定处于「满」的状态，输出可复现。
func main() {
	const capacity = 5
	const total = 20

	q := make(chan int, capacity)
	gate := make(chan struct{})
	done := make(chan int, 1)

	go func() {
		<-gate
		processed := 0
		for range q {
			processed++
		}
		done <- processed
	}()

	accepted, rejected := 0, 0
	for i := 0; i < total; i++ {
		select {
		case q <- i:
			accepted++
		default:
			rejected++ // 队列满：立即拒绝，而不是无限堆积
		}
	}

	close(gate)
	close(q)
	processed := <-done

	fmt.Printf("容量=%d，提交=%d\n", capacity, total)
	fmt.Printf("accepted=%d rejected=%d processed=%d\n", accepted, rejected, processed)
	fmt.Println("结论：拒绝发生在入口，内存和延迟都是常数；无界队列只是把 OOM 推迟到更晚")

	// 对比：无界收集（错误示范）
	unbounded := make([]int, 0, total)
	for i := 0; i < total; i++ {
		unbounded = append(unbounded, i)
	}
	fmt.Printf("无界队列长度=%d（全部接收，没有上限）\n", len(unbounded))
}
