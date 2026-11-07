package main

import (
	"errors"
	"fmt"
)

var failures int

func check(name string, ok bool, detail string) {
	if ok {
		fmt.Printf("PASS %s\n", name)
		return
	}
	failures++
	fmt.Printf("FAIL %s: %s\n", name, detail)
}

func add(a, b int) int { return a + b }

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func parseAge(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	return 18, nil
}

func main() {
	fmt.Println("== 1. 单元测试 ==")
	check("add(1,2)=3", add(1, 2) == 3, "期望 3")
	_, err := divide(1, 0)
	check("divide(1,0) 返回错误", err != nil, "期望错误，实际 nil")
	age, err := parseAge("abc")
	check("parseAge(\"abc\") 返回错误", err != nil, fmt.Sprintf("期望错误，实际 %d, %v", age, err))

	fmt.Println("== 2. 汇总 ==")
	fmt.Printf("3 个用例：%d 通过，%d 失败\n", 3-failures, failures)

	fmt.Println("== 3. 结论 ==")
	fmt.Println("最后一条是故意失败的用例，用来看失败输出长什么样")
	fmt.Println("单测不碰网络/DB；集成测试才起真实依赖，两者分开跑")
	fmt.Println("容器化：多阶段构建、非 root 用户、HEALTHCHECK 打 /healthz")
}
