# C01 · Node.js · 连接池与资源耗尽

## 问题复现

池子 3 个槽位，10 个请求各 100ms。没有等待超时，后 7 个一直排队，前端等到几秒后超时，而槽位还占着——故障被放大。

## 错误写法

```js
await pool.acquire(); // 没有超时，拿不到就一直等
```

## 正确写法

- 池大小按下游容量设；
- `acquire` 必须带超时，超时直接失败（返回 503），不要无限排队；
- 超时后要把 waiter 从队列里摘掉，否则 `release` 会唤醒一个已经失败的 Promise；
- 监控「等待时长」和「队列长度」。

本例用信号量模拟：无超时 `ok=10`、约 400ms；超时 50ms 时 `ok=3 timeout=7`、约 100ms。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/C01-connection-pool
node index.mjs
```

预期：第 1 段 `ok=10 timeout=0` 约 400ms；第 2 段 `ok=3 timeout=7` 约 100ms。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 连接池 | `database/sql` | `pg.Pool` / `undici.Pool` | SQLAlchemy pool | HikariCP | `deadpool` |
| 等待超时 | ctx | `connectionTimeoutMillis` | `pool_timeout` | `connectionTimeout` | `acquire_timeout` |
| 并发上限 | 信号量 / worker pool | `p-limit` | `Semaphore` | `Semaphore` | `tokio::sync::Semaphore` |

Node 的 `pg.Pool` 默认 `connectionTimeoutMillis=0`（无限等待），这是线上最常见的坑之一；`undici` 的 Pool 默认有 `connections` 上限但仍需配超时。
