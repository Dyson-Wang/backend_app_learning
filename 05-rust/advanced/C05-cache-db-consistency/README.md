# C05 缓存与 DB 双写一致性

## 问题复现

「删缓存 + 写库」之间来读请求：缓存已删 → 回源读旧值 → 写回缓存。写请求随后更新库，缓存永远旧。

## 错误写法

```rust
cache.delete(key);
db.update(key, v);      // 中间有读请求回填旧值
```

## 正确写法

```rust
db.update(key, v);
cache.delete(key);
// 延迟再删一次，兜住迟到的回填
```

## 怎么验证

```bash
cargo run --quiet
```

预期：前两种 `是否脏 = true`，双删后 `缓存命中 = false`，回源拿到 `v2`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 客户端 | `go-redis` | `ioredis` | `redis-py` | Lettuce | `redis-rs` |
| 延迟双删 | `time.AfterFunc` | `setTimeout` | `threading.Timer` | `ScheduledExecutor` | 线程 |
| 陷阱 | goroutine 泄漏 | 定时器不 unref | Timer 非 daemon | 线程池打满 | 线程 join |
