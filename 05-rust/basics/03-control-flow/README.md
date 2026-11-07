# 03 · Rust · 控制流：分支、循环、迭代

## 这个例子讲什么

Rust 的 `if` 和 `match` 都是**表达式**（有返回值），循环有 `loop` / `while` / `for` 三种。重点体会：`match` 穷尽性由编译器强制、`break` 能带值、带标签的循环能跳出外层。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/03-control-flow
cargo run
```

## 关键点

- `if` 用作值时所有分支类型必须一致且必须覆盖（有 `else`）；`match` 少一个分支会直接编译失败，这就是「穷尽性检查」。
- `500..=599` 区间模式、`400 | 401 | 403` 或模式、`_` 通配，都是 match 的常规武器。
- `loop { ... break v; }` 中 `break` 的值成为 `loop` 表达式的结果，适合「重试直到成功」。
- `for x in 1..=10` 走 `Iterator`，比 C 风格下标循环更不易越界；配合 `filter` / `collect` 声明式地生成集合。
- `'outer:` 标签 + `break 'outer` 解决嵌套循环跳出外层的需求，Go 里对应带标签的 `break`，Java 也有标签中断。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 分支 | `if` / `switch` | `if` / `switch` | `if` / `match`(3.10) | `if` / `switch` 表达式 | `if` / `match` 表达式，穷尽性强制 |
| 三元 | 无 | `?:` | `a if c else b` | `?:` | `if ... else ...` 表达式 |
| 循环 | 只有 `for`（含 `range`） | `for`/`while`/`for..of` | `for`/`while` | `for`/`while`/`for-each` | `loop`/`while`/`for` 三种 |
| 迭代协议 | `range` | 迭代器 / 生成器 | `__iter__`/`__next__` | `Iterable`/`Iterator` | `Iterator` trait |
| 跳出外层 | 标签 `break` | 标签 `break` | 无（靠标志位） | 标签 `break` | 标签 `break 'label` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
