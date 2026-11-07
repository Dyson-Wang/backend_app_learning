# 15 · Rust · 迭代器与生成器

## 这个例子讲什么

Rust 的迭代器协议：实现 `Iterator` 只需写一个 `next()`，就自动获得 `map` / `filter` / `take` / `chain` / `sum` / `collect` 等几十个适配器。重点是**惰性求值**、自定义迭代器、无限序列 + `take`、以及迭代器「一次性消费」的语义。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/15-iterators-and-generators
cargo run --quiet
```

## 关键点

- `impl Iterator for Counter` 只要求 `type Item` 和 `fn next`，其余方法全是默认实现。
- 适配器是惰性的：`map` 只是包一层，不调用 `next` 就不会执行闭包；`collect` / `sum` / `for` 才是消费者。
- `successors` 可以造无限序列，必须配 `take` 截断，否则 `collect` 会 OOM。
- 迭代器被 `collect` 或 `sum` 消费后所有权转移；`next()` 是 `&mut self`，会推进内部状态，不能重复从头读。
- Rust 没有 `yield` 生成器语法（`gen` 块仍在 nightly），「生成器」用「实现 Iterator 的结构体」或 `successors` 表达。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 迭代协议 | `range` / channel | `Symbol.iterator` | `__iter__` / `__next__` | `Iterator<T>` | `Iterator` trait |
| 生成器语法 | 无（用 channel 模拟） | `function*` + `yield` | `yield` / `yield from` | 无（用 `Stream`） | 无稳定语法，用结构体实现 |
| 惰性管道 | 手写循环 | 生成器链 | `itertools` / 生成器 | `Stream` | 适配器链 |
| 一次性消费 | channel 关闭 | 生成器耗尽 | 生成器耗尽 | `Iterator` 耗尽 | 迭代器 move |

Python 和 JS 有真正的 `yield` 生成器，写起来最省事；Go 靠 channel + goroutine 模拟；Java 的 `Stream` 是惰性管道但不是「生成器」；Rust 用 trait 实现，性能最好但没有语法糖。
