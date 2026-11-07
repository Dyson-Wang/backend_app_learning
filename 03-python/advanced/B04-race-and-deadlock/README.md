# B04 · Python · 竞态、死锁与调试

## 问题复现

两个线程同时执行 `v = counter; counter = v + 1`，都读到旧值，最终只加一次。GIL 只保证单条字节码原子，`counter += 1` 是「读-加-写」三条字节码，仍会丢更新。

死锁则是两个线程各持一把锁、互相等对方，表现为线程永久阻塞。

## 错误写法

```python
value = counter
counter = value + 1  # 读和写之间没有保护
```

## 正确写法

- 互斥：`with lock:` 包住读改写；
- 死锁：两把锁全局统一顺序，或用 `lock.acquire(timeout=...)`；
- 检测：`threading.Event.wait(timeout)` 或 `Thread.join(timeout)`；
- daemon 线程不会阻止进程退出，适合演示死锁而不挂住终端。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/B04-race-and-deadlock
python3 main.py
```

预期：`期望 2，实际 1`、`期望 4000，实际 4000`、`检测到死锁`、两行 `thread-1/2` 完成。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 竞态检测 | `-race` | async stack | `faulthandler` | JFR / thread dump | 编译器 |
| 互斥 | `sync.Mutex` | Promise 链 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 原子 | `sync/atomic` | 无 | 无（需锁） | `AtomicInteger` | `AtomicUsize` |
| 死锁检测 | `select + timeout` | `Promise.race` | `join(timeout)` | `tryLock(timeout)` | `recv_timeout` |

Python 有 GIL 反而容易让人误以为「不用加锁」，这是最常见的坑；真正需要并行计算时用 `multiprocessing`，GIL 就不存在了。
