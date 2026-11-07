# 20 · Rust · 异步编程与取消 / 上下文

## 这个例子讲什么

Rust 的 `Future` 是**惰性状态机**：不 `poll` 就不执行。本例用约 90 行标准库代码手写一个 `block_on` executor（`Wake` + `park/unpark`），再实现 `YieldOnce` 和 `Delay` 两个 Future，最后用 `Timeout` 演示「超时即取消」——内部 Future 不再被 poll。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/20-async-and-cancellation
cargo run --quiet
```

## 关键点

- `Future::poll` 返回 `Poll::Ready` 或 `Poll::Pending`；返回 `Pending` 前**必须**注册 waker，否则永远醒不过来。
- `block_on` 的核心就是 `loop { poll; Pending => park }`，waker 负责 `unpark`。
- 取消是**丢弃**：`Timeout` 返回 `Err` 后，内部 `Future` 再也不会被 poll，这就是 Rust 的取消语义（不需要 `cancel()` 方法）。
- `Pin<Box<F>>` 让 `Timeout` 可以安全地持有任意 `Future` 而不需要 `unsafe`。
- 生产不用自己写 executor：用 tokio / async-std，超时用 `tokio::time::timeout`，取消用 `CancellationToken` 或直接 drop。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 异步模型 | goroutine + 调度器 | 事件循环 + Promise | asyncio 事件循环 | 虚拟线程 / `CompletableFuture` | `Future` + 外部 executor |
| 超时 | `context.WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 取消 | `ctx.Done()` | `AbortController` | `Task.cancel()` | `Future.cancel` | drop future |
| 上下文传值 | `context.WithValue` | `AsyncLocalStorage` | `contextvars` | `ThreadLocal` / MDC | `tower` Extension / 参数传递 |

Go 的取消是显式 `ctx` 一路传；JS/Python 有内建事件循环，`await` 即可；Java 21 虚拟线程让「阻塞写法」也能高并发；Rust 把「谁来 poll」交给运行时库，标准库只定义 `Future` 协议，所以必须选一个 executor。
