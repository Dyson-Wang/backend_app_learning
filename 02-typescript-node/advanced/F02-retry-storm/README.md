# F02 重试风暴

## 问题复现

下游开始超时，100 个客户端同时重试 3 次，100 QPS 瞬间变 400 QPS，打在一个已过载的服务上。它本来 30 秒能恢复，现在 5 分钟起不来。

## 错误写法

```js
for (let i = 0; i < 3; i++) {
  if (await call()) return;      // 立即重试，无退避无抖动
}
```

## 正确写法

```js
for (let i = 0; i < maxAttempts; i++) {
  if (await call()) return;
  await sleep(2 ** i * 100 + Math.random() * 100);   // 退避 + 抖动
}
```

| 手段 | 作用 |
|---|---|
| 上限 | 别无限重试 |
| 退避 | 给下游恢复时间 |
| 抖动 | 避免同时重试 |
| 熔断 | 明确挂了就别再试 |

## 怎么验证

```bash
node index.mjs
```

预期：`400 QPS` → 退避后 `133 QPS` → 熔断后 `0 QPS`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 退避库 | `backoff` | `p-retry` | `tenacity` | Resilience4j | `backoff` crate |
| 熔断 | `gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |
| 陷阱 | 重试非幂等写 | 同 | 同 | 同 | 同 |
