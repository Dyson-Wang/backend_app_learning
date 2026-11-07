package main

import "fmt"

func main() {
	// if 可以带初始化语句，变量作用域限于 if/else
	if v := 7; v%2 == 0 {
		fmt.Println("even")
	} else {
		fmt.Println("odd:", v)
	}

	switch n := 3; n {
	case 1, 2:
		fmt.Println("small")
	case 3:
		fmt.Println("three")
		fallthrough
	case 4:
		fmt.Println("case 4 (also reached via fallthrough)")
	default:
		fmt.Println("other")
	}

	age := 20
	switch {
	case age < 18:
		fmt.Println("minor")
	case age < 65:
		fmt.Println("adult")
	default:
		fmt.Println("senior")
	}

	// Go 只有 for 一种循环关键字
	for i := 0; i < 3; i++ {
		fmt.Print(i, " ")
	}
	fmt.Println()

	i := 0
	for i < 3 {
		i++
	}
	fmt.Println("while-style result:", i)

	nums := []int{10, 20, 30}
	for idx, v := range nums {
		fmt.Printf("nums[%d]=%d\n", idx, v)
	}

	m := map[string]int{"a": 1, "b": 2}
	for k, v := range m {
		fmt.Printf("m[%s]=%d\n", k, v)
	}

	defer fmt.Println("defer: 3rd registered, last executed")
	defer fmt.Println("defer: 2nd registered")
	defer fmt.Println("defer: 1st registered")

	// 标签跳出多层循环
outer:
	for a := 0; a < 3; a++ {
		for b := 0; b < 3; b++ {
			if a*b == 4 {
				fmt.Println("break outer at", a, b)
				break outer
			}
		}
	}
}
