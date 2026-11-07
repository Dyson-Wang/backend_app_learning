# C01 · Go · 连接池与资源耗尽

## 问题复现

池子只有 3 个连接，10 个请求各要 100ms。没有等待超时的话，后 7 个会一直排队，前端等到几秒后才超时，而这时连接还占着，故障被放大。

## 错误写法

```go
conn := pool.Get() // 没有超时，拿不到就一直阻塞
defer pool.Put(conn)
```

## 正确写法

- 池大小按下游容量设（数据库 `max_connections` / 实例数），不是越大越好；
- 获取连接必须带超时（`AcquireTimeout` / `ConnMaxLifetime`）；
- 用完立刻归还，`defer` 是最稳的写法；
- 监控「等待时长」和「等待队列长度」，这两个指标先于错误率报警。

本例池大小 3、10 个请求、每个 100ms：无超时全部成功但耗时约 400ms；等待超时 50ms 时 3 个成功、7 个快速失败、耗时约 100ms。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/C01-connection-pool
go run .
```

预期：第 1 段 `ok=10 timeout=0`，耗时约 400ms；第 2 段 `ok=3 timeout=7`，耗时约 100ms。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 连接池 | `database/sql` 内建 | `pg.Pool` | SQLAlchemy `pool_size` | HikariCP | `deadpool` / `sqlx::Pool` |
| 等待超时 | `ConnMaxIdleTime` / ctx | `connectionTimeoutMillis` | `pool_timeout` | `connectionTimeout` | `acquire_timeout` |
| 池上限 | `SetMaxOpenConns` | `max` | `max_overflow` | `maximumPoolSize` | `max_size` |

Go 的 `database/sql` 自带池，不需要第三方；Node 的 `pg.Pool` 默认无限等待，必须显式配 `connectionTimeoutMillis`；HikariCP 的默认值最保守，但 `maximumPoolSize` 设太大反而更慢。
