# C03 游标分页替代 OFFSET

## 问题复现

按热度倒序，每页 3 条。看完第一页后插入一条高分内容，翻第二页 `OFFSET 3` 时原本第 3 条重复出现。

## 错误写法

```sql
SELECT * FROM items ORDER BY score DESC LIMIT 3 OFFSET 3;
```

## 正确写法

```sql
SELECT * FROM items WHERE (score, id) < ($1, $2)
ORDER BY score DESC, id ASC LIMIT 3;
```

```python
if after is not None and not ((it.score, -it.id) > (after.score, -after.id)):
    continue
```

把「分数降序 + id 升序」映射成元组 `(score, -id)` 的升序比较，一行就能表达。

## 怎么验证

```bash
python3 main.py
```

预期：OFFSET 重复 `1` 条，游标重复 `0` 条。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 排序 | `sort.Slice` | `Array.sort` | `sorted(key=)` | `Comparator` | `sort_by` |
| 元组比较 | 手写 | 逻辑或 | 原生元组比较 | `Comparator` 链 | 手写 |
| 深分页 | 游标 | 游标 | 游标 | 游标 | 游标 |
