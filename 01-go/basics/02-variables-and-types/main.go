package main

import "fmt"

// 自定义类型是全新的类型，不会和底层类型隐式混用
type UserID int64

const (
	StatusPending = iota
	StatusPaid
	StatusCancelled
)

func main() {
	var n int
	var s string
	var ok bool
	var ptr *UserID
	fmt.Printf("zero: n=%d s=%q ok=%v ptr=%v\n", n, s, ok, ptr)

	explicit := int64(42)
	inferred := 42
	fmt.Printf("types: %T %T\n", explicit, inferred)

	// 无类型常量可以参与运算，但 int64 变量不行
	uid := UserID(1001)
	uid = uid + 1
	// uid = uid + explicit // 取消注释会编译失败：mismatched types UserID and int64
	fmt.Printf("uid=%d (%T)\n", uid, uid)

	const maxRetries = 3
	fmt.Printf("const=%d statuses=%d,%d,%d\n",
		maxRetries, StatusPending, StatusPaid, StatusCancelled)

	name := "百万"
	fmt.Printf("string: bytes=%d runes=%d\n", len(name), len([]rune(name)))
	for i, r := range name {
		fmt.Printf("  rune[%d]=%c U+%04X\n", i, r, r)
	}

	x := 10
	p := &x
	*p = 20
	fmt.Printf("pointer: x=%d *p=%d\n", x, *p)
}
