# 03 · Python · 控制流：分支、循环、迭代

## 这个例子讲什么

`if/elif/else`、条件表达式、`for...else` / `while...else`、`break/continue`、海象运算符 `:=`，以及 `range` / `enumerate` / `zip` / `any` / `all` 这些迭代惯用法。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/03-control-flow/main.py
```

## 关键点

- **只有 `for` 和 `while`**：没有 `do...while`，没有 C 风格三段 `for`；`range()` 生成惰性整数序列。
- **`for...else` 是 Python 独有**：循环**没被 `break`** 才执行 `else`，常用来写「查找失败」分支，最容易被误读成 `if...else`。
- **链式比较**：`0 < n < 10` 等价于 `0 < n and n < 10`，且中间值只求值一次。
- `:=` 海象运算符在表达式里顺带赋值，适合 `while (line := f.readline()):` 这类写法。
- 条件表达式 `A if cond else B` 是表达式而非语句；`and/or` 短路且返回操作数本身。

## 和其他语言的差异

- **Go**：只有 `for` 一种循环（兼顾 while 与三段式）、`switch` 自动 `break`、`if` 可带初始化语句；没有 `for...else`、没有三元。
- **Node.js/TS**：`for...of` / `for...in`、`switch`、truthy/falsy 判断、有 `?:` 三元与 `?.` 可选链。
- **Java**：增强 for、`switch` 表达式（14+）、带标签的 `break`、无 `for...else`。
- **Rust**：`loop` / `while` / `for` 三种循环、`match`、循环标签 `'label: loop`、用 `if` 表达式代替三元。
- **Python**：`for...else`、`while...else`、链式比较、海象 `:=`、内置 `enumerate`/`zip`/`any`/`all`。
