# 05 · Node.js · 集合与迭代器

## 这个例子讲什么

覆盖 Array（`map`/`filter`/`reduce`/`sort`/`slice`/`splice`）、`Set`、`Map`、`WeakMap`，并手写一个实现 `[Symbol.iterator]` 的可迭代对象，说明 `for...of` / 展开运算符 `...` / `Array.from` 背后统一的迭代器协议。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/05-collections-and-iterators/index.mjs
```

## 关键点

- `map`/`filter`/`slice` 返回新数组；`push`/`splice`/`sort` 原地修改——`sort` 默认字典序，数字排序要 `(a,b)=>a-b`。
- `Set` 去重且保持插入顺序；`Map` 允许任意类型作键并保持插入顺序，比普通对象更适合当字典。
- `WeakMap` 键是弱引用、不阻止 GC、不可遍历，适合给对象挂私有数据（21 课讲内存）。
- 迭代器协议：对象实现 `[Symbol.iterator]()` 返回带 `next()` 的迭代器，`next()` 返回 `{value, done}`；`for...of`、`...`、`Array.from` 都靠它。
- 迭代器让自身也实现 `[Symbol.iterator]() { return this; }`，就能被反复 `for...of` 消费。
- 普通对象默认不可迭代，遍历要用 `Object.keys/values/entries`；生成器（`function*`）见 15 课。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 动态数组 | `slice` | `list` | `ArrayList` | `Vec` | `Array` |
| 哈希集合 | 用 `map[T]bool` 模拟 | `set` | `HashSet` | `HashSet` | `Set` |
| 有序字典 | 无（map 无序） | `dict`（3.7+ 有序） | `LinkedHashMap` | 需 `IndexMap` 库 | `Map`（插入有序） |
| 迭代协议 | `range` / 无统一协议 | `__iter__`/`__next__` | `Iterable`/`Iterator` | `Iterator` trait | `[Symbol.iterator]` + `next()` |
| 函数式算子 | 无，手写循环 | `map`/`filter`/itertools | Stream API | `Iterator::map/filter` | 数组自带 `map/filter/reduce` |

Go 的 `map` 遍历顺序是刻意随机化的，Rust 标准库也没有有序 Map；JS 的 `Map` 和 Python 3.7+ 的 `dict` 都保证插入顺序，横向对照时这一点最容易被忽略。
