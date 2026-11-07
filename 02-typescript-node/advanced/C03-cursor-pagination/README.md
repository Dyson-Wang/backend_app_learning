# C03 游标分页替代 OFFSET

## 问题复现

按热度倒序，每页 3 条。看完第一页后有人发布高分内容插到第 2 位，翻第二页 `OFFSET 3` 时，原本第 3 条又出现一次。

## 错误写法

```sql
SELECT * FROM items ORDER BY score DESC LIMIT 3 OFFSET 3;
```

## 正确写法

```sql
SELECT * FROM items WHERE (score, id) < ($1, $2)
ORDER BY score DESC, id ASC LIMIT 3;
```

```js
const later = it.score < after.score || (it.score === after.score && it.id > after.id);
if (!later) continue;
```

## 怎么验证

```bash
node index.mjs
```

预期：OFFSET 第二页重复 `1` 条，游标重复 `0` 条。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 排序 | `sort.Slice` | `Array.sort` | `sorted(key=)` | `Comparator` | `sort_by` |
| 元组比较 | 手写 | 逻辑或 | 元组比较 | `Comparator` 链 | 手写 |
| 深分页 | 游标 | 游标 | 游标 | 游标 | 游标 |
