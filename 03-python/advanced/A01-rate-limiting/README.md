# A01 · Python · 限流：令牌桶 / 滑动窗口

## 问题复现

接口不限流，单客户端一秒打几千次，CPU 打满，正常用户全部超时。固定窗口计数器还有边界问题：窗口切换瞬间能打进 2 倍流量。

## 错误写法

```python
if counter < 100:
    counter += 1  # 固定窗口，边界可翻倍
```

## 正确写法

令牌桶：桶里最多 `burst` 个令牌，按 `rate/s` 匀速补充，取到才放行。

- 允许突发，长期速率受控；
- 时钟通过 `clock` 可调用对象注入，测试用 `state["t"]` 推进，输出确定；
- 生产单机用内存桶，多实例用 Redis + Lua 做分布式限流。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/A01-rate-limiting
python3 main.py
```

预期：t=0 前 5 次 `allow=True`、第 6 次 `False`；t=0.5 通过 1 次再拒绝；t=1.5 通过 2 次再拒绝。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单机库 | `x/time/rate` | `rate-limiter-flexible` | `limits` | Guava `RateLimiter` | `governor` |
| 分布式 | Redis + Lua | Redis + Lua | Redis + Lua | Redisson | Redis + Lua |
| 时钟 | `time.Now` | `now` 回调 | 可调用对象 | `LongSupplier` | `Instant` |

Python 的 `time.monotonic()` 才是限流该用的时钟（不受系统时间调整影响）；`time.time()` 被 NTP 回拨时会导致令牌计算错乱。
