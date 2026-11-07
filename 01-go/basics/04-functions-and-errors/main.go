package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidInput = errors.New("invalid input")

// 自定义错误类型：实现 error 接口即可
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("divide(%d, 0): %w", a, ErrInvalidInput)
	}
	return a / b, nil
}

// 命名返回值：裸 return 可用，但别用过头
func double(n int) (result int, err error) {
	if n < 0 {
		err = fmt.Errorf("negative input: %d", n)
		return
	}
	result = n * 2
	return
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// recover 只在 defer 里有效，用来把 panic 转成错误边界（库代码不要滥用）
func safeCall(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	fn()
}

func main() {
	if v, err := divide(10, 2); err == nil {
		fmt.Println("10/2 =", v)
	}

	_, err := divide(1, 0)
	fmt.Println("err:", err)
	fmt.Println("errors.Is(ErrInvalidInput):", errors.Is(err, ErrInvalidInput))

	_, fileErr := os.Open("/definitely/not/exist")
	var pathErr *os.PathError
	if errors.As(fileErr, &pathErr) {
		fmt.Printf("errors.As -> *os.PathError op=%s path=%s\n", pathErr.Op, pathErr.Path)
	}

	var valErr *ValidationError
	wrapped := fmt.Errorf("create user: %w", &ValidationError{Field: "email", Reason: "empty"})
	if errors.As(wrapped, &valErr) {
		fmt.Printf("custom error -> field=%s reason=%s\n", valErr.Field, valErr.Reason)
	}

	r, e := double(-1)
	fmt.Printf("named return: result=%d err=%v\n", r, e)
	fmt.Println("variadic sum:", sum(1, 2, 3, 4))

	adder := func(base int) func(int) int {
		return func(x int) int { return base + x }
	}
	fmt.Println("closure:", adder(10)(5))

	safeCall(func() { panic("something went wrong") })
	fmt.Println("still running after recover")
}
