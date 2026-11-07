# C01 连接池与资源耗尽

## 问题复现

池只有 3 个连接，10 个并发请求各占 120ms。不配等待超时，7 个请求无限排队；下游一慢，队列爆炸，整个进程被拖死。

## 错误写法

```java
Connection c = pool.getConnection();     // 没有超时，永久阻塞
```

## 正确写法

```java
if (!sem.tryAcquire(40, TimeUnit.MILLISECONDS)) {
    throw new ServiceUnavailableException("pool exhausted");
}
```

各语言配置：

| 语言 | 参数 |
|---|---|
| Go | `sql.DB` + 请求级 `context` 超时 |
| Node.js | `pg.Pool({ connectionTimeoutMillis })` |
| Python | `asyncpg.create_pool(timeout=)` |
| Java | HikariCP `connectionTimeout` |
| Rust | `PgPoolOptions::acquire_timeout` |

## 怎么验证

```bash
javac Main.java && java Main
```

预期：第一种 `ok=10 timeout=0 耗时=400ms 左右`；第二种 `ok=3 timeout=7 耗时=120ms 左右`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 池实现 | `database/sql` | `pg.Pool` | `asyncpg` | HikariCP | `sqlx::Pool` |
| 等待超时 | `context` | `connectionTimeoutMillis` | `timeout` | `connectionTimeout` | `acquire_timeout` |
| 泄漏检测 | `ConnMaxLifetime` | 无内建 | `pool.recycle` | `leakDetectionThreshold` | 无内建 |
