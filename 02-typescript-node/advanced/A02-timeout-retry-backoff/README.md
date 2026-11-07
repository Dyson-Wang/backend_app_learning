# A02 · Node.js · 超时、重试、指数退避与抖动

## 问题复现

下游偶发 503，调用方立刻重试，重试又失败，再重试——三个实例同时这么做，下游从「偶发失败」变成「持续过载」。另一种情况是把 400 参数错误也重试，纯属浪费。

## 错误写法

```js
for (let i = 0; i < 3; i += 1) {
  try { return await call(); } catch { /* 无差别重试 */ }
}
```

没有退避、没有可重试判定、没有上限。

## 正确写法

- 错误带 `retryable` 标记：5xx / 超时 / 连接错误可重试，4xx 不重试；
- 指数退避：`base * 2^(attempt-1)`，本例 10 / 20 / 40ms；
- 有最大尝试次数，耗尽后返回最后一次错误；
- 生产必须加**抖动**（`delay * (0.5 + random()*0.5)`），否则所有客户端会在同一时刻重试，形成重试风暴。

本例为了输出可复现，去掉了抖动并在 README 里说明。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/A02-timeout-retry-backoff
node index.mjs
```

预期：case 1 打印 3 行重试后成功（`attempt: 4`）；case 2 立即返回 `permanent: true`；case 3 耗尽 3 次返回 `permanent: false`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 超时 | `context.WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 退避库 | `cenkalti/backoff` | `p-retry` | `tenacity` | Resilience4j / Failsafe | `backoff` |
| 取消传播 | context 一路传 | AbortSignal | Task.cancel | Future.cancel | drop future |

Go 的 context 让超时和取消沿着调用链自动传播；Node 的 `AbortSignal` 要显式传给 fetch；Python 的 `tenacity` 用装饰器最省事；Java 的 Resilience4j 把重试和熔断组合；Rust 的取消是 drop，超时靠 `select!`。
