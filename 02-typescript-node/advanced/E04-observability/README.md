# E04 · Node.js · 结构化日志 / trace id / 指标

## 问题复现

Node 服务里常见的日志是 `console.log('请求失败', err)`：自由文本、没有 trace id、没有状态码和耗时。异步回调一多，多个请求的日志交错在一起，根本分不清哪几行属于同一次请求。

## 错误写法

```js
console.log('请求失败', err);
```

- 无法按字段过滤和聚合；
- 异步交错后无法归属到具体请求；
- 没有指标，错误率只能靠人工统计。

## 正确写法

1. **结构化日志**：`JSON.stringify({ level, trace_id, method, path, status, latency_ms })`，一行一条。
2. **trace id 透传**：用 `AsyncLocalStorage` 在入口 `als.run({ traceId }, ...)`，异步链路里任何位置 `als.getStore().traceId` 都能拿到——这是 Node 版的 `context.Context`。
3. **RED 指标**：请求数、5xx 数、平均/最大耗时，用一个类累加，最后 `snapshot()` 输出。

注意 `als.run` 的回调要包住整个异步流程，否则 `await` 之后拿不到 store。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/E04-observability
node index.mjs
```

预期：10 行 JSON，`trace_id` 从 `trace-0001` 递增，两条 `level: "error"`；最后 `requests=10 errors=2 error_rate=20.00% avg_ms=55.7 max_ms=260`。数字固定，可复现。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 结构化日志 | `log/slog` | `pino` / 手写 JSON | `logging` + `json` | Logback + MDC | `tracing` |
| 上下文透传 | `context.Context` 显式传参 | `AsyncLocalStorage` 隐式 | `contextvars` 隐式 | MDC 线程本地 | `tracing::Span` |
| 指标 | Prometheus client | `prom-client` | `prometheus_client` | Micrometer | `metrics` |

Go 的上下文必须一路显式传，啰嗦但清晰；Node 的 `AsyncLocalStorage` 是隐式的，写起来最省事但要注意作用域；Python 的 `contextvars` 在 asyncio 里正确、多线程里要小心；Java 的 MDC 在虚拟线程下需要显式传递。
