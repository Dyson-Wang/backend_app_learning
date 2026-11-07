# 20 · Java · 异步编程与取消 / 上下文

## 这个例子讲什么

`CompletableFuture` 异步任务模型、`orTimeout` 超时、`cancel` 取消与异常沿链传播、`allOf` 并发等待、`whenComplete` 取消后清理、`ExecutorService` 收尾，以及上下文传值。

## 怎么跑

```
cd backend-roadmap/04-java/basics/20-async-and-cancellation
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **`CompletableFuture`**（Java 8）是 Java 的异步 / future 抽象，`supplyAsync(fn, pool)` 把任务交给线程池；`thenApply` 等构成链式流水线。
- **`orTimeout(40ms)`**（Java 9）给永不完成的 future 设超时，内部用守护调度线程，触发 `TimeoutException`——示例不 `sleep`，靠它演示超时。
- **`cancel(true)`** 对未完成 future 直接以 `CancellationException` 结束；`join()` 遇到取消抛 `CancellationException`、遇到普通异常抛 `CompletionException`。上游异常会沿 `thenApply` 链传播。
- **`allOf`** 并发等待多个任务；**`whenComplete`** 无论成功失败都会执行，适合做清理。
- **上下文传值**：Java 没有内置 context，示例显式把 `record Ctx` 沿链传递；`ThreadLocal` 在线程池里会串味，`ScopedValue` 仍是 preview，故不采用。
- 结束前 `shutdownNow()` 终止线程池，保证 2 秒内退出。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 异步模型 | `CompletableFuture` | goroutine + channel | `Promise`/`async` | `asyncio`/`await` | `Future` + 运行时 |
| 超时 | `orTimeout` | `context.WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `tokio::time::timeout` |
| 取消 | `cancel(true)`（协作式） | `context.WithCancel` + `select` | `AbortController` | `Task.cancel()` | `CancellationToken` |
| 并发等待 | `allOf` | `WaitGroup`/`select` | `Promise.allSettled` | `TaskGroup`/`gather` | `join!`/`try_join!` |
| 上下文 | 显式传值 / `ThreadLocal` | `context.Context` | `AsyncLocalStorage` | `contextvars` | `task_local!` |
