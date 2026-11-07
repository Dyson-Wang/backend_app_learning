# A01 · 限流：令牌桶

## 问题复现

突发流量打进来，应用自身没有上限，数据库连接池先被打满，然后所有接口一起变慢。限流是把「超出容量的请求」在入口处快速拒绝，而不是让它们排队拖垮整个系统。

本例用 `httptest` 起服务，20 个并发请求打一个 `burst=5, rate=1/s` 的桶，观察 200 与 429 的分布，以及 1.1 秒后令牌补充能否重新放行。

## 错误写法

**错法 1 · 每个请求各自 new 一个桶**：限流形同虚设，桶必须是所有请求共享的同一个实例。

**错法 2 · 先判断再扣减但没加锁**：

```go
if b.tokens >= 1 {   // 读
    b.tokens--       // 写，和读不在同一个临界区
    return true
}
```

10 个 goroutine 可能同时通过判断，桶被扣成负数，一次超发十几个请求。用 `go run -race` 可以抓到。

**错法 3 · 固定窗口计数**：`INCR key` + `EXPIRE 60` 在窗口交界处会放行两倍流量（第 59 秒一波、第 61 秒一波）。改用滑动窗口或令牌桶。

**错法 4 · 单机限流当成全局限流**：4 个实例各限 100 QPS，实际是 400 QPS。分布式限流要把计数放到 Redis（Lua 保证原子），或由网关统一做。

## 正确写法

`Allow()` 把「按经过时间补充令牌 → 夹紧到 burst → 判断 → 扣减」整段放在 `sync.Mutex` 里，保证判断与扣减的原子性。用浮点记录令牌数，避免定时器 goroutine，也不需要后台任务。

配套细节：

- 返回 429 时带 `Retry-After`，让客户端知道该等多久
- 限流中间件放在最外层，早于鉴权、早于业务逻辑，被拒的请求不该消耗 DB 连接
- 按**维度**分桶：全局 / 用户 / IP / 接口，通常要叠几层
- 桶参数（rate、burst）走配置中心，别硬编码
- 生产里单机桶 + 网关全局限流一起用：前者兜底，后者精确

## 怎么验证

```
cd backend-roadmap/01-go
go run ./advanced/A01-rate-limiting
# 期望：200:5 429:15，1.1s 后 refill -> 200，selftest PASS

# 起真实服务手动压
go run ./advanced/A01-rate-limiting -serve
for i in $(seq 1 10); do curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8081/api; done

# 竞态检测
go run -race ./advanced/A01-rate-limiting
```

相关底座概念见 `../../../00-foundations/`。做完把差异结论补到 `../../../99-cross-language/scenario-matrix.md`。
