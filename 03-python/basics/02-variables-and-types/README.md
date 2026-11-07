# 02 · Python · 变量、类型、常量、零值

## 这个例子讲什么

用 `type()` / `bool()` / `is` 观察 Python 的动态类型：名字只是绑定到对象的标签、空值只有 `None`、`int` 任意精度、`bool` 是 `int` 子类、`Final` 与全大写只是常量约定。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/02-variables-and-types/main.py
```

## 关键点

- **无零值**：Python 变量必须先赋值，未初始化名会 `NameError`；空值统一是 `None`（单例），没有 Go 那种「每种类型都有确定零值」。
- **常量不强制**：`Final[int]` 与 `UPPER_CASE` 只给 `mypy` / IDE 看，运行时改了也不报错；类型注解同理，不参与运行时校验。
- **动态但强类型**：变量能换类型，但 `1 + "2"` 会 `TypeError`，不会像 JS 那样隐式拼接。
- `int` 任意精度不会溢出；`float` 是双精度，超大才变 `inf`。
- `is` 比身份、`==` 比值；`bool` 是 `int` 子类，所以 `True + True == 2`。

## 和其他语言的差异

- **Go**：`var` / `:=`、`const` + `iota`、静态类型，每个类型有确定零值（`0` / `""` / `nil`）。
- **Node.js/TS**：`let` / `const`（`const` 绑定真不可改）、`undefined` 与 `null` 两种空值、TS 的 `number` 只有双精度浮点。
- **Java**：`final`、字段有默认值（`0` / `null`）但局部变量必须显式初始化、`var`（Java 10+）、包装类型。
- **Rust**：`let` / `let mut`、`const` / `static`、默认绑定不可变、没有 `null`（用 `Option<T>`）、编译期静态类型。
- **Python**：无真常量、`None` 单一空值、`int` 任意精度、`bool` 是 `int` 子类、注解不强制。
