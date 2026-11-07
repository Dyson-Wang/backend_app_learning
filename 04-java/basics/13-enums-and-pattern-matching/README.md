# 13 · Java · 枚举、模式匹配与代数数据类型

## 这个例子讲什么

带字段的 `enum`、`enum` 上的 `switch` 表达式、`sealed interface + permits` 封闭继承、`record` 作为带数据的变体、`instanceof` 模式匹配解构，以及「穷尽性」的讨论。

## 怎么跑

```
cd backend-roadmap/04-java/basics/13-enums-and-pattern-matching
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **`enum` 是真类型**，可带字段 / 构造器 / 方法；`values()`、`valueOf()`、`ordinal()`。`switch` 表达式作用于枚举是**稳定特性**。
- **`sealed` + `permits`**（Java 17）把实现集合封闭起来，配合 `record`（隐式 `final`）就是 Java 的 ADT：一个 `Shape` 只能是 `Circle` 或 `Rect`。
- **解构用 `instanceof` 模式匹配**（Java 16 稳定）：`s instanceof Circle c` 之后直接用 `c.r()`  accessor。**record patterns 与 switch 模式匹配在 17/21 仍是 preview，本系列不用 `--enable-preview`**。
- **穷尽性靠 `sealed` 约束「有哪些实现」，而不是靠编译器对 `switch` 做穷尽检查**；`if/else` 链编译器不验证穷尽，所以末尾必须自己 `throw` 兜底。
- 对照 Rust 的 `match` 编译期强制穷尽，Java 这里要弱一些。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 枚举 | 原生 `enum`（带字段） | `iota` + 自定义类型 | 无（`Object.freeze`/`as const`） | `enum.Enum` | `enum` 可带数据 |
| ADT | `sealed`+`record` | `interface` + type switch | 判别联合 + `as const` | `dataclass`+`Union` | `enum` 变体 |
| 模式匹配 | `instanceof` 模式（switch 仅 preview） | `type switch` | 无（TS 有 `switch` 收窄） | `match case`（3.10+） | `match` / `if let` |
| 穷尽检查 | 靠 `sealed` 约束实现集合 | 无（需 `default`） | 需 `never` 兜底 | 无 | **编译期强制穷尽** |
