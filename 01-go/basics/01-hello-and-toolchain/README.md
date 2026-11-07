# 01 · hello 与工具链

## 这个例子讲什么

最小的 Go 程序：包声明、import、main 函数，以及怎么看当前运行时信息。顺便确认工具链装好了。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/01-hello-and-toolchain

# 也可以编译成单二进制再执行
go build -o /tmp/brm-hello ./basics/01-hello-and-toolchain && /tmp/brm-hello

# 常用工具链命令
go mod tidy     # 整理依赖
go vet ./...    # 静态检查
gofmt -w .      # 格式化
go test ./...   # 跑测试
```

## 关键点

- **只有 `main` 包 + `main()` 才能生成可执行文件**，其他包都是库。
- `go run` 编译到临时目录再执行；`go build` 产出**单二进制**，这是 Go 部署简单的根本原因（对照 Java 需要 JVM、Node/Python 需要运行时）。
- `GOMAXPROCS` 决定同时执行用户级 goroutine 的 OS 线程数，默认等于 CPU 核数。
- 导出规则：标识符首字母大写才对外可见，没有 `public` 关键字。

## 和其他语言的差异

| 语言 | 入口与交付 |
|---|---|
| Go | `package main` + `func main()`，编译成单二进制 |
| Node.js | 直接执行 `.mjs` 文件，需要 Node 运行时 |
| Python | 直接执行 `.py`，需要 CPython 运行时 |
| Java | `public static void main`，编译成 `.class` 由 JVM 执行 |
| Rust | `fn main()` + `cargo run`，编译成单二进制 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
