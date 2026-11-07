# C01 连接池与资源耗尽

## 问题复现

连接池只有 3 个连接，10 个并发请求各占 120ms。不配等待超时时，7 个请求排队等，全部成功但耗时 400ms+；一旦下游变慢，队列会无限增长，最终拖垮整个进程。

## 错误写法

```python
conn = pool.acquire()      # 无超时，永久阻塞
```

## 正确写法

等待必须有上限，超时立刻失败并返回 503，把压力还给调用方：

```python
if not pool.acquire(wait_ms=40):
    raise ServiceUnavailable("pool exhausted")
```

对应各语言配置：

| 语言 | 参数 |
|---|---|
| Go | `sql.DB.SetConnMaxLifetime` + 请求级 `context` 超时 |
| Node.js | `pg.Pool({ connectionTimeoutMillis })` |
| Python | `asyncpg.create_pool(timeout=...)` |
| Java | HikariCP `connectionTimeout` |
| Rust | `sqlx::postgres::PgPoolOptions::acquire_timeout` |

## 怎么验证

```bash
python3 main.py
```

预期：第一种 `ok=10 timeout=0 耗时=400ms 左右`；第二种 `ok=3 timeout=7 耗时=100ms 左右`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 池实现 | `database/sql` 内建 | `pg.Pool` | `asyncpg` / `QueuePool` | HikariCP | `sqlx::Pool` |
| 等待超时 | `context` 控制 | `connectionTimeoutMillis` | `timeout` | `connectionTimeout` | `acquire_timeout` |
| 泄漏检测 | `SetConnMaxLifetime` | 无内建 | `pool.recycle` | `leakDetectionThreshold` | 无内建 |
