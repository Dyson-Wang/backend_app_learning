# F02 重试风暴

## 问题复现

下游开始超时，100 个客户端同时重试 3 次——本来 100 QPS 的流量瞬间变成 400 QPS，打在一个已经过载的服务上。它本来 30 秒能恢复，现在 5 分钟都起不来。

## 错误写法

```go
for i := 0; i < 3; i++ {
    if err := call(); err == nil { return nil }   // 立即重试，无退避无抖动
}
```

## 正确写法

```go
for i := 0; i < maxAttempts; i++ {
    if err := call(ctx); err == nil { return nil }
    sleep(backoff(i) + rand.Jitter())      // 指数退避 + 抖动
}
```

三件套缺一不可：

| 手段 | 作用 |
|---|---|
| 上限 | 别无限重试 |
| 退避 | 给下游恢复时间 |
| 抖动 | 避免所有客户端同时重试 |
| 熔断 | 下游明确挂了就别再试 |

## 怎么验证

```bash
go run .
```

预期：`400 QPS`（4 倍放大）→ 退避后 `133 QPS` → 熔断后 `0 QPS`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 退避库 | `backoff` | `p-retry` | `tenacity` | Resilience4j | `backoff` crate |
| 熔断 | `sony/gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |
| 常见坑 | 重试非幂等写 | 同 | 同 | 同 | 同 |
