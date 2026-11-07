# D02 死信队列与重试策略

## 问题复现

永久错误的消息被反复重投，卡住队头；或重试无上限，把下游打挂。

## 错误写法

```rust
loop {
    if handle(m).is_ok() { break; }
}
```

## 正确写法

用 enum 表达「可重试 / 永久」，让类型系统逼你分类：

```rust
enum Failure { Retryable, Permanent }
match classify(&err) {
    Failure::Permanent => to_dlq(m),
    Failure::Retryable if attempt >= max => to_dlq(m),
    Failure::Retryable => sleep(backoff(attempt)),
}
```

## 怎么验证

```bash
cargo run --quiet
```

预期：`成功 = 3，死信 = 2，总投递 = 9`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 错误分类 | `errors.Is` / 哨兵错误 | `instanceof` | 异常类型 | 异常子类 | enum |
| DLQ | 另建 topic | 同 | 同 | RabbitMQ DLX | 同 |
| 退避 | `time.Sleep` | `setTimeout` | `time.sleep` | `Thread.sleep` | `thread::sleep` |
