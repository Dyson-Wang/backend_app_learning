# B04 · Rust · 竞态、死锁与调试

## 问题复现

Rust 编译器不允许安全代码里出现数据竞争，所以「丢更新」这类问题在 `cargo check` 阶段就过不去。但**死锁仍然可能**：两个线程各持一把 `Mutex`，反序获取对方那把，双方永久阻塞。

## 错误写法

```rust
// 线程 A：先 a 后 b
let _ga = a.lock().unwrap();
let _gb = b.lock().unwrap();

// 线程 B：先 b 后 a  → 死锁
let _gb = b.lock().unwrap();
let _ga = a.lock().unwrap();
```

## 正确写法

- 统一加锁顺序：所有线程按同一顺序（如地址排序）获取锁；
- 用 `try_lock()` 代替 `lock()`，拿不到就放弃并重试；
- 用 `recv_timeout` / `join` 超时检测死锁；
- 计数器用 `AtomicUsize`，避免锁本身成为瓶颈。

本例第 1 段故意制造死锁，用 `recv_timeout(100ms)` 检测；第 3 段按同一顺序加锁，正常完成。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/B04-race-and-deadlock
cargo run --quiet
```

预期：`检测到死锁`、`期望 4000，实际 4000`、`未死锁`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 数据竞争 | 运行时 `-race` 检测 | 无（单线程） | GIL 不防读改写 | 无内建 | **编译期禁止** |
| 死锁 | `select + timeout` | `Promise.race` | `join(timeout)` | `tryLock(timeout)` | `try_lock` |
| 原子 | `sync/atomic` | 无 | 无 | `AtomicInteger` | `AtomicUsize` |

Rust 的 `Mutex` 和 `Arc` 一起用才能跨线程共享；`MutexGuard` 在作用域结束时自动释放，不会忘记 unlock。但「忘记统一顺序」编译器管不了，死锁依然要自己防。
