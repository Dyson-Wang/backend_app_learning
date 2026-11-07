# 20 · 异步编程与取消 / 上下文

## 这个例子讲什么

Go 没有 `async/await`：异步单元是 goroutine，结果靠 channel 回传，取消与超时统一交给 `context.Context`。本例覆盖 `WithTimeout`/`WithCancel`、`select` 监听取消、父子取消链传播、等待多个并发任务、取消后清理、上下文传值。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/20-async-and-cancellation
```

## 关键点

- **Context 三件事**：取消信号、截止时间、请求作用域的键值。约定作为函数**第一个参数** `ctx` 一路传下去，不要塞进结构体。
- **取消传播**：`context.WithCancel(parent)` 建出的 child 会随 parent 一起被取消；`pcancel()` 一次调用，`parent.Done()` 与 `child.Done()` 同时关闭。`cancel()` 必须调用（本例 `defer`）以释放 timer/资源。
- **`select` 是取消的落点**：长时间任务要在循环里 `case <-ctx.Done(): return` 主动响应；本例 `worker` 每次迭代重新检查，超时到点即 `return` 并做清理。`ctx.Err()` 给出 `context.Canceled` / `context.DeadlineExceeded`。
- **超时**：`context.WithTimeout(ctx, 25*time.Millisecond)` 到点自动取消，等价于带 deadline 的 `WithDeadline`。
- **并发等待**：`sync.WaitGroup` 适合「只等完成」，要「收结果/超时兜底」就用 channel + `select`，并给 `time.After` 一个兜底分支防止永久阻塞。
- **传值**：`WithValue` 只放请求级元数据（trace id 等），key 建议用自定义类型（本例 `ctxKey`）避免冲突；不要拿来传可选参数。

## 和其他语言的差异

| 语言 | 异步模型 | 取消 / 超时 | 上下文传值 |
|---|---|---|---|
| Go | goroutine + channel，无 async/await | `context.WithCancel/WithTimeout` + `select` | `context.WithValue` |
| Node.js | Promise / `async-await` + 事件循环 | `AbortController` / `AbortSignal.timeout` | `AsyncLocalStorage` |
| Python | `asyncio` Task / 协程 | `asyncio.timeout` / `TaskGroup` / `Task.cancel()` | `contextvars.ContextVar` |
| Java | `CompletableFuture` / 虚拟线程 | `orTimeout` / `Future.cancel` / `interrupt` | `ThreadLocal` / `ScopedValue` |
| Rust | 标准库**无异步运行时**，`tokio::spawn` 是生态方案 | `tokio::time::timeout` / `CancellationToken` / Drop | `tokio::task_local!` / `tracing` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
