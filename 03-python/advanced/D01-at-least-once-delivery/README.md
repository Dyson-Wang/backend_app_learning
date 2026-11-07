# D01 至少一次投递与重复消费

## 问题复现

消费者处理完、还没 ACK 就重启，Broker 重投，同一条消息处理两次。

## 错误写法

```python
balance += amount
```

## 正确写法

```python
if msg_id in seen:
    continue
seen.add(msg_id)
balance += amount
```

```sql
BEGIN;
INSERT INTO consumed (msg_id) VALUES (%s) ON CONFLICT DO NOTHING;
UPDATE accounts SET balance = balance + %s WHERE id = %s;
COMMIT;
```

## 怎么验证

```bash
python3 main.py
```

预期：无幂等 `100 -> 140`；有幂等 `实际处理 3 次，100 -> 130`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 去重载体 | `map` | `Set` | `set` | `Set` | `HashSet` |
| 事务边界 | `sql.Tx` | 驱动事务 | `conn.transaction()` | `@Transactional` | `Transaction` |
| 顺序保证 | 单分区单消费者 | 同 | 同 | 同 | 同 |
