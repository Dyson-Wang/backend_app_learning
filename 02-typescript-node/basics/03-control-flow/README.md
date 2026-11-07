# 03 · Node.js · 控制流：分支、循环、迭代

## 这个例子讲什么

覆盖 `if/else`、三元、`switch`、`for`/`while`/`do...while`、`break`/`continue`/标签 break，以及 `for...of` 与 `for...in` 的本质区别。重点看迭代靠的是「迭代器协议」而不是下标，所以 Map / Set / 字符串都能直接遍历。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/03-control-flow/index.mjs
```

## 关键点

- `switch` 用 `===` 严格匹配，漏写 `break` 会 fall-through 贯穿到下一个 case。
- 三元表达式适合简单二选一；嵌套三元可读性差，改回 `if`。
- `continue` 跳过本轮、`break` 结束本层；`标签: for` 配 `break 标签` 可一次跳出多层循环（Go 里对应 `break LABEL`）。
- `for...in` 遍历的是可枚举键（数组就是下标字符串），`for...of` 遍历的是值，依赖 `[Symbol.iterator]`，是 15 课生成器的基础。
- `for...of` 配合解构 `for (const [k, v] of map)` 是遍历 Map 的惯用写法。
- 普通对象默认不可迭代，要 `Object.keys/values/entries`（见 05）。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 分支 | `if/else`、`switch` | `if/elif/else` | `if/else`、`switch` | `if/else`、`match` | `if/else`、`switch`、三元 `?:` |
| 循环 | 只有 `for`（含 `range`） | `for ... in`、`while` | `for`、`while`、增强 for | `for`、`while`、`loop` | `for`/`for...of`/`for...in`/`while`/`do` |
| 迭代协议 | `range`、无统一协议 | `__iter__`/`__next__` | `Iterable`/`Iterator` | `Iterator` trait | `[Symbol.iterator]` + `for...of` |
| 跳出多层 | `break LABEL` | 无，靠标志位 | `break label` | `'label: loop` | `outer: for` + `break outer` |
| 模式分支 | `switch` + `type switch` | 3.10 `match case` | 17+ switch 表达式 | 强大 `match` | 无模式匹配，仅值/范围判断 |

Go 只有一种 `for` 关键字却覆盖全部循环场景；JS 循环种类最多，且「值迭代」与「键迭代」分成两套语法，这是最容易踩坑的地方。
