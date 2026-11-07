# A03 · Python · 熔断与降级

## 问题复现

下游超时后，调用方若继续请求，线程池会被占满、请求排队、上游超时被拖长。Python 的 `requests` 默认没有超时，一次下游故障就可能把整个 worker 卡死。熔断器把「下游不可用」变成一次本地状态判断。

## 错误写法

```python
for _ in range(3):
    try:
        return requests.get(url, timeout=1)
    except Exception:
        continue  # 下游已经挂了，重试只是放大流量
```

没有状态、没有计数、没有冷却，重试风暴会把下游按死。

## 正确写法

`Breaker` 用 `threading.Lock` 保护三态：

- `closed`：放行，失败累加；
- 连续失败达到阈值 → `open`，`allow()` 直接返回 `False`，**不发起网络调用**；
- `open` 冷却结束 → `half-open`，`probe_in_flight` 保证只有一个探测；
- 探测成功 → `closed`，失败 → 立即回 `open`。

时钟通过 `clock` 参数注入，测试用 `FakeClock` 推进，不依赖真实 sleep。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/A03-circuit-breaker
python3 main.py
```

预期：`state=closed` → 5 次失败后 `state=open` → `rejected: circuit open: fast fail` → 推进 0.2s 后 `state=half-open` → 探测成功回 `closed`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 状态保护 | `sync.Mutex` | 单线程 + 无 await 临界区 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 时钟注入 | 函数字段 | `now` 回调 | 可调用对象 | `LongSupplier` | `Arc<AtomicU64>` |
| 常见库 | `sony/gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |

Python 有 GIL，但 `time.sleep` / I/O 期间会释放，所以状态判断仍需要锁；Java 的 `synchronized` 最省事；Rust 的所有权让「忘记加锁」变成编译错误。
