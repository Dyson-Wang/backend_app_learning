# 09 · Node.js · 并发原语入门

## 这个例子讲什么

用一段确定性输出讲清 Node 单线程事件循环的任务调度顺序：同步代码 → ES 微任务（`queueMicrotask`/`Promise.then`）→ 宏任务（`setImmediate`），单独演示 Node 专属的 `process.nextTick` 高优先级队列，再演示 `Promise.all`/`race`/`allSettled`、`async/await` 错误传播，以及 Promise 执行器同步执行这一易错点。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/09-concurrency-primitives/index.mjs
```

无端口、无定时器长睡、无标准输入，瞬间结束。

## 关键点

- Node 只有一个 JS 线程，I/O 交给 libuv；「并发」是任务交错而非并行（真正并行的 `worker_threads` 见 19 课）。
- ES 微任务（`queueMicrotask` / `Promise.then`）先于 `setImmediate`（check 阶段的宏任务）；示例用数组收集顺序后一次打印，不依赖真实耗时，输出稳定可复现。
- `process.nextTick` 是 Node 专属队列，会连同回调里嵌套排入的 nextTick 一起、在事件循环前进前被完整排空；它与 ES 微任务的相对先后在 ESM 与 CJS 顶层脚本里并不相同，别依赖这个细节写业务逻辑。
- `Promise.all` 全成功才返回、任一失败即 reject；`allSettled` 永远收齐所有状态；`race` 取最先敲定者。
- `async` 函数永远返回 Promise，`throw` 变成 rejected，靠 `try { await } catch` 接住。
- `new Promise(executor)` 的 executor 是**同步**执行的，只有 `.then` 回调进微任务队列。
- 任何同步阻塞（大循环、`*Sync` API、巨型 JSON 解析）都会卡住整个进程的事件循环。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 并发载体 | goroutine（M:N） | 线程 / `asyncio` 协程 | `Thread` / 虚拟线程 | `thread` / async task | 事件循环 + Promise |
| 调度 | Go 运行时抢占 | GIL + 事件循环 | OS 线程 / ForkJoinPool | 需 tokio 运行时 | libuv 单线程事件循环 |
| 等待多个 | `sync.WaitGroup` / `errgroup` | `asyncio.gather` | `CompletableFuture.allOf` | `join!` / `FuturesUnordered` | `Promise.all` |
| 通信 | channel | `queue.Queue` / `asyncio.Queue` | `BlockingQueue` | `mpsc` channel | 回调 / `MessageChannel`（19 课） |
| 取消 | `context.Context` | `asyncio.CancelledError` | `Future.cancel` | `CancellationToken` | `AbortSignal`（20 课） |

Go 的 goroutine 与 Rust/Java 的线程是抢占式并行，能真正吃满多核；Node 与 Python `asyncio` 是协作式单线程调度，CPU 密集任务必须挪到 worker/子进程。取消与上下文传值在 Node 里对应 `AbortSignal` + `AsyncLocalStorage`，见 20 课。
