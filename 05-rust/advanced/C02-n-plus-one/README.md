# C02 N+1 查询与批量写

## 问题复现

列表页 10 条订单 + 下单人姓名：先查订单（1 次），再循环查人（N 次），共 11 次往返。

## 错误写法

```rust
for o in db.list_orders() {
    let name = db.get_account(o.account_id);   // 每条一次
}
```

## 正确写法

```rust
let ids: Vec<usize> = orders.iter().map(|o| o.account_id)
    .collect::<BTreeSet<_>>().into_iter().collect();
let mapping = db.get_accounts(&ids);           // WHERE id = ANY($1)
```

## 怎么验证

```bash
cargo run --quiet
```

预期 `查询 11 次` vs `查询 2 次`，`结果一致：true`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 惯用解法 | 显式批量 | ORM 预加载 | `selectinload` | `JOIN FETCH` | join / `find_with_related` |
| 隐形陷阱 | 循环查 DB | 循环 await | 循环查 ORM | 懒加载代理 | 循环 await |
| 检测 | 慢日志 | 查询计数 | `echo=True` | `show_sql` | query 日志 |
