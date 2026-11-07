# 15 · 迭代器与生成器

## 这个例子讲什么

Go 没有 `yield` 关键字，也没有内建迭代器协议。本例用「返回 `func(yield func(T) bool)` 的 push 迭代器」手写一套惰性管道：`Map` / `Filter` / `Take` / `Chain` / `Collect`，演示惰性求值、无限序列取前 N 项、自定义生成器，以及 channel 支撑的迭代器「只能消费一次」的陷阱。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/15-iterators-and-generators
```

## 关键点

- **push 迭代器协议**：`Seq[T] = func(yield func(T) bool)`。迭代器主动把值推给 `yield`，`yield` 返回 `false` 即提前停止。这和 Go 1.23 的 `iter.Seq` / `range-over-func` 是同一套协议。
- **版本边界**：本仓库 `go.mod` 是 **go 1.21**，而 `iter.Seq` 与 `for x := range seq`（range over func）是 **Go 1.23** 才有的，1.21 里 `for range` 一个函数会编译失败。所以这里全程手动调用 Seq，不写 range-over-func。
- **惰性**：`Map`/`Filter`/`Take` 只返回闭包、不碰数据；直到 `Collect` 真正带 `yield` 调用时才逐项计算。无限 `Nats()` 配合 `Take(n)` 才能安全取前 N 项。
- **组合子实现要点**：`Filter` 跳过时 `yield` 不调用但仍返回 `true` 继续；`Take` 用计数达到 n 后返回 `false` 掐断上游；`Chain` 用 `more` 标志决定前半段是否提前中止，从而避免继续读后半段。
- **一次性消费陷阱**：`ChanSeq` 背后是 channel，`Collect` 一次就抽干，第二次得到空切片；而 `From`/`Nats` 这类纯函数式 Seq 可以反复消费。channel、文件、网络流都是一次性资源。
- 没有 `for...of` / `Symbol.iterator` 之类语法糖，消费端必须自己写 `yield` 回调或 `Collect`。

## 和其他语言的差异

| 语言 | 迭代协议 | 惰性 / 无限序列 | 生成器 |
|---|---|---|---|
| Go | 1.21 手写 push `func(yield func(T) bool)`；1.23 才有 `iter.Seq` + range-over-func | 手写 `Take` 掐断 | 无 `yield` 关键字，用闭包模拟 |
| Node.js | `Symbol.iterator` / `next()`，`for...of` 语法糖 | Generator 天然惰性 | `function*` + `yield`、`yield*` |
| Python | `__iter__` / `__next__`、`StopIteration` | 生成器惰性，`itertools.islice` 取前 N | `yield` / `yield from`、`itertools` |
| Java | `Iterable` / `Iterator`；Stream API | `Stream` 惰性，`limit(n)` | 无原生生成器，借 `Stream` / `Spliterator` |
| Rust | `Iterator` trait（`next`），`for` 语法糖 | 适配器默认惰性，`take(n)` | 无 `yield`（稳定版），自定义 `impl Iterator` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
