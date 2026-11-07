# 20 · Python · 异步编程与取消 / 上下文

## 这个例子讲什么

`asyncio.run` 任务模型、`asyncio.timeout`（3.11+）、`asyncio.TaskGroup`（3.11+）、`CancelledError` 与 `finally` 清理、`ExceptionGroup` + `except*`、`contextvars.ContextVar` 的任务间隔离与 `copy_context`。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/20-async-and-cancellation/main.py
```

## 关键点

- **`asyncio.timeout`**：超时会取消块内协程，退出时把 `CancelledError` 转成内建 `TimeoutError`（不是 `asyncio.TimeoutError`）；被取消的协程用 `finally` 做清理。
- **`TaskGroup`**：`async with` 退出时等待所有子任务；任一任务抛异常会**取消同组其余任务**，并把异常包成 `ExceptionGroup`，用 `except*` 按类型分流捕获。
- `task.cancel()` 只是**请求**取消，真正生效点在下一次 `await`；`await task` 会抛 `CancelledError`。
- `ContextVar` 的值随 `Task` 自动复制一份，子任务 `set()` 不影响父任务——这正是 Go `context` 传值、Node `AsyncLocalStorage` 的对应物。
- 全部 `sleep` 都是毫秒级，整个脚本 0.2 秒内自行结束，不监听端口、不常驻。

## 和其他语言的差异

- **Go**：`context.WithTimeout` / `WithCancel` + `select` + goroutine 取消传播，`errgroup` 聚合错误。
- **Node.js/TS**：`AbortController` / `AbortSignal.timeout`、`Promise.allSettled`、`AsyncLocalStorage`。
- **Java**：`CompletableFuture.orTimeout`、`ExecutorService.shutdownNow`、`Future.cancel(true)`。
- **Rust**：标准库无异步运行时，生产用 `tokio::time::timeout` / `CancellationToken` / `JoinSet`。
- **Python**：`asyncio.run`、`asyncio.timeout`、`TaskGroup`、`CancelledError`、`ExceptionGroup` + `except*`、`contextvars`。
