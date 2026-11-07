# A03 · Go · 熔断与降级

## 问题复现

下游服务超时后，调用方如果不加保护，会持续把请求打过去：连接池被占满、goroutine 堆积、上游超时被拖长，最终整条链路雪崩。熔断器的作用是「下游明显不可用时，直接快速失败」，把恢复时间留给下游。

## 错误写法

```go
// 每次都直接调用，失败就重试
resp, err := client.Do(req)
if err != nil {
    return retry(req) // 下游已经挂了，重试只会放大流量
}
```

没有状态机、没有失败计数、没有冷却窗口，重试风暴会直接把下游按死。

## 正确写法

`Breaker` 维护 `closed / open / half-open` 三态：

- `closed`：正常放行，失败累加；
- 连续失败达到阈值 → `open`，之后所有请求**不发网络调用**，直接返回错误；
- `open` 持续到冷却时间结束 → `half-open`，只放一个探测请求；
- 探测成功 → `closed`，失败 → 立刻回 `open`。

关键是 `Allow()` 在 `open` 状态下不碰下游，这才是「降级」。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/A03-circuit-breaker
go run .
```

预期输出依次出现 `state=closed`、5 次失败后 `state=open`、`rejected: circuit open`、推进 200ms 后 `state=half-open`、探测成功回到 `closed`。本例用可注入时钟，不依赖真实 sleep，输出可复现。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 状态保护 | `sync.Mutex` | 单线程无锁 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 时钟注入 | 函数字段 | 传 `now` 函数 | 传 `clock` 可调用 | `Clock` 接口 | `fn() -> Instant` |
| 典型库 | `sony/gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |

Node 单线程不需要锁，但 `await` 期间仍可能并发进入，所以状态判断要和「是否已有探测在飞」一起做；Java 的 Resilience4j 把熔断、限流、重试做成了组合器；Rust 的 `Mutex` 让状态机天然线程安全，但要小心别在持锁时做 I/O。
