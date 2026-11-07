# C02 N+1 查询与批量写

## 问题复现

列表页要展示 10 条订单和下单人名字。代码先查订单列表（1 次），再循环查每个下单人（N 次），总共 11 次往返。列表一长就是灾难。

## 错误写法

```python
for o in db.list_orders():
    acc = db.get_account(o.account_id)   # 每条一次
```

## 正确写法

一次把需要的 id 全查回来，在内存里拼：

```python
orders = db.list_orders()
ids = {o.account_id for o in orders}
mapping = db.get_accounts(ids)           # WHERE id = ANY($1)
```

ORM 里的对应做法：

| 语言 | 做法 |
|---|---|
| Go | 手写 `WHERE id = ANY($1)` + map 拼装 |
| Node.js | Prisma `include` / Drizzle `with` |
| Python | SQLAlchemy `selectinload` |
| Java | JPA `@EntityGraph` / `JOIN FETCH` |
| Rust | SeaORM `find_with_related` / 手写 join |

## 怎么验证

```bash
python3 main.py
```

预期 `查询 11 次` vs `查询 2 次`，`结果一致：True`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 惯用解法 | 显式批量查询 | ORM 预加载 | `selectinload` | `JOIN FETCH` | join / 预加载 |
| 隐形陷阱 | 循环里查 DB | 循环 `await` | 循环查 ORM | 懒加载代理 | 循环 `await` |
| 检测手段 | 慢日志 + 查询计数 | 查询计数 | `echo=True` | `hibernate.show_sql` | `sqlx::query` 日志 |
