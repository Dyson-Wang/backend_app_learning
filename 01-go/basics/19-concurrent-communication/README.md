# 19 · 并发通信：channel / 队列 / 同步原语

## 这个例子讲什么

把 09 课的并发原语推向工程模型：有界队列与背压、生产消费、扇出扇入（fan-out/fan-in）、`WaitGroup` 等待多任务、`close` 关闭信号广播，以及 `sync/atomic` 与 `sync.Mutex` 两种共享状态保护。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/19-concurrent-communication
go run -race ./basics/19-concurrent-communication   # 竞态检测
```

## 关键点

- **有界队列 = 带缓冲 channel**：`make(chan int, 3)`，缓冲满时发送方阻塞，天然形成背压；无缓冲 channel 则同时承担同步。
- **关闭信号传播**：发送方 `close(ch)`，接收方 `for range ch` 会在取完剩余值后自动结束。`close` 是**一对多广播**，本例 `close(done)` 一次唤醒 3 个等待者，比挨个发信号干净。
- **扇出扇入**：多个 worker 从同一 `src` 竞争消费是扇出；每个 worker 的结果写自己的 `out` channel，再用 `merge` 汇总成一条流是扇入。`merge` 用 `WaitGroup` 等所有输入关闭后才 `close(out)`，避免提前关流。
- **等待多个任务**：`sync.WaitGroup` 的 `Add/Done/Wait` 是标准做法；并发结果顺序不定，所以本例收集后 `sort.Ints` 保证输出可复现。
- **共享状态**：`atomic.AddInt64` 适合单纯计数，`sync.Mutex` 适合保护一段临界区；读多写少换 `sync.RWMutex`。去掉 `mu.Lock` 用 `go run -race` 能抓到真实数据竞争。

## 和其他语言的差异

| 语言 | 队列 / 通信 | 等待多任务 | 共享状态 |
|---|---|---|---|
| Go | channel（有界/无界）、`close` 广播 | `sync.WaitGroup`、`select` | `sync.Mutex`/`RWMutex`、`sync/atomic` |
| Node.js | 事件循环单线程；`worker_threads` + `MessageChannel` | `Promise.all` / `allSettled` | 几乎无需锁（单线程），跨线程用 `SharedArrayBuffer`/`Atomics` |
| Python | `queue.Queue`、`asyncio.Queue` | `concurrent.futures` / `asyncio.gather` | `threading.Lock`、`multiprocessing` 绕开 GIL |
| Java | `BlockingQueue`（`ArrayBlockingQueue` 有界）、`CountDownLatch` | `ExecutorService` / `CompletableFuture.allOf` | `synchronized`、`ReentrantLock`、`AtomicInteger` |
| Rust | `std::sync::mpsc`（多生产者单消费者）、`tokio::sync::mpsc` | `JoinHandle` / `tokio::join!` | `Arc<Mutex<T>>`、`AtomicUsize` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
