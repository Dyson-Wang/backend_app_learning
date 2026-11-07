# C02 · Go · N+1 查询与批量写

## 问题复现

查 10 个订单，再为每个订单查一次用户：1 + 10 = 11 次数据库往返。每个往返 1ms，接口就多 10ms；如果是跨机房 20ms，就是 200ms。数据量涨到 100 条，查询次数线性涨到 101 次。

## 错误写法

```go
for _, o := range orders {
    u := db.UserByID(o.UserID) // 每条一次查询
    out = append(out, ...)
}
```

## 正确写法

- 收集所有 `UserID`，一次 `WHERE id IN (...)` 查回来，内存里做 map 关联；
- 或者用 JOIN，让数据库一次返回；
- 批量写同理：`INSERT ... VALUES (...), (...)` 或 `COPY`，不要循环单条插入。

本例用计数器直接对比：11 次 vs 2 次。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/C02-n-plus-one
go run .
```

预期：`N+1 ... 查询 11 次`，`批量 ... 查询 2 次`，`第一条结果一致：true`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| ORM | GORM / sqlc | Prisma / TypeORM | SQLAlchemy / Django ORM | JPA / MyBatis | Diesel / SeaORM |
| 预加载 | `Preload` | `include` | `joinedload` | `@EntityGraph` | `preload` |
| 批量插入 | 手写多值 INSERT | `createMany` | `bulk_insert_mappings` | `saveAll` / JDBC batch | 手写 / `QueryBuilder` |

ORM 的「懒加载」是 N+1 的主要来源：Django 的 `select_related` / `prefetch_related`、JPA 的 `@EntityGraph`、GORM 的 `Preload` 都是用来关掉懒加载的。开 `log` 看 SQL 条数是最快的排查方式。
