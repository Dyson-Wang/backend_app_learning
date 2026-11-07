# 05 · Java · 集合与迭代器

## 这个例子讲什么

`List/Set/Map/Queue` 四大接口与常用实现、`Iterator.remove`、`for-each` 快速失败、`Comparator` 排序、`List.of` 不可变集合。

## 怎么跑

```
cd backend-roadmap/04-java/basics/05-collections-and-iterators
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- 集合框架分**接口**（`List/Set/Map/Queue`）与**实现**（`ArrayList/HashMap/HashSet/ArrayDeque`），面向接口编程。
- `Iterator` 遍历时想删元素必须用 `it.remove()`，且要先 `next()`。
- `for-each` 期间直接增删集合会触发 `ConcurrentModificationException`（**快速失败**是 best-effort，不保证每种情况都抛）。
- `List.of` / `Set.of` / `Map.of` 返回**不可变**集合，修改抛 `UnsupportedOperationException`。
- `HashMap` 不保证顺序；要稳定顺序用 `LinkedHashMap`（插入序）或 `TreeMap`（键序）。示例按固定键序打印，保证输出可复现。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 动态数组 | `ArrayList` | slice | `Array` | `list` | `Vec<T>` |
| 哈希映射 | `HashMap` | `map[K]V` | `Map` / 对象 | `dict` | `HashMap` |
| 去重集合 | `HashSet` | `map[T]struct{}` | `Set` | `set` | `HashSet` |
| 迭代 | `Iterator` / for-each | `range` | `for-of` / `Symbol.iterator` | `__iter__`/`__next__` | `Iterator` trait |
| 不可变 | `List.of` | 无（约定） | `Object.freeze` | `tuple`/`frozenset` | 默认不可变需 `mut` |
