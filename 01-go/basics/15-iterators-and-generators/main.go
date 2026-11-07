package main

import "fmt"

// Seq 是「push 迭代器」：一个接收 yield 回调的函数。
// 调用时把值推给 yield；yield 返回 false 表示消费者要提前停止。
// 这正是 Go 1.23 range-over-func / iter.Seq 的协议，但 1.21 只能手动调用。
type Seq[T any] func(yield func(T) bool)

// From 把切片包装成 Seq（可反复消费）
func From[T any](xs []T) Seq[T] {
	return func(yield func(T) bool) {
		for _, x := range xs {
			if !yield(x) {
				return
			}
		}
	}
}

// Nats 是无限自然数序列：靠 yield 返回 false 才能停下
func Nats() Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// 惰性组合子：都只返回新的 Seq，构造阶段不消费任何值
func Map[T, U any](s Seq[T], f func(T) U) Seq[U] {
	return func(yield func(U) bool) {
		s(func(x T) bool { return yield(f(x)) })
	}
}

func Filter[T any](s Seq[T], p func(T) bool) Seq[T] {
	return func(yield func(T) bool) {
		s(func(x T) bool {
			if p(x) {
				return yield(x)
			}
			return true // 跳过当前值，但继续迭代
		})
	}
}

func Take[T any](s Seq[T], n int) Seq[T] {
	return func(yield func(T) bool) {
		i := 0
		s(func(x T) bool {
			if i >= n {
				return false
			}
			i++
			return yield(x)
		})
	}
}

func Chain[T any](a, b Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		more := true
		a(func(x T) bool {
			more = yield(x)
			return more
		})
		if more {
			b(yield)
		}
	}
}

// Collect 消费一个 Seq，把值收进切片
func Collect[T any](s Seq[T]) []T {
	var out []T
	s(func(x T) bool {
		out = append(out, x)
		return true
	})
	return out
}

// ChanSeq 由 channel 支撑：channel 耗尽后无法重放，演示「一次性消费陷阱」
func ChanSeq[T any](ch <-chan T) Seq[T] {
	return func(yield func(T) bool) {
		for x := range ch {
			if !yield(x) {
				return
			}
		}
	}
}

func main() {
	// 1. 手动用 yield 驱动迭代器
	fmt.Print("手动迭代: ")
	From([]int{10, 20, 30})(func(x int) bool {
		fmt.Printf("%d ", x)
		return true
	})
	fmt.Println()

	// 2. 惰性求值：Map/Filter/Take 只构造管道，Collect 才真正计算
	pipeline := Take(
		Filter(
			Map(Nats(), func(n int) int { return n * n }),
			func(n int) bool { return n%2 == 0 },
		),
		5,
	)
	fmt.Println("惰性: 上面只构造了 Seq，下面 Collect 才触发计算")
	fmt.Printf("管道(平方中偶数前 5 个): %v\n", Collect(pipeline))

	// 3. chain 拼接两个序列
	merged := Chain(From([]string{"a", "b"}), From([]string{"c", "d"}))
	fmt.Printf("Chain: %v\n", Collect(merged))

	// 4. 无限序列取前 N 项：不会真的无限跑
	fmt.Printf("Nats 前 6 项: %v\n", Collect(Take(Nats(), 6)))

	// 5. 一次性消费陷阱：channel 支撑的 Seq 只能读一次
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	one := ChanSeq(ch)
	fmt.Printf("一次性(第 1 次): %v\n", Collect(one))
	fmt.Printf("一次性(第 2 次): %v（channel 已空）\n", Collect(one))

	// 6. 提前停止：yield 返回 false 即中断
	fmt.Print("提前停止(到 3 即停): ")
	Nats()(func(n int) bool {
		fmt.Printf("%d ", n)
		return n < 3
	})
	fmt.Println()
}
