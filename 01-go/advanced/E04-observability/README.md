# E04 · Go · 结构化日志 / trace id / 指标

## 问题复现

线上排查最常见的困境：日志是 `fmt.Printf("请求失败")` 这种自由文本，没有 trace id、没有状态码、没有耗时，出问题只能靠时间戳猜。或者有日志但没有指标，无法回答「错误率是多少、P99 多慢」。

## 错误写法

```go
log.Printf("用户请求失败: %v", err)
```

- 没有结构化字段，日志系统无法按 status / path 聚合；
- 没有 trace id，跨服务串不起来；
- 没有指标，只能人肉数日志。

## 正确写法

三件事分开做：

1. **结构化日志**：一行一个 JSON，字段固定（`level` / `trace_id` / `method` / `path` / `status` / `latency_ms`），采集端直接解析。
2. **trace id 透传**：入口生成，放进 context 或请求头，沿途所有日志都带上同一个 id。
3. **RED 指标**：Rate（请求数）、Errors（5xx 数）、Duration（平均/最大耗时），用 `sync.Mutex` 保护累加。

本例把三者串在一个 `handle()` 里：先取 trace id，再记指标，最后打日志。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/E04-observability
go run .
```

预期：10 行 JSON 日志，`trace_id` 从 `trace-0001` 递增；最后输出 `requests=10 errors=2 error_rate=20.00%`，`max_ms=260`。所有数字固定，可复现。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 结构化日志 | `log/slog` | `pino` / 手写 JSON | `logging` + `json` | Logback + MDC | `tracing` |
| trace 上下文 | `context.Context` | `AsyncLocalStorage` | `contextvars` | MDC / ThreadLocal | `tracing::Span` |
| 指标 | Prometheus client | `prom-client` | `prometheus_client` | Micrometer | `metrics` |

Go 的 `slog` 已进标准库；Node 用 `AsyncLocalStorage` 才能在异步回调里拿到 trace id；Python 的 `contextvars` 是 asyncio 场景的正确选择；Java 的 MDC 依赖线程本地变量，虚拟线程下要注意传递；Rust 的 `tracing` 把 span 和日志合成一个概念。
