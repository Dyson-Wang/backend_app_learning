# C02 · Node.js · N+1 查询与批量写

## 问题复现

查 10 个订单，再为每个订单查一次用户：1 + 10 = 11 次往返。每个往返 1ms 就多 10ms；跨机房 20ms 就是 200ms。数据涨到 100 条，查询变 101 次。

## 错误写法

```js
for (const o of orders) {
  const u = await db.userById(o.userId); // 每条一次
}
```

## 正确写法

- 收集去重后的 `userId`，一次 `WHERE id IN (...)` 查回来，用 `Map` 关联；
- 或直接 JOIN；
- 批量写用 `INSERT ... VALUES (...), (...)` / `createMany`，不要循环单条插。

本例用计数器对比：11 次 vs 2 次。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/C02-n-plus-one
node index.mjs
```

预期：`N+1 ... 查询 11 次`，`批量 ... 查询 2 次`，`第一条结果一致：true`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| ORM | GORM / Prisma | Prisma / TypeORM | SQLAlchemy / Django | JPA / MyBatis | Diesel / SeaORM |
| 预加载 | `Preload` | `include` | `joinedload` | `@EntityGraph` | `preload` |
| 批量插入 | 多值 INSERT | `createMany` | `bulk_insert_mappings` | JDBC batch | 手写 |

Prisma 的 `include` 只做 JOIN 或额外一次 IN 查询，不会 N+1；TypeORM 的 `relations` 默认懒加载，很容易踩坑。排查手段：打开 SQL 日志数条数。
