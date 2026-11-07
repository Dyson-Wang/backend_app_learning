# C04 · Go · 缓存穿透 / 击穿 / 雪崩

## 问题复现

热点 key 过期的一瞬间，50 个并发请求同时发现「缓存没有」，于是全部回源数据库。数据库连接池只有几十个，瞬间被占满，接口 P99 从 20ms 涨到几秒——这就是缓存击穿。本例用一个 20ms 的假数据库 + 回源计数器，把「50 并发 → 50 次回源」直接打出来。

## 错误写法

```go
v, ok := cache[key]
if !ok {
    v = db.Query(key) // 所有并发请求都走到这里
    cache[key] = v
}
```

「查缓存 → 回源 → 写缓存」三步不是原子的，并发越高回源次数越多，缓存反而成了放大器。

## 正确写法

按 key 粒度加锁 + 双重检查：

1. 先查缓存，命中直接返回；
2. 未命中时取该 key 专属的 `sync.Mutex`（不同 key 不互相阻塞）；
3. 拿到锁后**再查一次缓存**——等锁期间别人可能已经回源完成；
4. 仍未命中才回源，写回缓存后释放锁。

这样同一 key 的并发请求只会有一个真正打数据库，其余等锁后直接读缓存。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/C04-cache-stampede
go run .
```

预期：`无保护 ... 回源 50 次`，`互斥+双重检查 ... 回源 1 次`。耗时也会从 20ms 级降到接近 20ms（并发等待）。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单飞实现 | `sync.Mutex` per key / `singleflight` | in-flight Promise Map | `threading.Lock` per key | `ConcurrentHashMap.computeIfAbsent` | `Mutex` + `HashMap` |
| 天然优势 | goroutine 轻量 | 单线程，用 Promise 去重最简洁 | GIL 不解决 I/O 并发 | `computeIfAbsent` 原子 | 所有权保证安全 |
| 常见坑 | 锁粒度太粗 | 忘记 finally 清理 in-flight | 锁内做 I/O 阻塞其他线程 | `computeIfAbsent` 里做长耗时操作 | 持锁跨 await（tokio） |

Node 的 in-flight Promise 方案最优雅：把「正在回源的 Promise」存进 Map，后来者直接 `await` 同一个 Promise；Java 的 `computeIfAbsent` 一行搞定但要避免在映射函数里做慢操作；Rust 在 async 场景下不能跨 `await` 持 `std::sync::Mutex`，要用 `tokio::sync::Mutex`。
