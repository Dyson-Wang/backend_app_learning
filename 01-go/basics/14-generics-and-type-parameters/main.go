package main

import "fmt"

// Number 是自定义 constraint：用 | 列出允许的类型，~ 表示底层类型匹配
type Number interface {
	~int | ~int64 | ~float64
}

// Sum 泛型函数：T 受 Number 约束，才能用 + 和零值
func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

// Map 泛型函数：两个类型参数，返回新的切片
func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs))
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

// IndexOf 用内建 comparable 约束：只有可比较类型才能用 ==
func IndexOf[T comparable](xs []T, target T) int {
	for i, x := range xs {
		if x == target {
			return i
		}
	}
	return -1
}

// Stack 泛型类型：结构体自身带类型参数
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// Set 泛型集合：底层用 map[T]struct{}，要求 T 可比较
type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{m: make(map[T]struct{})}
}
func (s *Set[T]) Add(v T)      { s.m[v] = struct{}{} }
func (s *Set[T]) Has(v T) bool { _, ok := s.m[v]; return ok }
func (s *Set[T]) Len() int     { return len(s.m) }

// Pair 多类型参数的泛型类型
type Pair[A, B any] struct {
	First  A
	Second B
}

func main() {
	// 1. 泛型函数 + 类型推断：编译器从实参推出 T
	fmt.Printf("Sum(int)   = %d\n", Sum([]int{1, 2, 3, 4}))
	fmt.Printf("Sum(float) = %v\n", Sum([]float64{1.5, 2.5}))

	// 2. 显式类型实参（推断不出时才需要）
	fmt.Printf("Sum 显式类型实参 = %d\n", Sum[int64]([]int64{10, 20}))

	// 3. Map：T=int, U=string，由参数与闭包推断
	strs := Map([]int{1, 2, 3}, func(n int) string {
		return fmt.Sprintf("#%d", n)
	})
	fmt.Printf("Map int->string = %v\n", strs)

	// 4. comparable 约束
	fmt.Printf("IndexOf: %d, %d\n",
		IndexOf([]string{"a", "b", "c"}, "b"), IndexOf([]string{"a", "b"}, "z"))

	// 5. 泛型类型 Stack：同一份代码，不同类型实例
	var is Stack[int]
	is.Push(10)
	is.Push(20)
	top, ok := is.Pop()
	fmt.Printf("Stack[int]: Pop = %d, ok = %t\n", top, ok)
	var ss Stack[string]
	ss.Push("go")
	sv, _ := ss.Pop()
	fmt.Printf("Stack[string]: Pop = %q\n", sv)

	// 6. 泛型集合 Set + 泛型类型 Pair
	set := NewSet[string]()
	set.Add("go")
	set.Add("rust")
	set.Add("go")
	fmt.Printf("Set[string]: Len = %d, Has(\"go\") = %t, Has(\"java\") = %t\n",
		set.Len(), set.Has("go"), set.Has("java"))
	p := Pair[int, string]{First: 7, Second: "七"}
	fmt.Printf("Pair[int,string] = {%d %q}\n", p.First, p.Second)
}
