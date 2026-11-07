# A02 · Python · 超时、重试、指数退避与抖动

## 问题复现

下游偶发 503，调用方立刻重试，重试又失败——多个实例同时这么做，下游从「偶发失败」变成「持续过载」。另一种情况是把 400 参数错误也重试，纯属浪费。

## 错误写法

```python
for _ in range(3):
    try:
        return call()
    except Exception:  # 无差别重试
        pass
```

没有退避、没有可重试判定、没有上限。

## 正确写法

- 错误带 `retryable`：5xx / 超时 / 连接错误可重试，4xx 不重试；
- 指数退避 `base * 2^(attempt-1)`，本例 10 / 20 / 40ms；
- 最大尝试次数，耗尽返回最后一次错误；
- 生产必须加抖动（`delay * random.uniform(0.5, 1.0)`），否则所有客户端同时重试形成风暴。

本例去掉抖动保证输出可复现，README 里说明。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/A02-timeout-retry-backoff
python3 main.py
```

预期：case 1 打印 3 行重试后成功（`attempt: 4`）；case 2 立即返回 `permanent: True`；case 3 耗尽 3 次返回 `permanent: False`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 超时 | `context.WithTimeout` | `AbortSignal.timeout` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 退避库 | `cenkalti/backoff` | `p-retry` | `tenacity` | Resilience4j | `backoff` |
| 取消 | context | AbortSignal | `Task.cancel()` | `Future.cancel` | drop future |

Python 的 `tenacity` 用装饰器最省事（`@retry(stop=stop_after_attempt(3), wait=wait_exponential())`）；asyncio 场景要用 `asyncio.sleep` 而不是 `time.sleep`，否则会阻塞整个事件循环。
