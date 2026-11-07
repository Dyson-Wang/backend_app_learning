# C04 · Java · 缓存穿透 / 击穿 / 雪崩

## 问题复现

热点 key 过期瞬间，50 个线程同时发现缓存为空，全部打到数据库。数据库连接池只有几十个，立刻被打满，接口 P99 从 20ms 涨到几秒。本例用 20ms 的假数据库 + 回源计数器复现。

## 错误写法

```java
String v = cache.get(key);
if (v == null) {
    v = db.query(key); // 所有并发线程都到这里
    cache.put(key, v);
}
```

`ConcurrentHashMap` 只保证单个操作线程安全，不保证「查-回源-写」这个组合的原子性。

## 正确写法

按 key 加锁 + 双重检查：

1. `data.get` 命中直接返回；
2. 未命中用 `keyLocks.computeIfAbsent` 取该 key 的锁对象；
3. `synchronized(lock)` 后**再查一次** `data`；
4. 仍未命中才 `db.query` 并写回。

`computeIfAbsent` 本身是原子的，保证同一个 key 只创建一个锁对象；不同 key 用不同锁，不会互相阻塞。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/C04-cache-stampede
javac Main.java && java Main
```

预期：`无保护 ... 回源 50 次`，`按 key 加锁 + 双重检查 ... 回源 1 次`，耗时都在 20ms 量级。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单飞实现 | per-key mutex | in-flight Promise Map | per-key `Lock` | `computeIfAbsent` + `synchronized` | `Mutex` + `HashMap` |
| 原子创建锁 | `sync.Map` | 单线程无需 | `setdefault` | `computeIfAbsent` | `entry().or_insert_with` |
| 易错点 | 锁粒度 | 忘记 finally | 锁内 I/O | `computeIfAbsent` 里做慢操作 | 跨 await 持 std Mutex |

Java 的 `computeIfAbsent` 是标准答案，但映射函数里不能做耗时操作（会持有桶锁）；Go 的 `singleflight` 更省事；Node 用 Promise 去重最自然；Rust 的 `entry` API 等价于 `computeIfAbsent`。
