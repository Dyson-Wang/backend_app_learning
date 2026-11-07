# A02 · Java · 超时、重试、指数退避与抖动

## 问题复现

下游偶发 503，调用方立刻重试，重试又失败——多个实例同时这么做，下游从「偶发失败」变成「持续过载」。把 400 参数错误也重试则纯属浪费。

## 错误写法

```java
for (int i = 0; i < 3; i++) {
    try { return call(); } catch (Exception e) { /* 无差别重试 */ }
}
```

没有退避、没有可重试判定、没有上限。

## 正确写法

- 自定义 `DownstreamException` 带 `retryable`；
- 指数退避 `base * 2^(attempt-1)`，本例 10 / 20 / 40ms；
- 最大尝试次数，耗尽返回最后一次错误；
- 生产用 Resilience4j `Retry` 或 Failsafe，并加抖动。

Java 的受检异常让「必须处理」变成编译期约束，但也导致 lambda 里抛异常需要自定义函数式接口（本例的 `Attempt`）。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/A02-timeout-retry-backoff
javac Main.java && java Main
```

预期：case 1 打印 3 行重试后成功（`attempt=4`）；case 2 立即 `permanent=true`；case 3 耗尽 3 次 `permanent=false`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 超时 | `context.WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 退避库 | `cenkalti/backoff` | `p-retry` | `tenacity` | Resilience4j / Failsafe | `backoff` |
| 异常模型 | 返回值 | throw / reject | 异常 | 受检 / 非受检 | `Result<T, E>` |

Java 的受检异常强制调用方处理，但也让「重试策略」必须显式建模；Go 用返回值 + `errors.Is` 判断可重试；Rust 用 `Result` + 自定义错误枚举。
