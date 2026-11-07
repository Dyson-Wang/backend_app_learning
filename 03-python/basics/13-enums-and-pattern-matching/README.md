# 13 · Python · 枚举、模式匹配与代数数据类型

## 这个例子讲什么

`enum.Enum` / `IntEnum`、`dataclass` + `typing.Union` 当 ADT、`match case` 的类模式/字面量/捕获/通配/守卫/OR/序列/映射模式，以及穷尽性检查的局限。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/13-enums-and-pattern-matching/main.py
```

## 关键点

- `Enum` 成员是单例，有 `.name`/`.value`；`IntEnum` 成员可直接当 `int` 参与比较与运算。
- **ADT 用 `Union[Circle, Rect]`（或 `Circle | Rect`）+ `dataclass`** 表达；`match` 的**类模式** `case Circle(r=r)` 按类型+属性解构。
- `match case` 支持守卫 `if`、OR `|`、序列 `[a,*rest]`、映射 `{"k": v}` 与值模式 `Color.RED`；`case _` 兜底。
- **没有编译期穷尽性检查**：漏分支运行期才暴露；`typing.assert_never` 是给 `mypy` 的信号，运行时调用会抛 `AssertionError`。
- 与 Rust / Java 的 `sealed` 不同，Python 的 `Union` 不封闭子类集合。

## 和其他语言的差异

- **Go**：无原生 enum，用 `iota` + 自定义类型 + `type switch` 模拟，无模式匹配。
- **Node.js/TS**：`.mjs` 无 enum（TS 的 `enum` / `as const` 是编译期），用判别联合 + `switch` 穷尽性。
- **Java**：`enum` + `sealed interface` + `record` + `instanceof` 模式匹配，穷尽性靠 `sealed` 而非 switch 检查。
- **Rust**：`enum` 带数据 + `match` **编译期强制穷尽** + `if let` / `let else`，`Option`/`Result` 本身就是 ADT。
- **Python**：`enum.Enum`/`IntEnum`、`dataclass` + `Union`、`match case`、`assert_never` 给静态检查器，无编译期穷尽性。
