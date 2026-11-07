# C03 游标分页替代 OFFSET

## 问题复现

按热度倒序每页 3 条，看完第一页后插入一条高分内容，翻第二页时原本第 3 条重复出现。

## 错误写法

```sql
SELECT * FROM items ORDER BY score DESC LIMIT 3 OFFSET 3;
```

## 正确写法

```sql
SELECT * FROM items WHERE (score, id) < ($1, $2)
ORDER BY score DESC, id ASC LIMIT 3;
```

```rust
let later = it.score < a.score || (it.score == a.score && it.id > a.id);
if !later { continue; }
```

## 怎么验证

```bash
cargo run --quiet
```

预期：OFFSET 重复 `1` 条，游标重复 `0` 条。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 排序 | `sort.Slice` | `Array.sort` | `sorted(key=)` | `Comparator` | `sort_by` + `then` |
| 元组比较 | 手写 | 逻辑或 | 原生元组 | `Comparator` 链 | 手写 |
| 深分页 | 游标 | 游标 | 游标 | 游标 | 游标 |
