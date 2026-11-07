# F01 慢查询拖垮连接池

## 问题复现

池 2 个连接，10 个并发，正常 30ms。一条全表扫描 SQL 跑 300ms，占住 1 个连接，第 6 个请求起等超 100ms 超时。

## 错误写法

```rust
sqlx::query("SELECT * FROM orders WHERE user_name = $1")
    .bind(name).fetch_all(&pool).await?;      // 无索引、无超时
```

## 正确写法

```rust
tokio::time::timeout(Duration::from_millis(100),
    sqlx::query(...).fetch_all(&pool)).await??;
```

```sql
SET statement_timeout = '100ms';
CREATE INDEX idx_orders_user_name ON orders (user_name);
```

## 怎么验证

```bash
cargo run --quiet
```

预期：正常 `成功 = 10，超时 = 0，总耗时 = 150ms`；混入慢查询 `成功 = 5，超时 = 5`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 请求级超时 | `QueryContext` | `query({ timeout })` | `timeout=` | `@Transactional(timeout=)` | `tokio::time::timeout` |
| 池上限 | `SetMaxOpenConns` | `pg.Pool max` | `max_size` | `maximumPoolSize` | `max_connections` |
| 特点 | 显式 ctx | 配置项 | 配置项 | 注解 | **超时是组合子** |
