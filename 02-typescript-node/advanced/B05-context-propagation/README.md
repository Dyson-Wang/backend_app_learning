# B05 · Node.js · 上下文取消与跨服务传播

## 问题复现

前端 50ms 就断开，后端还在跑 200ms 的查询。Node 里更常见的是：某层忘了传 `signal`，`fetch` 继续跑，连接和内存都占着。

## 错误写法

```js
async function service() {
  return repo(); // 丢掉了 signal
}
```

## 正确写法

- `AbortController` 在最外层创建，`signal` 一路往下传；
- 所有可取消的异步操作都要监听 `abort` 事件并清理定时器；
- `fetch(url, { signal })` 原生支持；
- 错误用 `controller.abort(reason)` 带原因，便于区分「超时」和「客户端断开」。

本例第 3 段故意不传 signal，耗时回到 200ms。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/B05-context-propagation
node index.mjs
```

预期：第 1 段 `deadline exceeded`，约 50ms；第 2 段 `client closed`，约 20ms；第 3 段约 200ms。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 取消载体 | `context.Context` | `AbortSignal` | `asyncio.Task` | `CompletableFuture` / 中断 | drop future |
| 传值 | `WithValue` | `AsyncLocalStorage` | `contextvars` | `ThreadLocal` / MDC | `tower` Extension |
| 超时 | `WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |

`AbortSignal.timeout(50)` 是最简洁的超时写法，但它只负责触发，不负责清理资源——定时器、监听器还是要在 `abort` 回调里自己清。
