# E04 · Java · 结构化日志 / trace id / 指标

## 问题复现

`log.info("请求失败 " + e.getMessage())` 这类日志没有字段、没有 trace id，线程池里多个请求的日志交错后无法归属。出了 500 只能靠时间戳猜是哪次请求。

## 错误写法

```java
log.info("用户请求失败: " + e.getMessage());
```

- 自由文本，采集端无法按 status / path 聚合；
- 没有 trace id，跨服务无法串联；
- 没有指标，错误率只能人工统计。

## 正确写法

1. **结构化日志**：一行一个 JSON，字段固定。
2. **trace id**：`AtomicLong` 生成 `trace-0001`，生产环境用 MDC 放进日志上下文。
3. **RED 指标**：`synchronized` 保护请求数、5xx 数、总耗时、最大耗时。

本例手写 JSON 是为了零依赖；生产用 Logback + `net.logstash.logback.encoder.LogstashEncoder` 或 Jackson，配合 MDC 自动带上 trace id。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/E04-observability
javac Main.java && java Main
```

预期：10 行 JSON，`trace_id` 从 `trace-0001` 递增，两条 `"level":"error"`；最后 `requests=10 errors=2 error_rate=20.00% avg_ms=55.7 max_ms=260`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 结构化日志 | `log/slog` | `pino` | `logging` + Formatter | Logback + MDC | `tracing` |
| trace 透传 | `context.Context` | `AsyncLocalStorage` | `contextvars` | MDC / `ThreadLocal` | `Span` |
| 指标 | Prometheus client | `prom-client` | `prometheus_client` | Micrometer | `metrics` |

Java 的 MDC 基于 `ThreadLocal`，线程池复用线程时必须 `MDC.clear()`，否则 trace id 会串到下一个请求；虚拟线程下要用 `ThreadLocal` 的替代方案或显式传递；Micrometer 是 Spring Boot 默认的指标门面。
