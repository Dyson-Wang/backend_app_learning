# A01 · Node.js · 限流：令牌桶 / 滑动窗口

## 问题复现

接口没有限流，单个客户端一秒打几千次，CPU 被打满，正常用户全部超时。固定窗口计数器又有边界问题：窗口切换的瞬间可以打进 2 倍流量。

## 错误写法

```js
if (counter < 100) counter += 1; // 固定窗口，边界可翻倍
```

## 正确写法

令牌桶：桶里最多 `burst` 个令牌，按 `rate/s` 匀速补充，取到令牌才放行。

- 允许突发：空闲一段时间后可以一次性用掉 burst；
- 长期速率被 rate 限制；
- 用注入的 `now()` 推进逻辑时钟，输出完全确定。

生产上单机用内存桶，多实例用 Redis + Lua（`INCR` + `EXPIRE` 或令牌桶脚本）做分布式限流。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/A01-rate-limiting
node index.mjs
```

预期：t=0 前 5 次 `allow=true`、第 6 次 `false`；t=500 通过 1 次再拒绝；t=1500 通过 2 次再拒绝。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单机实现 | `golang.org/x/time/rate` | 手写 / `rate-limiter-flexible` | 手写 / `limits` | Guava `RateLimiter` | `governor` |
| 分布式 | Redis + Lua | Redis + Lua | Redis + Lua | Redisson | Redis + Lua |
| 时钟 | `time.Now` 注入 | `now` 回调 | `time.monotonic` | `LongSupplier` | `Instant` / 原子时钟 |

Node 单线程不需要锁，但限流状态要放在进程内单例，否则每个 worker 各限各的；多进程部署必须用 Redis 共享计数。
