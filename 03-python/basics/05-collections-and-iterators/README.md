# 05 · Python · 集合与迭代器

## 这个例子讲什么

`list/tuple/set/dict` 四类内置集合、三种推导式、切片与 `sorted(key=)`、`collections` 的 `deque/Counter/defaultdict`，以及 `iter()` / `next()` 迭代协议与生成器的一次性消费。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/05-collections-and-iterators/main.py
```

## 关键点

- **`set` 无序去重、`tuple` 不可变可作 dict 键、`list` 可变**；`frozenset` 是不可变版 `set`。
- 推导式 `[..]` / `{k:v}` / `{..}` 分别产出 list / dict / set，比等价 `for` 循环更直观也更快。
- `defaultdict(list)` 免掉「键不存在先初始化」的样板；`Counter` 一行完成词频统计。
- **迭代协议**：`iter(obj)` 拿迭代器、`next(it)` 推进，耗尽抛 `StopIteration`，`for` 循环正是靠它收尾。
- 生成器/迭代器**一次性消费**：`list(gen)` 第二次就是空；要复用得转成 list 或重新生成（15 课展开）。

## 和其他语言的差异

- **Go**：`slice` / `map` / `array` + `range`；泛型容器需自己写或引第三方，标准库没有推导式。
- **Node.js/TS**：`Array` / `Map` / `Set` / `Object`、`for...of` 依赖 `Symbol.iterator` 迭代器协议。
- **Java**：`List` / `Map` / `Set` 集合框架、`Iterator`、`Stream` 惰性管道、`Collections` 工具类。
- **Rust**：`Vec` / `HashMap` / `HashSet` + `Iterator` trait + `collect()`，所有权决定迭代时 move 还是 borrow。
- **Python**：`list/dict/set/tuple/deque/Counter`、推导式、`iter()`/`next()` 协议、切片、`sorted(key=)`。
