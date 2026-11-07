# C04 · Rust · 缓存穿透 / 击穿 / 雪崩

## 问题复现

热点 key 过期瞬间，50 个线程同时发现缓存为空，全部回源。数据库连接被打满，接口 P99 暴涨。本例用 20ms 的假数据库 + 原子计数器复现。

## 错误写法

```rust
if let Some(v) = cache.lock().unwrap().get(key) {
    return v.clone();
}
let value = db.query(key); // 所有并发线程都到这里
cache.lock().unwrap().insert(key.to_string(), value);
```

每次操作都加锁，但「查-回源-写」三步之间没有互斥。

## 正确写法

按 key 加锁 + 双重检查：

1. 查 `data`，命中返回；
2. 从 `key_locks` 取该 key 的 `Arc<Mutex<()>>`，`entry().or_insert_with` 保证只创建一次；
3. 锁住 key 后**再查一次** `data`；
4. 仍未命中才 `db.query` 并写回。

注意每次 `lock()` 都放进独立作用域，避免 `MutexGuard` 跨越 `return` 造成借用错误——这是 Rust 版最容易踩的坑。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/C04-cache-stampede
cargo run --quiet
```

预期：`无保护 ... 回源 50 次`，`按 key 加锁 + 双重检查 ... 回源 1 次`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单飞实现 | per-key mutex | in-flight Promise Map | per-key `Lock` | `computeIfAbsent` | `entry().or_insert_with` |
| 锁类型 | `sync.Mutex` | 无 | `threading.Lock` | `synchronized` | `Arc<Mutex<()>>` |
| 易错点 | 锁粒度 | 忘记 finally | 锁内 I/O | 慢操作进 computeIfAbsent | guard 跨 return / 跨 await |

Rust 的 `MutexGuard` 会在作用域结束时自动释放，比手动 unlock 安全；但 async 里不能跨 `.await` 持 `std::sync::Mutex`，要用 `tokio::sync::Mutex`。
