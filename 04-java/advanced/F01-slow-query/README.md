# F01 慢查询拖垮连接池

## 问题复现

池 2 个连接，10 个并发，正常 30ms。一条全表扫描 SQL 跑 300ms，占住 1 个连接，第 6 个请求起等超 100ms 超时。错误率 50%，数据库 CPU 却不高。

## 错误写法

```java
jdbc.queryForList("SELECT * FROM orders WHERE user_name = ?", name);   // 无索引、无超时
```

## 正确写法

```java
@Transactional(timeout = 1)      // 秒
// 或 HikariCP connectionTimeout + PostgreSQL statement_timeout
```

```sql
CREATE INDEX idx_orders_user_name ON orders (user_name);
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：正常 `成功 = 10，超时 = 0，总耗时 = 150ms`；混入慢查询 `成功 = 5，超时 = 5`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 请求级超时 | `QueryContext` | `query({ timeout })` | `timeout=` | `@Transactional(timeout=)` | `acquire_timeout` |
| 池上限 | `SetMaxOpenConns` | `pg.Pool max` | `max_size` | `maximumPoolSize` | `max_connections` |
| 陷阱 | 忘关 rows | 忘 release | 忘 close | 连接泄漏 | 忘 drop |
