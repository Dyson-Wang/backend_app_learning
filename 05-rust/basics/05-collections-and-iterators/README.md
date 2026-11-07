# 05 · Rust · 集合与迭代器

## 这个例子讲什么

标准库四类常用集合：`Vec` 动态数组、`HashMap` 无序映射、`BTreeMap` 有序映射、`HashSet` 去重集合，以及贯穿所有集合的 `Iterator` 惰性适配器（`map` / `filter` / `collect` / `enumerate` / `sum`）。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/05-collections-and-iterators
cargo run
```

## 关键点

- `Vec<T>` 堆上连续存储；`sort` 需要 `&mut self`，索引 `nums[0]` 越界会 panic。
- `HashMap` 顺序**不确定**，`BTreeMap` 按键有序；本例把 `HashSet` 转成有序 `Vec` 再打印，正是为了让输出稳定可复现。
- `entry(key).or_insert(default)` 返回 `&mut V`，可原地更新，避免「先查再插」两次哈希。
- 迭代器**惰性**：`map` / `filter` 只是记录步骤，直到 `collect` / `sum` / `for` 才真正执行；且消费后不能复用（一次性）。
- `iter()` 借用、`into_iter()` 拿走所有权、`iter_mut()` 可变借用——三者由类型系统区分。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 动态数组 | slice | `Array` | `list` | `ArrayList` | `Vec<T>` |
| 映射 | `map` | `Map` / 对象 | `dict` | `HashMap` / `TreeMap` | `HashMap` / `BTreeMap` |
| 集合 | 无内置 | `Set` | `set` | `HashSet` | `HashSet` / `BTreeSet` |
| 迭代 | `range` | `for..of` / 迭代器 | `iter` / 生成器 | `Iterator` / `Stream` | `Iterator` trait + 适配器 |
| 惰性管道 | 手写 | 数组方法链（急切） | 生成器 / `itertools` | `Stream`（惰性） | 适配器默认惰性 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
