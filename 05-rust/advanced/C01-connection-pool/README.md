# C01 连接池与资源耗尽

## 问题复现

池只有 3 个连接，10 个并发请求各占 120ms。不配等待超时，7 个请求无限排队；下游一慢，队列爆炸。

## 错误写法

```rust
let conn = pool.acquire(None);      // 永久阻塞
```

## 正确写法

```rust
match pool.acquire(Some(Duration::from_millis(40))) {
    true => run_query(conn),
    false => return Err(ServiceUnavailable),
}
```

用 `Condvar::wait_timeout` 实现，不要忙等。

## 怎么验证

```bash
cargo run --quiet
```

预期：第一种 `ok=10 timeout=0 耗时=400ms 左右`；第二种 `ok=3 timeout=7 耗时=120ms 左右`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 池实现 | `database/sql` | `pg.Pool` | `asyncpg` | HikariCP | `sqlx::Pool` / `bb8` |
| 等待超时 | `context` | `connectionTimeoutMillis` | `timeout` | `connectionTimeout` | `acquire_timeout` |
| 实现手段 | channel / 信号量 | 队列 | `Semaphore` | `Semaphore` | `Condvar::wait_timeout` |
