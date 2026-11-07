# C03 游标分页替代 OFFSET

## 问题复现

列表按热度倒序，每页 3 条。用户看完第一页（1/2/3），此时有人发布了一条高分内容，插到了第 2 位。用户翻第二页时 `OFFSET 3`，原本的第 3 条又出现了一次——**重复**；反过来删数据就会**漏读**。

## 错误写法

```sql
SELECT * FROM items ORDER BY score DESC LIMIT 3 OFFSET 3;
```

`OFFSET` 数的是「行位置」，前面的增删会让位置整体偏移。

## 正确写法

用上一页最后一条的排序值做锚点：

```sql
SELECT * FROM items
WHERE (score, id) < ($1, $2)
ORDER BY score DESC, id ASC
LIMIT 3;
```

```go
if after != nil && !(it.Score < after.Score || (it.Score == after.Score && it.ID > after.ID)) {
    continue
}
```

锚点必须带上唯一列（`id`）做 tiebreak，否则同分数据会互相跳过。

## 怎么验证

```bash
go run .
```

预期：`OFFSET` 第二页与第一页重复 `1` 条；游标分页重复 `0` 条。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 排序 | `sort.Slice` | `Array.sort` | `sorted(key=)` | `Comparator` | `sort_by` |
| 元组比较 | 手写 || | 元组比较 | `Comparator` 链 | 手写 || |
| 深分页 | 游标 / 延迟关联 | 同 | 同 | 同 | 同 |
