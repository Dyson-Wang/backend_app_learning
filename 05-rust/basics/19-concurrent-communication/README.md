# 19 · Rust · 并发通信：channel / 队列 / 同步原语

## 这个例子讲什么

Rust 并发的两条正统路线：**消息传递**（`mpsc::channel` / `sync_channel`）和**共享状态**（`Arc<Mutex<T>>`）。覆盖单生产者、多生产者、有界队列背压、`join` 等待、关闭信号传播——正好是生产里最常用的五种模式。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/19-concurrent-communication
cargo run --quiet
```

## 关键点

- `mpsc` = multi-producer, single-consumer；多生产者靠 `tx.clone()`，消费者只能有一个。
- 多生产者场景**必须 `drop(tx)`**，否则原始 sender 一直活着，`rx.iter()` 永远等不到通道关闭。
- `sync_channel(n)` 有界，队列满时 `send` 阻塞——这就是最朴素的背压，比无界队列安全得多。
- `Arc<Mutex<T>>` 是共享可变的标配；`lock()` 返回 `LockResult`，`.unwrap()` 会在持锁线程 panic 时传播。
- 关闭信号用「发一个空消息 + `try_recv`」或「drop sender 让 `recv` 返回 `Err`」，不要用共享 bool 轮询。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 通信原语 | channel + `select` | `MessageChannel` / 事件 | `queue.Queue` | `BlockingQueue` | `mpsc` channel |
| 共享状态 | `sync.Mutex` | 单线程，无锁 | `threading.Lock`（GIL） | `synchronized` / `Lock` | `Arc<Mutex<T>>` |
| 等待全部完成 | `sync.WaitGroup` | `Promise.all` | `ThreadPoolExecutor` | `CompletableFuture.allOf` | `join()` |
| 关闭传播 | `close(ch)` | `AbortController` | 哨兵值 / `Event` | `interrupt` / poison | drop sender / 哨兵消息 |

Go 的 `select` 能同时等多个 channel，Rust 标准库没有等价物（要 `crossbeam`）；Java 的 `BlockingQueue` 最成熟；Python 有 GIL，多线程只适合 I/O 密集；Rust 用所有权在编译期拦住「多个线程同时改同一块内存」。
