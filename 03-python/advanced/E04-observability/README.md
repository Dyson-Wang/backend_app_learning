# E04 · Python · 结构化日志 / trace id / 指标

## 问题复现

`logging.info("请求失败 %s", err)` 这种自由文本日志，在并发下会交错、无法按字段聚合，也没有 trace id 把同一次请求的日志串起来。出故障时只能靠时间戳猜。

## 错误写法

```python
logging.info("用户请求失败: %s", err)
```

- 无结构化字段，采集端无法按 status / path 聚合；
- 无线程/协程上下文，日志归属不明；
- 无指标，错误率靠人工数。

## 正确写法

1. **结构化日志**：自定义 `JsonFormatter`，把 `level` / `trace_id` / `method` / `path` / `status` / `latency_ms` 输出成一行 JSON。
2. **trace id 透传**：用 `contextvars.ContextVar`，它在 asyncio 和多线程下都能正确隔离；入口 `set()`，出口 `reset()`。
3. **RED 指标**：请求数、5xx 数、平均/最大耗时。

`logger.error("", extra=extra)` 把业务字段挂到 `LogRecord` 上，Formatter 再取出来——这是 Python 标准库最干净的扩展方式。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/E04-observability
python3 main.py
```

预期：10 行 JSON，`trace_id` 从 `trace-0001` 递增，两条 `"level": "error"`；最后 `requests=10 errors=2 error_rate=20.00% avg_ms=55.7 max_ms=260`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 结构化日志 | `log/slog` | `pino` / 手写 JSON | `logging` + 自定义 Formatter | Logback + MDC | `tracing` |
| 上下文透传 | `context.Context` 显式 | `AsyncLocalStorage` | `contextvars` | MDC / ThreadLocal | `tracing::Span` |
| 指标 | Prometheus client | `prom-client` | `prometheus_client` | Micrometer | `metrics` |

Python 的 `contextvars` 是 3.7+ 的标准方案，比 `threading.local` 更适合 asyncio；Go 坚持显式传 context；Node 的 `AsyncLocalStorage` 是隐式的；Java 的 MDC 在线程池里需要手动清理，否则会串号。
