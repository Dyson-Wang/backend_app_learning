# 15 · Node.js · 迭代器与生成器

## 这个例子讲什么

用 `function*`/`yield` 写自定义生成器，实证迭代器协议（生成器既是迭代器也是可迭代对象）；再用手写的 `map`/`filter`/`take`/`chain` 搭惰性管道，演示惰性求值、`yield*` 拼接、无限序列取前 N 项，以及「生成器只能消费一次」的陷阱。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/15-iterators-and-generators/index.mjs
```

无端口、无定时器；`naturals()` 是惰性无限生成器，被 `take(N)` 限定后正常结束。

## 关键点

- 生成器函数 `function*` 返回一个对象，既有 `next()` 又有 `[Symbol.iterator]()`（返回自身），因此可直接 `for...of` / 展开。
- `yield` 暂停并交出值；`return` 值只在 `{done:true}` 那次出现，`for...of` 与展开都拿不到。
- 惰性：构造 `map/filter/take` 管道时一次都不计算，只有消费才拉动；示例里 `take(2)` 让上游只计算到够用为止。
- `take` 用 `return` 提前结束，会触发 `IteratorClose` 调用上游生成器的 `return()`，从而正确关闭无限生成器（`finally` 会被执行）——这是「无限序列不会变死循环」的关键。
- `yield* it` 把另一个可迭代对象的所有值委托产出，用来实现 `chain`。
- `while (true) yield n++` 定义的是**惰性**无限序列，不是运行期死循环；必须配 `take`/`break` 限定消费量。
- 一次性陷阱：生成器消费完就空了，第二次 `[...gen]` 得到 `[]`；数组/`Set`/`Map` 每次 `for...of` 都新建迭代器，可反复遍历。要复用生成器结果就先摊成数组。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 生成器 | 1.21+ `iter.Seq` / push 函数 | `yield` / `yield from` | Stream / `Iterator` | `Iterator` trait | `function*`/`yield` |
| 迭代器协议 | `range` / `func(yield)` | `__iter__`/`__next__` | `hasNext`/`next` | `Iterator::next` | `[Symbol.iterator]` + `next()` |
| 惰性算子 | 手写 / `iter` 包 | `itertools` | Stream API | `map/filter/take/chain` | 手写生成器 / 数组方法 |
| 无限序列 | 手写 | `itertools.count` | `Stream.iterate` | `0..` / `iterate` | `while(true) yield` |
| 一次性 | 是（Seq） | 是（生成器） | 是（Iterator） | 是（Iterator） | 是（生成器） |

Python 的生成器/`itertools` 与 JS 的 `function*`/惰性管道几乎一一对应（`yield from` ≈ `yield*`）；Rust 的 `Iterator` 是零成本抽象、默认惰性且编译期单态化；Go 直到 1.21 才有 `iter.Seq` 这类标准迭代器；Java 走 Stream API 且不可复用。共同点是「迭代器一次性消费」，数组类容器则可反复遍历。
