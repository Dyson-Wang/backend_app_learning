# F02 重试风暴

## 问题复现

下游超时，100 个客户端同时重试 3 次，100 QPS 变 400 QPS，打在一个已过载的服务上。

## 错误写法

```rust
for _ in 0..3 {
    if call().await.is_ok() { return Ok(()); }   // 立即重试
}
```

## 正确写法

```rust
for i in 0..max_attempts {
    if call().await.is_ok() { return Ok(()); }
    let base = Duration::from_millis(100 * (1 << i));
    let jitter = Duration::from_millis(rand_u64() % 100);
    tokio::time::sleep(base + jitter).await;
}
```

## 怎么验证

```bash
cargo run --quiet
```

预期：`400 QPS` → 退避后 `133 QPS` → 熔断后 `0 QPS`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 退避库 | `backoff` | `p-retry` | `tenacity` | Resilience4j | `backoff` crate |
| 熔断 | `gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |
| 陷阱 | 重试非幂等写 | 同 | 同 | 同 | 同 |
