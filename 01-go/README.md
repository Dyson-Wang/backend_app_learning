# Go · 01-go

**运行时模型**：CSP 并发（goroutine / channel）· GC · 静态类型 · 编译成单二进制

## 环境准备

本机当前**未安装 Go**（2026-11-07 检测），示例代码尚未编译验证。安装后再跑：

```bash
brew install go          # 或官网 pkg
cd 01-go
go build ./...
go vet ./...
go run ./basics/01-hello-and-toolchain
```

`go.mod` 只声明标准库，不需要 `go mod download`。

## 目录

- `basics/` —— 基础入门 10 课（统一编号）
- `advanced/` —— 生产场景（统一 A01 ~ F03 编号）

## 为什么第一门是 Go

概念最少、单二进制、一个 `net/http` 就能起服务。先用它把「路由 / 中间件 / 数据层 / 并发 / 部署」这条骨架走通，后面换语言时只需要替换实现，不用重新理解后端是什么。

## 与你现有项目的对照

`linnet-server` 是 Rust + axum + tokio + sqlx。Go 这一层的 `middleware` 链可以对照 axum 的 `tower` Layer，`context.Context` 对照 Tokio 的取消机制——这两个对照点在 `../../99-cross-language/` 里会展开。

## 跨语言对照

每课 / 每个场景做完，把「这门语言和其他语言的差异」补到 `../99-cross-language/scenario-matrix.md`。同一编号横向打开，是这套库唯一正确的读法。
