# D01 至少一次投递与重复消费

## 问题复现

消费者处理完消息、还没 ACK 就重启了，Broker 认为没送达，重投一次。同一条消息被处理两次，余额多加了一次。

## 错误写法

```go
balance += e.Amount      // 来几次加几次
```

## 正确写法

消费端维护去重表，**和业务写在同一事务**：

```go
if seen[e.ID] { return }   // 已处理，直接 ACK
seen[e.ID] = true
balance += e.Amount
```

```sql
BEGIN;
INSERT INTO consumed (msg_id) VALUES ($1) ON CONFLICT DO NOTHING;
-- 影响 0 行 = 重复，ROLLBACK
UPDATE accounts SET balance = balance + $2 WHERE id = $3;
COMMIT;
```

## 怎么验证

```bash
go run .
```

预期：无幂等 `100 -> 140`；有幂等 `实际处理 3 次，100 -> 130`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 去重载体 | `map` / Redis SETNX | `Map` / Redis | `set` / Redis | `Set` / Redis | `HashSet` |
| 事务边界 | `sql.Tx` | 驱动事务 | `with conn.transaction()` | `@Transactional` | `sqlx::Transaction` |
| 顺序保证 | 单分区单消费者 | 同 | 同 | 同 | 同 |
