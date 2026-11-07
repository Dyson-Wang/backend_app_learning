# C04 · Node.js · 缓存穿透 / 击穿 / 雪崩

## 问题复现

Node 是单线程，很多人以为「不会并发」。但 `await db.query()` 会把执行权交回事件循环，50 个请求会在 `await` 之前全部通过「缓存为空」的检查，然后一起回源。本例用 20ms 的假数据库把这件事打出来。

## 错误写法

```js
if (cache.has(key)) return cache.get(key);
const v = await db.query(key); // 50 个请求都会到这里
cache.set(key, v);
```

`has` 和 `set` 之间隔着 `await`，这个间隙就是并发窗口。

## 正确写法

用 **in-flight Promise Map** 做 singleflight：

1. 命中 `data` 直接返回；
2. 命中 `inflight` 就 `return` 同一个 Promise，后来者不再回源；
3. 未命中则创建回源 Promise，存进 `inflight`，并在 `finally` 里删除；
4. 所有等待者共享同一次查询结果。

关键是 `finally` 清理：失败也要删，否则失败的 key 会被永久卡住。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/C04-cache-stampede
node index.mjs
```

预期：`无保护 ... 回源 50 次`，`in-flight 去重 ... 回源 1 次`。耗时都在 20~30ms 量级，说明并发窗口真实存在。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单飞实现 | `singleflight` / per-key mutex | in-flight Promise Map | per-key `Lock` | `computeIfAbsent` | `Mutex` + `HashMap` |
| 并发单元 | goroutine | Promise | 线程 / asyncio Task | 线程池 | 线程 / tokio task |
| 清理时机 | `defer` 解锁 | `.finally()` | `finally` | `computeIfAbsent` 自动 | `Drop` / guard |

Node 的 Promise 去重写法最短，但要记得 `finally`；Go 的 `singleflight` 是标准答案；Java 的 `computeIfAbsent` 原子但不适合慢操作；Rust 的 `Drop` guard 能保证异常路径也解锁。
