# B01 · Go · 事务与隔离级别、乐观锁

## 问题复现

两个事务都读到 `balance=100`，各自扣 10 后写回 90。数据库里应该是 80，实际是 90——第二次写入覆盖了第一次，这就是**丢失更新**。读已提交隔离级别下，这种交错完全合法。

## 错误写法

```go
balance, _ := read()          // 两个事务都读到 100
write(balance - 10)           // 都写 90
```

读和写之间没有任何版本校验。

## 正确写法

- **悲观锁**：`SELECT ... FOR UPDATE`，读的时候就锁住行，简单但吞吐低、易死锁；
- **乐观锁**：加 `version` 列，`UPDATE ... SET balance=?, version=version+1 WHERE id=? AND version=?`，影响行数为 0 就重试；
- **可串行化**：数据库替你检测冲突，冲突时回滚，代价是并发度最低。

本例用内存 `cas(expectedVersion, newBalance)` 模拟乐观锁，冲突时重试，最终 80。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/B01-transactions-and-isolation
go run .
```

预期：无锁 `实际 90`；乐观锁 `T2 用旧 version 写入失败`，重试后 `balance=80`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 事务 API | `database/sql` `Tx` | `pg` / `mysql2` | `psycopg` / SQLAlchemy | JDBC / Spring `@Transactional` | `sqlx` / `tokio-postgres` |
| 乐观锁 | 手写 version 条件 | 手写 | 手写 / SQLAlchemy `version_id_col` | JPA `@Version` | 手写 |
| 隔离级别 | `sql.TxOptions` | 连接参数 | 连接参数 | `@Transactional(isolation=...)` | `sqlx` 事务配置 |

Java 的 JPA `@Version` 是开箱即用的乐观锁；Go/Rust 需要自己写 `WHERE version=?` 并检查 `RowsAffected`；Python 的 SQLAlchemy 有 `version_id_col`。注意「读已提交」下乐观锁仍可能遇到幻读，范围查询要配合间隙锁或可串行化。
