# 20 · Node.js · 异步编程与取消 / 上下文

## 这个例子讲什么

把 Node 异步编程的四块拼齐：`AsyncLocalStorage` 做上下文传值（跨 await/微任务/宏任务不丢）、`AbortController` 主动取消并把信号传播给多个子任务、`AbortSignal.timeout` 超时自动取消、`AbortSignal.any` 组合多个取消源、`Promise.allSettled` 并发等待，以及取消后在 `abort` 监听里清理定时器资源。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/20-async-and-cancellation/index.mjs
```

只用 `node:async_hooks`；唯一的真实等待是 `AbortSignal.timeout(30)` 的 30ms，远低于限制。

## 关键点

- `AsyncLocalStorage.run(store, fn)` 把整条异步调用链绑定到一个 store，`als.getStore()` 在 `await` / `Promise` / `setImmediate` 之后仍能读到，等价于给每个请求挂 `requestId`，不必层层传参（10 课的访问日志正该用它）。
- `AbortController` + `signal`：调用 `controller.abort(reason)` 会同步触发所有监听者并把 `signal.reason` 传播给共享该信号的任务——这就是**取消传播**。
- `AbortSignal.timeout(ms)` 到点自动 abort，`reason.name === 'TimeoutError'`；它的内部定时器是 **unref 的**，不会阻止进程退出，所以必须让任务自身的 handle（示例里 500ms 定时器）撑住事件循环。
- `AbortSignal.any([...])` 把多个取消源合成一个：任一触发即取消，并沿用最先触发者的 `reason`。
- **取消后清理**：在 `signal.addEventListener('abort', ..., {once:true})` 里 `clearTimeout(timer)` 释放资源，否则被取消的任务仍可能拖着定时器让进程不退出。
- `Promise.allSettled` 永远收齐所有结局（`status`/`reason`），不像 `Promise.all` 遇到一个 reject 就短路，适合批量任务的取消汇总。
- `AbortSignal` 是 Web 标准 API，Node 的 `fetch`、`fs` 流、定时器等内置接口普遍接受 `{ signal }` 参数，形成统一的取消协议。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 取消原语 | `context.Context` | `asyncio.CancelledError` | `Future.cancel` / 中断 | `CancellationToken`（tokio） | `AbortSignal` / `AbortController` |
| 超时 | `context.WithTimeout` | `asyncio.timeout` / `wait_for` | `orTimeout` / `TimeoutException` | `tokio::time::timeout` | `AbortSignal.timeout` |
| 组合取消 | `context` 树派生 | `TaskGroup` | 组合 `CompletableFuture` | `CancellationToken` 子 token | `AbortSignal.any` |
| 上下文传值 | `context.WithValue` | `contextvars` | `ThreadLocal` | `task_local!` | `AsyncLocalStorage` |
| 并发等待 | `errgroup` / `WaitGroup` | `asyncio.gather(return_exceptions)` | `allOf` | `join!` / `FuturesUnordered` | `Promise.allSettled` |
| 清理 | `defer` + `<-ctx.Done()` | `finally` / `async with` | `finally` | `Drop` | `abort` 监听 + `finally` |

Node 的 `AbortSignal` 直接对应 Go 的 `context.Context`：都是显式沿调用链传递的取消/超时信号；`AsyncLocalStorage` 对应 Go 的 `context.WithValue` / Python 的 `contextvars`（但 ALS 是隐式透传，不需显式传参）。Java 用 `Future.cancel` + `ThreadLocal`，Rust 标准库没有异步运行时，生产上靠 tokio 的 `CancellationToken` 与 `timeout`。Go 的 `errgroup.WithContext` 把「取消 + 等待多个任务」合成一体，是最接近本课 `AbortSignal + allSettled` 组合的写法。
