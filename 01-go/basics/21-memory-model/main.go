package main

import (
	"fmt"
	"runtime"
)

type Point struct{ X, Y int }

// escapesToHeap 返回局部变量的地址：x 必须「逃逸」到堆，否则返回后指针悬空
func escapesToHeap() *int {
	x := 42
	return &x
}

// byValue 演示值传递的拷贝：形参是副本
func byValue(p Point) Point {
	p.X = 999
	return p
}

func main() {
	// 1. 值 vs 引用：结构体赋值是深拷贝，指针/引用共享同一对象
	p := Point{1, 2}
	q := p
	q.X = 99
	fmt.Printf("值拷贝: p = %+v, q = %+v（改 q 不动 p）\n", p, q)
	ref := &p
	ref.X = 42
	fmt.Printf("指针共享: p = %+v（经指针改到了原对象）\n", p)
	fmt.Printf("值传递: byValue 返回 = %+v, 原 p = %+v（未被修改）\n", byValue(p), p)

	// 2. slice / map 是「引用语义」的描述符，底层数据共享
	s := []int{1, 2, 3}
	t := s
	t[0] = 9
	fmt.Printf("slice 共享底层数组: s = %v（t 改了 s 也变）\n", s)
	m := map[string]int{"a": 1}
	m2 := m
	m2["a"] = 7
	fmt.Printf("map 引用语义: m[\"a\"] = %d\n", m["a"])

	// 3. 逃逸：局部变量的地址被返回，编译器把它分配到堆
	ep := escapesToHeap()
	*ep = 100
	fmt.Printf("逃逸到堆: 返回局部变量指针后仍可读写 = %d\n", *ep)

	// 4. 常见内存陷阱
	// (a) append 扩容后指向新数组，副本不再共享
	base := []int{1, 2, 3} // len=cap=3
	grown := append(base, 4)
	grown[0] = 9
	fmt.Printf("陷阱-append 扩容: base = %v, grown = %v（已分离）\n", base, grown)
	// (b) 子切片仍持有整个大数组，导致大对象无法回收
	big := make([]byte, 1024)
	head := big[:3]
	fmt.Printf("陷阱-子切片: head len=%d cap=%d（仍引用 1024 字节底层数组）\n", len(head), cap(head))
	// (c) 想切断引用要 copy 到新切片
	detached := append([]byte(nil), head...)
	fmt.Printf("陷阱-修正: copy 后 len=%d, 已脱离大数组 = %t\n",
		len(detached), cap(detached) < cap(head))

	// 5. GC 与对象生命周期
	runtime.GC()
	fmt.Println("GC: runtime.GC() 主动建议一次回收；Go 是并发三色标记 GC，无所有权/引用计数/手动 free")
	fmt.Println("生命周期: 对象是否可达由 GC 判定，逃逸分析决定栈/堆；用 go build -gcflags=-m 看逃逸报告")
}
