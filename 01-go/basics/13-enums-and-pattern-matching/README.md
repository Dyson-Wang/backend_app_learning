# 13 · 枚举、模式匹配与代数数据类型

## 这个例子讲什么

Go 没有原生 `enum` 和模式匹配。本例用 `iota` + 自定义类型模拟枚举，用「接口 + 未导出标记方法」构造封闭的代数数据类型（ADT），再用 `type switch` 做分支匹配与解构，并诚实展示「没有穷尽性检查」这一局限。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/13-enums-and-pattern-matching
```

## 关键点

- **iota 枚举**：`const` 组里 `iota` 从 0 递增；配合自定义类型 `type Status int` 和 `String()` 方法才有类型安全和可读输出。`1 << iota` 可做位掩码。
- **封闭接口技巧**：接口里放一个**未导出方法** `sealedExpr()`，包外类型无法实现它，于是实现者被限制在本包内，类型集合近似「封闭」。这是模拟 `sealed`/ADT 的常用手法。
- **type switch 当模式匹配**：`switch v := e.(type)` 按动态类型分支，`v` 已是具体类型，可直接取字段（即「解构」）。
- **没有编译期穷尽性检查**：新增一个 `Expr` 实现类型而忘了改 `describe`，编译器**不会报错**，只能靠 `default` 运行时兜底（本例 `describe(nil)` 演示）。这是与 Rust/Java sealed 最大的差距。
- 没有 `match` 卫语句、没有 `if let`、没有 `let else`、没有解构绑定；type assertion `x.(T)` 的 comma-ok 形式用于安全取具体类型。

## 和其他语言的差异

| 语言 | 枚举 | 带数据 / ADT | 模式匹配与穷尽性 |
|---|---|---|---|
| Go | `iota` + 自定义类型（无关键字） | 接口 + 未导出标记方法 | `type switch`；**无编译期穷尽性检查** |
| Node.js | 无原生 enum；`Object.freeze` 或 TS `enum`/`as const` | 判别联合 `{kind, ...}` | `switch(kind)` + TS `never` 兜底；运行时无类型 |
| Python | `enum.Enum` / `IntEnum` | `dataclass` + `typing.Union` | 3.10+ `match case`；无编译期穷尽性（靠类型检查器） |
| Java | 原生 `enum` | `sealed interface ... permits` + `record` | `instanceof` 模式匹配、`switch` 表达式；穷尽性靠 `sealed` 约束 |
| Rust | 原生 `enum`，可带数据 | `enum` + `Option`/`Result` | `match` **强制穷尽**、`if let` / `let else`、解构 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
