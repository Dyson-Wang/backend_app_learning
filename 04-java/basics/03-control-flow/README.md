# 03 · Java · 控制流：分支、循环、迭代

## 这个例子讲什么

if/else、三元表达式、`switch` 表达式（箭头、无穿透）、for / 增强 for / while / do-while、break / continue 与带标签 break。

## 怎么跑

```
cd backend-roadmap/04-java/basics/03-control-flow
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- `switch` 表达式（Java 14+）用 `->`，**不会穿透**到下一个 case，可直接返回值，且必须写 `default`；传统 `:` 语句形式才会穿透。
- `case 3, 4, 5 ->` 支持多个标签合并。
- **增强 for** 适用于任何 `Iterable` 和数组，底层就是 `Iterator`。
- Java 没有 `goto`，但 `break 标签` / `continue 标签` 能精确跳出多层循环。
- do-while 至少执行一次；三元 `? :` 是表达式，有返回值。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 循环 | for/while/do-while | 只有 `for` | for/while/for-of | for/while | loop/while/for |
| 多分支 | `switch` 表达式/语句 | `switch`（默认不穿透） | `switch` | `match`（3.10+） | `match`（穷尽） |
| 穿透 | 传统 case 穿透，箭头不穿透 | 默认不穿透 | 需 `break` | 无 fallthrough | 无 fallthrough |
| 多层跳出 | `break 标签` | `break 标签` | 标签 break | 无（用异常/标志） | `'label: loop` |
| 迭代协议 | `Iterable`/`Iterator` | `range` | `for-of`/`Symbol.iterator` | `__iter__`/`__next__` | `Iterator` trait |
