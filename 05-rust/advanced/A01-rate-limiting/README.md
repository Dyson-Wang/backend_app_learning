# A01 · Rust · 限流：令牌桶 / 滑动窗口

## 问题复现

接口不限流，单客户端一秒打几千次，CPU 打满，正常用户全部超时。固定窗口计数器还有边界问题：窗口切换瞬间能打进 2 倍流量。

## 错误写法

```rust
if counter.load(Ordering::Relaxed) < 100 {
    counter.fetch_add(1, Ordering::Relaxed); // 固定窗口，边界可翻倍
}
```

## 正确写法

令牌桶：最多 `burst` 个令牌，按 `rate/s` 补充。

- 允许突发，长期速率受控；
- 时钟用 `Arc<AtomicU64>`，`advance(500)` 推进，测试不依赖真实 `Instant`；
- `allow(&mut self)` 需要可变借用，保证状态修改是独占的；
- 生产单机用 `governor`，分布式用 Redis + Lua。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/A01-rate-limiting
cargo run --quiet
```

预期：t=0 前 5 次 `allow=true`、第 6 次 `false`；t=500 通过 1 次再拒绝；t=1500 通过 2 次再拒绝。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单机库 | `x/time/rate` | `rate-limiter-flexible` | `limits` | Guava `RateLimiter` | `governor` |
| 分布式 | Redis + Lua | Redis + Lua | Redis + Lua | Redisson | Redis + Lua |
| 状态可变性 | 指针 + 锁 | 对象属性 | 实例属性 | 实例字段 | `&mut self` |

Rust 的 `&mut self` 让「同一个桶不能被两个线程同时改」变成编译期约束；要跨线程共享就用 `Mutex<TokenBucket>`，代价是每次 `allow` 都要加锁。
