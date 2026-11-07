# 23 · 测试、基准与文档

## 这个例子讲什么

入口 `main.go` 自带一套可直接 `go run` 的**表驱动微型测试**与手动计时演示；README 补充标准 `go test` / `Benchmark` / 覆盖率 / 文档注释的写法。入口不依赖任何测试框架，保证能独立运行。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/23-testing-benchmarks-and-docs          # 本课入口，直接跑
go doc  ./basics/23-testing-benchmarks-and-docs.FizzBuzz # 查看文档注释
```

标准 `go test` 写法（把下面内容放进同目录的 `main_test.go` 后运行 `go test -v -cover`、`go test -bench=. -benchmem`）：

```go
func TestFizzBuzz(t *testing.T) {
	cases := []struct{ n int; want string }{
		{1, "1"}, {3, "Fizz"}, {5, "Buzz"}, {15, "FizzBuzz"},
	}
	for _, c := range cases {
		if got := FizzBuzz(c.n); got != c.want {
			t.Errorf("FizzBuzz(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func BenchmarkFizzBuzz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FizzBuzz(i % 15)
	}
}
```

## 关键点

- **表驱动测试**：用例塞进 `[]struct{...}` 切片，一个循环跑完；加用例只改数据表，不改逻辑，是 Go 社区最主流的测试风格。
- **失败信息**：`t.Errorf("got=%q want=%q input=%d", ...)` 把「实际/期望/输入」写全；本例 2 段故意用错误期望演示 FAIL 行格式，但全部真实用例通过、退出码 0。
- **覆盖率思路**：让每条分支（15/3/5/默认）至少有一条用例；`go test -cover` 量化，`go test -coverprofile` 生成报告，`go tool cover -html` 可视化。
- **基准测试**：`testing.B` 的 `b.N` 由框架自动调整，`go test -bench=. -benchmem` 看 ns/op 与分配次数；本例入口用 `time.Now()/time.Since()` 做等价的手动计时，但只打印「耗时 > 0」这类稳定结论，不把具体耗时当结论。
- **文档注释**：紧挨声明上方的 `// Name ...` 注释即 godoc，`go doc` / `pkg.go.dev` 都能读；导出函数应有文档注释。

## 和其他语言的差异

| 语言 | 测试框架 | 基准 | 文档 |
|---|---|---|---|
| Go | 标准库 `testing` + 表驱动，`go test -cover` | `testing.B` + `go test -bench` | `//` 注释 + `go doc` |
| Node.js | `node:test` / `node:assert`，`--test` | `--test` 的 benchmark 或第三方 | JSDoc |
| Python | `unittest` / `pytest` | `timeit`，`pytest-benchmark` | docstring + `pydoc` |
| Java | JUnit 5 | JMH（`org.openjdk.jmh`） | Javadoc `/** */` |
| Rust | `#[cfg(test)]` + `cargo test` | `#[bench]`（nightly）/ `criterion` | `///` 文档注释 + `cargo doc`（rustdoc） |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
