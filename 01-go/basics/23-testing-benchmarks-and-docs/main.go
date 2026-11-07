package main

import (
	"fmt"
	"time"
)

// FizzBuzz 是带文档注释的导出函数：3 的倍数返回 Fizz，5 的倍数返回 Buzz，
// 15 的倍数返回 FizzBuzz，其余返回数字字符串。可用 go doc 查看这段注释。
func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return fmt.Sprint(n)
	}
}

func main() {
	// 1. 表驱动测试：用例集中在切片里，逐个断言
	cases := []struct {
		name string
		n    int
		want string
	}{
		{"n=1", 1, "1"},
		{"n=3", 3, "Fizz"},
		{"n=5", 5, "Buzz"},
		{"n=15", 15, "FizzBuzz"},
		{"n=7", 7, "7"},
		{"n=30", 30, "FizzBuzz"},
	}
	pass, fail := 0, 0
	for _, c := range cases {
		got := FizzBuzz(c.n)
		if got != c.want {
			fail++
			fmt.Printf("FAIL %-6s: got=%q want=%q input=%d\n", c.name, got, c.want, c.n)
		} else {
			pass++
			fmt.Printf("PASS %-6s: %q\n", c.name, got)
		}
	}
	fmt.Printf("汇总: %d 通过, %d 失败（退出码仍为 0）\n", pass, fail)

	// 2. 失败信息长什么样：故意用错误期望演示格式，不影响上面结果
	wrong := "Fizz!"
	fmt.Printf("失败信息示例: FAIL badcase: got=%q want=%q input=%d\n", FizzBuzz(3), wrong, 3)

	// 3. 覆盖率思路：每个分支至少一条用例
	fmt.Println("覆盖率思路: 6 个用例覆盖 15/3/5/默认 四个分支；用 go test -cover 可量化")

	// 4. 计时/基准演示：手动 micro-benchmark（对应 testing.B 的 Benchmark）
	const rounds = 20000
	start := time.Now()
	hits := 0
	for i := 0; i < rounds; i++ {
		if FizzBuzz(i%15) == "FizzBuzz" {
			hits++
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("计时演示: %d 次调用, 命中 FizzBuzz = %d, 计时有效(>0) = %t\n",
		rounds, hits, elapsed > 0)

	// 5. 文档注释生成方式
	fmt.Println("文档: 函数上方的 // 注释即 godoc；go doc ./basics/23-testing-benchmarks-and-docs 可查看")
}
