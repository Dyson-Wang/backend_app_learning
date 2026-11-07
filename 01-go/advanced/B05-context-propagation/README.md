# B05 · Go · 上下文取消与跨服务传播

## 问题复现

前端 50ms 就超时断开了，后端还在跑 200ms 的查询，连接、goroutine、数据库资源全被占着。更糟的是中间层某处写了 `context.Background()`，把取消信号截断，上游取消了也没用。

## 错误写法

```go
func service() {
    return repo(context.Background()) // 丢掉了调用方的 ctx
}
```

## 正确写法

- `ctx` 作为第一个参数一路显式传递；
- 只有最外层 `WithTimeout` / `WithCancel`，中间层**只传不改**；
- 每个可能阻塞的操作都要 `select { case <-ctx.Done(): ... }`；
- 错误用 `%w` 包装，保留 `context.Canceled` / `DeadlineExceeded` 语义，方便上层用 `errors.Is` 判断。

本例第 3 段故意演示 `context.Background()` 截断，耗时回到 200ms。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/B05-context-propagation
go run .
```

预期：第 1 段 `err = context deadline exceeded`，耗时约 50ms；第 2 段 `context canceled`，约 20ms；第 3 段耗时约 200ms。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 取消载体 | `context.Context` | `AbortSignal` | `asyncio.Task` / `CancelledError` | `CompletableFuture` / 中断 | drop future / `CancellationToken` |
| 传值 | `context.WithValue` | `AsyncLocalStorage` | `contextvars` | `ThreadLocal` / MDC | `tower` Extension |
| 超时 | `WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |

Go 是唯一把「取消 + 超时 + 传值」统一进一个 `Context` 接口的语言，代价是每个函数都要多一个参数；Node/Python 靠隐式上下文，写起来省事但容易忘记传；Rust 的取消是「drop future」，最彻底也最不直观。
