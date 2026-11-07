# B05 上下文取消与跨服务传播

## 问题复现

上游给 50ms 预算，下游要 200ms。中间层不传截止时间，下游跑满 200ms——上游早断了。

## 错误写法

```rust
downstream(None, work_ms);      // 忘了传
```

## 正确写法

```rust
let left = d.remaining_ms();
if (left as u64) < work_ms {
    thread::sleep(Duration::from_millis(left as u64));
    return "deadline exceeded";
}
```

Rust 没有隐式上下文，deadline 只能靠参数或 task-local 传；忘记传会表现为类型不匹配，比 Go 的 `context.Background()` 误用更容易被发现。

## 怎么验证

```bash
cargo run --quiet
```

预期：第一种约 50ms 返回 `deadline exceeded`；第二种 200ms 返回 `ok`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 载体 | `context.Context` | `AbortSignal` | 显式参数 | `orTimeout` | `tokio::time::timeout` |
| 忘记传的后果 | 静默失去取消 | 静默失去取消 | 静默失去取消 | 静默失去取消 | 编译期类型不匹配 |
| 跨服务 | header | header | header | header | header |
