# E04 · Rust · 结构化日志 / trace id / 指标

## 问题复现

`println!("请求失败: {}", err)` 这类日志没有字段、没有 trace id，并发下无法归属。线上排查时既不能按 status 聚合，也不能把同一次请求的多条日志串起来。

## 错误写法

```rust
println!("用户请求失败: {}", err);
```

自由文本、无结构、无上下文、无指标。

## 正确写法

1. **结构化日志**：手写一行 JSON，字段固定（零依赖；生产用 `tracing` + `tracing-subscriber` 的 JSON 层）。
2. **trace id**：`AtomicU64` 生成 `trace-0001`，生产用 `tracing::Span` 或请求头透传。
3. **RED 指标**：`Mutex<Metrics>` 保护请求数、5xx 数、总耗时、最大耗时。

Rust 的所有权让「指标累加」必须显式 `lock()`，不会出现忘加锁的静默数据竞争。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/E04-observability
cargo run --quiet
```

预期：10 行 JSON，`trace_id` 从 `trace-0001` 递增，两条 `"level":"error"`；最后 `requests=10 errors=2 error_rate=20.00% avg_ms=55.7 max_ms=260`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 结构化日志 | `log/slog` | `pino` | `logging` + Formatter | Logback + MDC | `tracing` |
| trace 透传 | `context.Context` | `AsyncLocalStorage` | `contextvars` | MDC | `Span` |
| 指标 | Prometheus client | `prom-client` | `prometheus_client` | Micrometer | `metrics` |

Rust 的 `tracing` 把「span（一次请求）」和「event（一条日志）」统一成一个模型，比单纯的日志库更适合异步链路；Go 的 `slog` 是 1.21 才进标准库；Java 的 MDC 要小心线程池复用。
