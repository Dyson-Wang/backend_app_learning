# SQL 与 PostgreSQL

严格说 SQL 不是「后端语言」，但它是所有后端语言共用的底座，优先级高于第 6 门语言。

## 必须手写过的四件事

1. **事务与隔离级别**：`READ COMMITTED` / `REPEATABLE READ` / `SERIALIZABLE` 下分别跑一遍并发扣库存，亲眼看到超卖和幻读。
2. **索引与执行计划**：`EXPLAIN ANALYZE` 看懂 `Seq Scan` 与 `Index Scan`，知道哪些写法会让索引失效（函数包裹列、隐式类型转换、前导模糊匹配）。
3. **N+1 查询**：先用循环查 100 次，再用 `JOIN` / `IN` 批量查一次，对比耗时。
4. **连接池耗尽**：把池大小设为 2，开 20 个并发慢查询，观察排队与超时。

`schema.sql` 给了一个可直接执行的最小库表，5 门语言的 `advanced/B01` 都会用它。

## 隔离级别速查

| 级别 | 脏读 | 不可重复读 | 幻读 | PostgreSQL 默认 |
|---|---|---|---|---|
| READ UNCOMMITTED | 可能 | 可能 | 可能 | PG 实际等同 READ COMMITTED |
| READ COMMITTED | 不会 | 可能 | 可能 | 是 |
| REPEATABLE READ | 不会 | 不会 | 不会（PG 用 MVCC 快照） | 否 |
| SERIALIZABLE | 不会 | 不会 | 不会 | 否 |

## 常见坑

- **长事务**：事务里做 HTTP 调用 / 发消息，会长时间占着连接和锁。
- **先删后插**代替 upsert：并发下产生唯一键冲突，应使用 `INSERT ... ON CONFLICT DO UPDATE`。
- **分页用 `OFFSET`**：深分页越来越慢，改用游标分页（见 `advanced/C03`）。
- **隐式类型转换**：字符串列传数字导致索引失效。
- **没给查询设 statement_timeout**：一条慢 SQL 就能把连接池堵死（见 `advanced/F01`）。

## 接 homelab

`docker-services` 已提供 PostgreSQL。用环境变量传入连接串，不要把主机、端口、口令写进示例代码：

```
export BRM_PG_DSN='postgres://user:password@host:5432/brm?sslmode=disable'
```

`schema.sql` 假设库名 `brm`，可按实际环境调整。
