# A02 · Rust · 超时、重试、指数退避与抖动

## 问题复现

下游偶发 503，调用方立刻重试，重试又失败——多个实例同时这么做，下游从「偶发失败」变成「持续过载」。把 400 参数错误也重试则纯属浪费。

## 错误写法

```rust
for _ in 0..3 {
    if let Ok(v) = call() { return Ok(v); } // 无差别重试
}
```

没有退避、没有可重试判定、没有上限。

## 正确写法

- `DownstreamError` 带 `retryable`，用 `Result` 表达；
- 指数退避 `base * 2^(attempt-1)`，本例 10 / 20 / 40ms；
- 用枚举 `Outcome` 区分「成功 / 永久失败 / 耗尽」，调用方必须 match 所有分支；
- 生产用 `backoff` crate 或 `tokio-retry`，并加抖动。

Rust 的 `Result` 让「错误必须被处理」成为编译期约束，枚举让「三种结局」不会漏。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/A02-timeout-retry-backoff
cargo run --quiet
```

预期：case 1 打印 3 行重试后 `Success { attempt: 4, ... }`；case 2 `Permanent { attempt: 1, ... }`；case 3 `Exhausted { attempt: 3, ... }`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 超时 | `context.WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 退避库 | `cenkalti/backoff` | `p-retry` | `tenacity` | Resilience4j | `backoff` |
| 错误表达 | `error` 返回值 | throw / reject | 异常 | 受检 / 非受检异常 | `Result<T, E>` |

Rust 的 `?` 运算符让错误传播最短，但重试逻辑必须显式写；Go 用 `errors.Is` 判断哨兵错误；Python/Java 用异常类型区分。
