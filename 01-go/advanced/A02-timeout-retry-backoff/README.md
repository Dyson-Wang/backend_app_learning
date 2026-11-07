# A02 · 超时、重试、指数退避与抖动

## 问题复现

下游服务变慢，上游没有超时，请求越堆越多直到连接池耗尽；加了重试却用紧循环，把一次小抖动放大成流量洪峰，这就是 `F02 重试风暴`。

本例三个自测：瞬时故障重试到成功、永久错误（4xx）绝不重试、ctx 超时后立刻停止重试。

## 错误写法

**错法 1 · 没有超时**：`http.Get(url)` 用默认 `http.Client`，Timeout 为 0，下游卡住就永远等下去。

**错法 2 · 重试所有错误**：把 400/401/404 也重试 5 次，纯浪费，还可能重复触发非幂等操作。只对网络错误、超时、429、5xx 重试。

**错法 3 · 固定间隔 + 无抖动**：

```go
for i := 0; i < 5; i++ {
    if call() == nil { break }
    time.Sleep(100 * time.Millisecond)  // 所有实例同一时刻重试
}
```

1000 个客户端同时失败就会同时重试，形成周期性洪峰，下游刚恢复又被打挂。

**错法 4 · 不遵守取消**：客户端已经断开（`ctx.Done()`），重试循环还在继续烧资源。

**错法 5 · 重试非幂等的写操作**：POST 创建订单被重试两次，产生两笔订单。重试前先确认 `B02` 幂等已就位。

## 正确写法

`DoWithRetry` 四个要素：

1. **每次尝试前检查 `ctx.Err()`**，等待期间用 `select` 同时监听 `ctx.Done()` 和 timer
2. **用 `errors.Is` 识别永久错误**，立即返回，不浪费额度
3. **指数退避 + 全抖动**：`ceiling = base * 2^(attempt-1)`，实际 delay 在 `[0, ceiling]` 随机取值，把重试时刻打散
4. **限制最大次数**，最终用 `%w` 包装最后一个错误，保留可判断的错误链

还有两条生产经验：总重试预算要小于上游的超时（否则上游已经放弃，下游还在重试）；重试要配合熔断（`A03`），下游真挂了就别再打了。

跨语言对照点：Go `context` / Node `AbortSignal` / Python `asyncio.timeout` / Java `orTimeout` / Rust `tokio::time::timeout`——同一条原则的五种落地形态。

## 怎么验证

```
cd backend-roadmap/01-go
go run ./advanced/A02-timeout-retry-backoff
# 期望：case1 第 3 次成功；case2 calls=1；case3 返回 context deadline exceeded；selftest PASS
```

相关底座概念见 `../../../00-foundations/`。做完把差异结论补到 `../../../99-cross-language/scenario-matrix.md`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 超时载体 | `context.WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 退避库 | `backoff` | `p-retry` | `tenacity` | Resilience4j | `backoff` crate |
| 永久错误识别 | `errors.Is` + 哨兵 | 自定义 error 类 | 异常类型 | 异常子类 | enum |
