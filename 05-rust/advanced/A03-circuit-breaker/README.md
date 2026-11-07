# A03 · Rust · 熔断与降级

## 问题复现

下游超时后，如果调用方持续请求，会占满连接池与 tokio worker，最终整条链路雪崩。熔断器把「下游不可用」变成一次本地状态判断，不再产生网络 I/O。

## 错误写法

```rust
for _ in 0..3 {
    if let Ok(r) = client.call().await { return Ok(r); } // 失败就重试
}
```

没有状态、没有计数、没有冷却，重试会放大流量。

## 正确写法

`Breaker` 用 `Mutex<Inner>` 保护 `closed / open / half-open`：

- `open` 时 `allow()` 返回 `Err`，不发起请求；
- 冷却结束转 `half-open`，`probe_in_flight` 保证只有一个探测；
- 探测成功回 `closed`，失败立即回 `open`。

时钟用 `Arc<AtomicU64>`，测试里 `advance(200)` 推进，不依赖真实 `Instant`。Rust 的优势是「状态只能通过 `Mutex` 访问」，忘记加锁编译不过。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/A03-circuit-breaker
cargo run --quiet
```

预期：`state=closed` → 5 次失败后 `state=open` → `rejected: circuit open: fast fail` → `advance(200)` 后 `state=half-open` → 探测成功回 `closed`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 状态保护 | `sync.Mutex` | 单线程 + 无 await 临界区 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 时钟注入 | 函数字段 | `now` 回调 | 可调用对象 | 内部 `long now` | `Arc<AtomicU64>` |
| 生产库 | `sony/gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |

Rust 在 async 场景下不要跨 `.await` 持有 `std::sync::Mutex`（会阻塞 executor），应换 `tokio::sync::Mutex`；Go 的 `gobreaker` 用 `sync.Mutex` + 函数时钟；Java 的 Resilience4j 是组合式装饰器。
