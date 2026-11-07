# 13 · Rust · 枚举、模式匹配与代数数据类型

## 这个例子讲什么

Rust 的 `enum` 不是「一组常量」，而是**带数据的和类型（sum type）**：每个变体可以携带不同字段。配合 `match` 的穷尽性检查、`if let`、`let else`、`@` 绑定和 `matches!`，构成 Rust 表达业务状态的主力工具。`Option<T>` / `Result<T, E>` 就是标准库里的 enum。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/13-enums-and-pattern-matching
cargo run --quiet
```

## 关键点

- `match` 必须覆盖所有变体，漏一个直接编译错误——这是「新增状态忘了处理」这类线上事故的编译期拦截。
- 变体可以带命名字段（`Circle { radius }`）、元组字段（`Triangle(f64, f64)`）或什么都不带。
- 守卫 `if *radius > 10.0` 让匹配带条件；`{ .. }` / `(..)` 忽略不需要的字段。
- `if let` 处理单分支，`let else` 处理「不匹配就退出」，后者比嵌套 `if let` 少一层缩进。
- `@` 绑定同时做范围匹配和取值，`matches!` 把匹配结果转成 `bool`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 枚举 | `iota` + 常量，无数据 | 无原生，`Object.freeze` / TS union | `enum.Enum`，无数据 | `enum` 可带字段和方法 | `enum` 带数据的和类型 |
| 模式匹配 | `type switch`（仅类型） | `switch`（值） | `match case`（结构化） | `instanceof` + `switch`（17 无模式匹配） | `match` 全功能 |
| 穷尽性检查 | 无 | 无 | 无（需 `assert_never`） | `sealed` 提示，编译器不强制 | **编译器强制** |
| 空值表达 | `nil` | `null` / `undefined` | `None` | `null` / `Optional` | `Option<T>` |

Go 的 `type switch` 只能按具体类型分支，没有解构；Python 3.10 的 `match case` 很接近但没有穷尽性检查；Java 靠 `sealed interface` + `instanceof` 手动保证；Rust 是唯一把「漏处理变体」变成编译错误的。
