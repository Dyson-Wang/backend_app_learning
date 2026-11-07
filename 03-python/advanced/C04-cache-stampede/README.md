# C04 · Python · 缓存穿透 / 击穿 / 雪崩

## 问题复现

热点 key 过期瞬间，50 个线程同时发现缓存为空，全部回源数据库。Python 有 GIL，但 `time.sleep` 和 I/O 会释放 GIL，所以并发窗口真实存在——本例用 20ms 的假数据库打出来。

## 错误写法

```python
value = cache.get(key)
if value is None:
    value = db.query(key)  # 50 个线程都会到这里
    cache[key] = value
```

「查 → 回源 → 写」不是原子的，线程越多回源越多。

## 正确写法

按 key 加锁 + 双重检查：

1. 先查缓存，命中直接返回；
2. 未命中取该 key 专属的 `threading.Lock`，不同 key 不互相阻塞；
3. 拿到锁后**再查一次**，等锁期间别人可能已经回源完成；
4. 仍未命中才回源并写回。

`key_locks.setdefault(key, threading.Lock())` 在全局锁内完成，保证同一个 key 只有一个锁对象。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/C04-cache-stampede
python3 main.py
```

预期：`无保护 ... 回源 50 次`，`按 key 加锁 + 双重检查 ... 回源 1 次`，耗时都在 20ms 量级。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单飞实现 | per-key mutex / `singleflight` | in-flight Promise Map | per-key `Lock` | `computeIfAbsent` | `Mutex` + `HashMap` |
| 并发单元 | goroutine | Promise | 线程 / asyncio | 线程池 | 线程 / tokio |
| 易错点 | 锁粒度 | 忘记 finally 清理 | 锁内做 I/O | 慢操作放进 `computeIfAbsent` | 跨 await 持 std Mutex |

Python 的 `asyncio` 版本要用 `asyncio.Lock` 而不是 `threading.Lock`；Java 的 `computeIfAbsent` 一行解决但要保证映射函数不慢；Node 的 Promise 去重最简洁；Rust 在 async 里要用 `tokio::sync::Mutex`。
