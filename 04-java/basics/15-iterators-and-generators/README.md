# 15 · Java · 迭代器与生成器

## 这个例子讲什么

`Iterator` 迭代协议、自定义惰性迭代器、一次性消费陷阱、`Stream` 惰性管道 `map`/`filter`/`take(limit)`/chain、无限序列取前 N 项、流只能消费一次。

## 怎么跑

```
cd backend-roadmap/04-java/basics/15-iterators-and-generators
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **Java 没有 `yield` 生成器**：惰性序列用 `Iterator` 手写，或用 `Stream` / `Spliterator`。`next()` 被调一次才算产出一个值。
- **一次性消费**：`Iterator` 与 `Stream` 都不能复用；示例里第二次 `drain(counter)` 得到 0、第二次 `forEach` 抛 `IllegalStateException`。想重复遍历就重新构造。
- `Stream.iterate(seed, f)` 造**无限序列**，必须靠 `limit(n)`（take）截断，否则终止操作永不返回——这正是「禁止无限循环」要注意的。
- `Stream` 的中间操作（`map`/`filter`）是**惰性**的，只有终止操作（`sum`/`toList`/`forEach`）才真正触发遍历。
- `Stream.concat` 对应 `chain`；`String.lines()`、`Stream.toList()`（Java 16+）都是现成工厂。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 生成器 | 无 `yield`，用 `Iterator`/`Stream` | `iter.Seq`/push 迭代器（1.23） | `function*` / `yield` | `yield` / `itertools` | 无 `yield`，`Iterator` trait |
| 惰性管道 | `Stream` | 手写 chain/take | 生成器 + 辅助函数 | 生成器表达式 | `Iterator` 适配器 |
| take / chain | `limit` / `Stream.concat` | 手写 | 手写 / 库 | `islice` / `chain` | `.take(n)` / `.chain()` |
| 一次性消费 | 是（Iterator/Stream） | 是（迭代器） | 是（生成器） | 是（生成器） | 是（`Iterator` 消耗） |
| 迭代协议 | `hasNext`/`next` | `range` / `Seq` | `Symbol.iterator` | `__iter__`/`__next__` | `Iterator::next` |
