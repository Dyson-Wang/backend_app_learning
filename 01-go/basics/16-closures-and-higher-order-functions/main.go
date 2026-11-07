package main

import "fmt"

// Adder 返回闭包：捕获参数 base，base 随闭包存活到函数返回之后（逃逸到堆）
func Adder(base int) func(int) int {
	return func(delta int) int {
		return base + delta
	}
}

// newCounter 演示捕获可变状态：被捕获的 n 独立于每次调用
func newCounter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

// Apply / Keep / Reduce 是典型高阶函数：接收函数作为参数
func Apply(xs []int, f func(int) int) []int {
	out := make([]int, len(xs))
	for i, x := range xs {
		out[i] = f(x)
	}
	return out
}

func Keep(xs []int, p func(int) bool) []int {
	out := []int{}
	for _, x := range xs {
		if p(x) {
			out = append(out, x)
		}
	}
	return out
}

func Reduce(xs []int, init int, f func(int, int) int) int {
	acc := init
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

// Compose 返回「返回函数」的高阶组合
func Compose(f, g func(int) int) func(int) int {
	return func(x int) int { return f(g(x)) }
}

// Greeter 演示方法值 / 方法表达式，即 Go 里的「self 绑定」
type Greeter struct{ prefix string }

func (g Greeter) Greet(name string) string {
	return g.prefix + name
}

func main() {
	// 1. 函数作为值
	op := func(a, b int) int { return a + b }
	fmt.Printf("函数作为值: op(3,4) = %d\n", op(3, 4))

	// 2. 返回函数 + 捕获变量的生命周期
	add10 := Adder(10)
	fmt.Printf("闭包: add10(5) = %d, add10(7) = %d（base=10 一直存活）\n", add10(5), add10(7))

	// 3. 捕获可变状态：不同闭包各自独立，同一个变量被多处闭包共享
	c1, c2 := newCounter(), newCounter()
	fmt.Printf("可变捕获: c1 独立 = %d,%d；c2 独立 = %d\n", c1(), c1(), c2())
	shared := 0
	bump := func() int { shared++; return shared }
	fmt.Printf("共享捕获: bump(),bump() = %d,%d（同一变量被共享）\n", bump(), bump())

	// 4. 循环变量捕获陷阱：显式遮蔽才每个闭包拿到自己的值
	fs := make([]func() int, 0, 3)
	for i := 0; i < 3; i++ {
		i := i // Go 1.22 起循环变量默认每次迭代新建；1.21 需手写这行
		fs = append(fs, func() int { return i })
	}
	fmt.Printf("循环捕获(遮蔽): %d,%d,%d\n", fs[0](), fs[1](), fs[2]())

	// 5. 常见高阶函数 + 返回函数的组合
	nums := []int{1, 2, 3, 4}
	fmt.Printf("Apply(*2) = %v\n", Apply(nums, func(n int) int { return n * 2 }))
	fmt.Printf("Keep(偶数) = %v\n", Keep(nums, func(n int) bool { return n%2 == 0 }))
	fmt.Printf("Reduce(求和) = %d\n", Reduce(nums, 0, func(a, b int) int { return a + b }))
	twiceThenSquare := Compose(func(n int) int { return n * n }, func(n int) int { return n * 2 })
	fmt.Printf("Compose(先乘2再平方)(3) = %d\n", twiceThenSquare(3))

	// 6. 回调与 self 绑定：方法值已绑定接收者，方法表达式需显式传接收者
	g := Greeter{prefix: "你好, "}
	bound := g.Greet // 方法值：绑定了 g
	fmt.Printf("方法值(已绑定 self): %s\n", bound("Go"))
	unbound := Greeter.Greet // 方法表达式：第一个参数是接收者
	fmt.Printf("方法表达式(显式传 self): %s\n", unbound(g, "闭包"))
}
