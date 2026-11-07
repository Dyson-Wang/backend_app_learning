# B04 · Java · 竞态、死锁与调试

## 问题复现

两个线程同时执行 `int v = counter; counter = v + 1`，都读到旧值，最终只加一次。死锁则是两个线程各持一把锁、反序获取对方的锁，表现为线程永久 `BLOCKED`。

## 错误写法

```java
int value = counter;
counter = value + 1; // 读改写不是原子的
```

## 正确写法

- 互斥：`synchronized` 或 `ReentrantLock`；
- 死锁：统一加锁顺序，或用 `tryLock(timeout)` 放弃重试；
- 检测：`CountDownLatch.await(timeout)`、`jstack`、JFR；
- daemon 线程不阻止 JVM 退出，适合演示死锁而不挂住终端。

`tryLock(50, MILLISECONDS)` 是 Java 比 `synchronized` 强的地方：`synchronized` 一旦阻塞就无法超时退出。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/B04-race-and-deadlock
javac Main.java && java Main
```

预期：`期望 2，实际 1`、`期望 4000，实际 4000`、`检测到死锁`、`ab2 放弃 d` / `ba2 放弃 c` 之类输出。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 竞态检测 | `-race` | async stack | `faulthandler` | JFR / `jstack` | 编译器 |
| 互斥 | `sync.Mutex` | Promise 链 | `threading.Lock` | `synchronized` / `ReentrantLock` | `Mutex<T>` |
| 超时加锁 | `TryLock` | 无 | `acquire(timeout)` | `tryLock(timeout)` | `try_lock` |
| 原子 | `sync/atomic` | 无 | 无 | `AtomicInteger` | `AtomicUsize` |

Java 的 `jstack` 能直接打印死锁线程和持有的锁，是排查线上死锁最快的手段；Rust 在编译期禁止数据竞争，但死锁仍可能发生。
