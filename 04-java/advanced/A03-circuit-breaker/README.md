# A03 · Java · 熔断与降级

## 问题复现

下游超时后，Tomcat 线程池会被慢请求占满，新请求排队，最终整个服务不可用。熔断器把「下游不可用」变成一次本地状态判断，避免继续消耗线程。

## 错误写法

```java
for (int i = 0; i < 3; i++) {
    try { return client.call(); } catch (Exception e) { /* 重试 */ }
}
```

没有状态机、没有失败计数、没有冷却，重试会把下游压死，同时把当前线程也拖住。

## 正确写法

`Breaker` 用 `synchronized` 保护 `closed / open / half-open`：

- `open` 时 `allow()` 直接返回错误字符串，**不发远程调用**；
- 冷却结束转 `half-open`，`probeInFlight` 保证只有一个探测；
- 探测成功 `success()` 清零并回 `closed`，失败 `failure()` 立即回 `open`。

`now` 是内部逻辑时钟，`advance(200)` 推进它，测试不依赖 `Thread.sleep`。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/A03-circuit-breaker
javac Main.java && java Main
```

预期：`state=closed` → 5 次失败后 `state=open` → `rejected: circuit open: fast fail` → `advance(200)` 后 `state=half-open` → 探测成功回 `closed`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 状态保护 | `sync.Mutex` | 单线程 + 无 await 临界区 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 时钟注入 | 函数字段 | `now` 回调 | 可调用对象 | 内部 `long now` | `Arc<AtomicU64>` |
| 生产库 | `sony/gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |

Java 的 `synchronized` 最直接，但要注意别在同步块里做远程调用；Resilience4j 把熔断、限流、重试、舱壁组合成装饰器链；Go 的 `gobreaker` 用 `sync.Mutex` + 函数时钟，思路和本例一致。
