# 09 · 并发原语入门

## 这个例子讲什么

**五种并发模型第一次正面相遇的一课。** Go 的答案：goroutine 做并发单元、channel 做通信、`context` 做取消与超时。

## 怎么跑

```
cd backend-roadmap/01-go

go run ./basics/09-concurrency-primitives
go run -race ./basics/09-concurrency-primitives   # 竞态检测，务必试一次
```

## 关键点

- **CSP**：「通过通信共享内存，而不是通过共享内存通信」。但 Go 并不禁止共享，`sync.Mutex` 一样常用；读多写少用 `sync.RWMutex`。
- goroutine 是用户态轻量协程（初始栈约 2KB），由 Go 调度器复用到少量 OS 线程上，开几万个不稀奇。
- **channel 的关闭语义**：发送方负责 `close`，接收方 `range` 或 comma-ok；向已关闭的 channel 发送会 panic；无缓冲 channel 同时承担同步。
- `select` 多路复用，`time.After` 给一个兜底超时；循环里反复 `time.After` 会累积 timer，长生命周期用 `time.NewTimer` 复用。
- `context.Context` 承载**取消信号、截止时间、请求作用域的值**，应作为函数第一个参数一路传下去；跨服务调用要把它转成 HTTP header / MQ 属性。
- `go run -race` 能抓到真实数据竞争。把本例的 `mu.Lock/Unlock` 去掉再跑一次，看检测报告长什么样。

## 和其他语言的差异

| 语言 | 并发单元 | 通信 / 共享 | 取消 |
|---|---|---|---|
| Go | goroutine | channel / `sync.Mutex` | `context.Context` |
| Node.js | Promise + 事件循环 | 单线程，`AsyncLocalStorage` 传上下文 | `AbortSignal` |
| Python | asyncio Task / 多进程 | `asyncio.Queue` / `Lock` | `Task.cancel()` |
| Java | 平台线程 / 虚拟线程 | `BlockingQueue` / `Lock` | `interrupt` / 超时 |
| Rust | `tokio::spawn` | `mpsc` / `Arc<Mutex<T>>` | `CancellationToken` / Drop |

注意 Node.js 是唯一「默认不需要锁」的：不是因为更安全，而是因为只有一个线程执行 JS，代价是任何同步阻塞都会拖垮整个进程。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
