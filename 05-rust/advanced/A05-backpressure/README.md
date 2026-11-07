# A05 · Rust · 背压与并发上限

## 问题复现

`mpsc::channel()` 是无界的，生产者来多少收多少；`tokio::spawn` 无限制地起任务同理。内存持续上涨，最终 OOM。本例用 `sync_channel(5)` 提交 20 个任务，输出「接收 5、拒绝 15」。

## 错误写法

```rust
let (tx, rx) = mpsc::channel(); // 无界
for i in 0..20 { tx.send(i).unwrap(); } // 永不阻塞
```

无界 channel 的 `send` 永不失败，背压完全失效。

## 正确写法

```rust
match tx.try_send(i) {
    Ok(()) => accepted += 1,
    Err(TrySendError::Full(_)) => rejected += 1, // 满则拒绝
    Err(TrySendError::Disconnected(_)) => break,
}
```

- `sync_channel(capacity)` 是硬上限；
- `send` 阻塞、`try_send` 立即返回 `Full`，语义要明确选；
- `TrySendError` 必须同时处理 `Full` 和 `Disconnected`，否则编译不过——这是 Rust 强制你处理背压分支的好处。

本例用 `Barrier` 先挡住消费者，保证提交阶段队列必满，输出确定。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/A05-backpressure
cargo run --quiet
```

预期：`accepted=5 rejected=15 processed=5`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 有界队列 | `chan` + `cap` | 数组 + 长度判断 | `queue.Queue(maxsize)` | `ArrayBlockingQueue` | `sync_channel` |
| 满时行为 | `select default` | 手动 `if` | `put_nowait` 抛 `Full` | `offer` 返回 `false` | `try_send` 返回 `Full` |
| 错误处理 | 无错误，走 default | 无 | 异常 | 布尔返回值 | `Result` 必须 match |

Rust 把「队列满」编码进 `Result`，不处理就编译不过；Go 用 `select default` 静默走另一条分支；Java 返回布尔值容易被忽略；Python 抛异常但可能被裸 `except` 吞掉。
