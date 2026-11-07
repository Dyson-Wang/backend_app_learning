# 15 · Python · 迭代器与生成器

## 这个例子讲什么

迭代协议 `__iter__`/`__next__`、生成器 `yield` 的惰性求值、`yield from`、`itertools` 的 `chain`/`takewhile`/`islice`、一次性消费陷阱、无限序列取前 N 项，以及 `iter(callable, sentinel)`。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/15-iterators-and-generators/main.py
```

## 关键点

- **迭代协议**：`iter(obj)` 拿迭代器、`next(it)` 推进，抛 `StopIteration` 结束；自定义类实现 `__iter__`（返回 self）+ `__next__` 即可被 `for` 消费。
- **惰性求值**：调用生成器函数只返回生成器对象，**函数体一行都不执行**；`yield` 之后才逐项产出，适合大文件 / 无限序列。
- **一次性消费**：`list(gen)` 第二次为空；要复用得 `list(...)` 缓存或重新生成。
- `yield from sub` 把子可迭代对象「摊平」进当前生成器；`itertools.count()` 是无限序列，必须配 `islice` 取前 N 项才不会爆内存。
- `iter(callable, sentinel)` 反复调用直到返回哨兵值停止，适合按块读文件。

## 和其他语言的差异

- **Go**：1.23 的 `iter.Seq` / range-over-func，或手写 `func(yield func(T) bool)` push 迭代器；标准库偏命令式。
- **Node.js/TS**：`function*` 生成器 + `Symbol.iterator` 迭代器协议，`for...of` 消费。
- **Java**：`Iterator` / `Iterable`，`Stream` 提供惰性 `map`/`filter`/`limit` 管道。
- **Rust**：`Iterator` trait + `map`/`filter`/`take`/`chain` 适配器 + `collect()`，零成本抽象。
- **Python**：`__iter__`/`__next__` 协议、`yield`/`yield from`、`itertools`、`iter(callable, sentinel)`。
