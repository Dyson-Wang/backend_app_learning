package main

import (
	"fmt"
	"sort"
)

// 泛型函数：Go 1.18+ 支持，标准库里没有内置 map/filter 就自己写
func MapSlice[T any, R any](in []T, fn func(T) R) []R {
	out := make([]R, 0, len(in))
	for _, v := range in {
		out = append(out, fn(v))
	}
	return out
}

func main() {
	var arr [3]int
	arr[0] = 1
	fmt.Printf("array: %v len=%d\n", arr, len(arr))

	s := []int{1, 2, 3}
	s = append(s, 4, 5)
	fmt.Printf("slice: %v len=%d cap=%d\n", s, len(s), cap(s))

	// 共享底层数组陷阱：a 的 cap 还够，append 会覆写 base[2]
	base := []int{1, 2, 3, 4}
	a := base[:2]
	b := append(a, 99)
	fmt.Printf("shared backing: base=%v a=%v b=%v\n", base, a, b)
	fmt.Println("  -> base[2] 被覆写了，这是最经典的 slice 坑")

	safe := make([]int, 2)
	copy(safe, base[:2])
	safe = append(safe, 77)
	fmt.Printf("with copy:      base=%v safe=%v\n", base, safe)

	m := map[string]int{"go": 1, "rust": 2}
	m["python"] = 3
	if v, ok := m["java"]; !ok {
		fmt.Println("comma-ok: java missing, zero value =", v)
	}
	delete(m, "go")
	fmt.Printf("map: %v len=%d\n", m, len(m))

	words := []string{"banana", "apple", "cherry"}
	sort.Slice(words, func(i, j int) bool { return words[i] < words[j] })
	fmt.Println("sorted:", words)

	doubled := MapSlice([]int{1, 2, 3}, func(n int) int { return n * 2 })
	fmt.Println("generic map:", doubled)

	for i, v := range words {
		fmt.Printf("  words[%d]=%s\n", i, v)
	}
}
