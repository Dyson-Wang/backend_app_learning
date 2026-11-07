# 02 · Rust · 变量、类型、常量、零值

## 这个例子讲什么

Rust 的类型底线：`let` 默认不可变、类型可推断也可显式标注、`const`/`static` 两种常量、**没有隐式零值也没有 null**。未初始化变量直接用是编译错误，想表达「空」只能用 `Option<T>`。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/02-variables-and-types
cargo run
```

## 关键点

- `let x = 5;` 不可变，改值必须 `let mut x`；这与 Go/Java/Python 的默认可变相反。
- **遮蔽（shadowing）**：`let count = "十五";` 重新绑定同名变量并换类型，合法；这和 `mut` 赋值是两回事。
- **零值**：Rust 没有 Go 那种「声明即零值」，局部变量必须先初始化；`i32::default()` / `String::default()` 通过 `Default` trait 拿到零值。
- 数值间**无隐式转换**：`i32` 加 `i64` 编译不过，必须写 `as`。字面量后缀（`8080u16`、`0.75f64`）或类型标注可定型。
- `const` 是编译期内联、无地址；`static` 有固定内存地址、生命周期 `'static`。`i32` 到 `i128`/`usize` 等标量类型全部在编译期确定。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 语言 | 变量声明 | 未初始化 | 空值 | 常量 |
|---|---|---|---|---|
| Go | `var` / `:=` | 自动零值 | `nil` | `const` + `iota` |
| Node.js | `let` / `const` | `undefined` | `null` / `undefined` | `const`（绑定不可变，值可变） |
| Python | 赋值即声明 | 不存在（必须先赋值） | `None` | 无真常量，约定全大写 |
| Java | `int x;` / `var` | 字段有默认值，局部变量须先赋值 | 引用 `null` | `final` |
| Rust | `let` / `let mut` | **编译错误**，无隐式零值 | 无 null，用 `Option<T>` | `const`（内联）/ `static` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
