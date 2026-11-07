# F01 慢查询拖垮连接池

## 问题复现

连接池 2 个连接，10 个并发请求，正常每条 30ms。某天一条 SQL 走了全表扫描，跑 300ms——它占住 1 个连接，其余 9 个请求排队，第 6 个开始等超过 100ms 直接超时。接口错误率 50%，但数据库 CPU 并不高。

## 错误写法

```go
rows, err := db.Query("SELECT * FROM orders WHERE user_name = $1", name)  // 无索引
// 也没有 statement_timeout
```

## 正确写法

```go
ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
defer cancel()
db.QueryContext(ctx, ...)          // 请求级超时
```

```sql
SET statement_timeout = '100ms';   -- 兜住漏网的慢 SQL
CREATE INDEX idx_orders_user_name ON orders (user_name);
```

连接池还要配 `SetConnMaxLifetime`，避免连接被长期占用。

## 怎么验证

```bash
go run .
```

预期：正常 `成功 = 10，超时 = 0，总耗时 = 150ms`；混入慢查询 `成功 = 5，超时 = 5`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 请求级超时 | `QueryContext` | `statement_timeout` | `statement_timeout` | `@Transactional(timeout=)` | `PgPoolOptions` |
| 连接池上限 | `SetMaxOpenConns` | `max` | `pool.max_size` | `maximumPoolSize` | `max_connections` |
| 慢日志 | `log_min_duration_statement` | 同 | 同 | 同 | 同 |
