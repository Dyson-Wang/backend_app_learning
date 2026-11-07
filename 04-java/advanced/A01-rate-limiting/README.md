# A01 · Java · 限流：令牌桶 / 滑动窗口

## 问题复现

接口不限流，单客户端一秒打几千次，Tomcat 线程池被占满，正常用户全部超时。固定窗口计数器还有边界问题：窗口切换瞬间能打进 2 倍流量。

## 错误写法

```java
if (counter.get() < 100) counter.incrementAndGet(); // 固定窗口，边界可翻倍
```

## 正确写法

令牌桶：最多 `burst` 个令牌，按 `rate/s` 补充。

- 允许突发，长期速率受控；
- 用内部 `nowMs` 逻辑时钟推进，测试不依赖 `Thread.sleep`；
- 生产单机用 Guava `RateLimiter`，分布式用 Redisson `RRateLimiter`。

Guava 的 `RateLimiter.create(2.0)` 就是平滑预热令牌桶，`tryAcquire()` 非阻塞、`acquire()` 阻塞，语义要明确选。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/A01-rate-limiting
javac Main.java && java Main
```

预期：t=0 前 5 次 `allow=true`、第 6 次 `false`；t=500 通过 1 次再拒绝；t=1500 通过 2 次再拒绝。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单机库 | `x/time/rate` | `rate-limiter-flexible` | `limits` | Guava `RateLimiter` | `governor` |
| 分布式 | Redis + Lua | Redis + Lua | Redis + Lua | Redisson | Redis + Lua |
| 阻塞语义 | `Wait(ctx)` | 无 | `acquire()` | `acquire()` | `until_ready()` |

Guava 的 `RateLimiter` 是不可变的速率配置，改速率要重建；Redisson 的分布式限流基于 Redis + Lua，多实例共享同一个桶。
