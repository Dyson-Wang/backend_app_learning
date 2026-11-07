# C02 N+1 查询与批量写

## 问题复现

列表页展示 10 条订单和下单人姓名：先查订单（1 次），再循环查每个人（N 次），共 11 次往返。

## 错误写法

```java
for (Order o : db.listOrders()) {
    db.getAccount(o.accountId);      // 每条一次
}
```

JPA 里更隐蔽：`@ManyToOne(fetch = LAZY)` 在循环里触发代理加载。

## 正确写法

```java
List<Order> orders = db.listOrders();
Map<Integer, Account> m = db.getAccounts(ids);   // WHERE id IN (...)
```

| 语言 | 做法 |
|---|---|
| Go | `WHERE id = ANY($1)` + map 拼装 |
| Node.js | Prisma `include` / Drizzle `with` |
| Python | SQLAlchemy `selectinload` |
| Java | `@EntityGraph` / `JOIN FETCH` |
| Rust | SeaORM `find_with_related` |

## 怎么验证

```bash
javac Main.java && java Main
```

预期 `查询 11 次` vs `查询 2 次`，`结果一致：true`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 惯用解法 | 显式批量 | ORM 预加载 | `selectinload` | `JOIN FETCH` | join |
| 隐形陷阱 | 循环查 DB | 循环 await | 循环查 ORM | 懒加载代理 | 循环 await |
| 检测 | 慢日志 | 查询计数 | `echo=True` | `show_sql` | query 日志 |
