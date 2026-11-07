# D01 至少一次投递与重复消费

## 问题复现

消费者处理完、还没 ACK 就重启，Broker 重投，同一条消息处理两次。

## 错误写法

```java
balance += e.amount;
```

## 正确写法

```java
if (seen.contains(e.id)) return;
seen.add(e.id);
balance += e.amount;
```

```sql
BEGIN;
INSERT INTO consumed (msg_id) VALUES (?) ON CONFLICT DO NOTHING;
UPDATE accounts SET balance = balance + ? WHERE id = ?;
COMMIT;
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：无幂等 `100 -> 140`；有幂等 `实际处理 3 次，100 -> 130`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 去重载体 | `map` | `Set` | `set` | `Set` | `HashSet` |
| 事务边界 | `sql.Tx` | 驱动事务 | `conn.transaction()` | `@Transactional` | `Transaction` |
| 顺序保证 | 单分区单消费者 | 同 | 同 | 同 | 同 |
