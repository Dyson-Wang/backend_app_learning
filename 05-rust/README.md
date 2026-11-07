# Rust · 05-rust

**运行时模型**：所有权 + 借用检查 · 无 GC · async / Tokio · 静态强类型 · 单二进制

## 环境准备

本机**未安装 cargo / rustc**，示例尚未编译验证。安装后：

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
cd 05-rust/basics/01-hello-and-toolchain
cargo run
```

`01` 只用标准库；`10` 用标准库 `TcpListener` 手写一个最小 HTTP 服务，保证零依赖。`advanced/` 会引入 axum / tokio / sqlx，每个目录单独给 `Cargo.toml`。

## 目录

- `basics/` —— 基础入门 10 课（统一编号）
- `advanced/` —— 生产场景（统一 A01 ~ F03 编号）

## 重点补什么

编译器在编译期就消灭数据竞争和悬垂引用，代价是所有权与生命周期的心智负担。`advanced/` 里会重点对照：Go 的 `context.Context` vs Tokio 取消、`sync.Mutex` vs `tokio::sync::Mutex`、GC 语言里「随手共享」在 Rust 里要怎么表达。

## 与你现有项目的对照

你已经有 `linnet-server`（axum 0.8 + tokio + sqlx + tower-http + JWT/bcrypt + tracing + reqwest），这已经是相当完整的生产级后端。这一门不必从零入门，**可以直接拿它当 `advanced/` 的真实案例库**：鉴权对应 `E01/E02`、request-id + trace 对应 `E04`、流式响应对应 `D04`、限流中间件对应 `A01`。

## 跨语言对照

每课 / 每个场景做完，把「这门语言和其他语言的差异」补到 `../99-cross-language/scenario-matrix.md`。同一编号横向打开，是这套库唯一正确的读法。
