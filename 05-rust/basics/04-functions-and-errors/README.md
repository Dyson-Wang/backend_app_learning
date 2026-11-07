# 04 · Rust · 函数与错误处理

## 这个例子讲什么

Rust 函数靠「最后一个表达式」返回，参数必须写类型。核心是没有异常机制：可失败返回 `Result<T, E>`，可能为空返回 `Option<T>`，`?` 负责把错误向上传播，调用方用 `match` 显式处理。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/04-functions-and-errors
cargo run
```

## 关键点

- `fn add(a: i32, b: i32) -> i32 { a + b }`：`a + b` 无分号是返回值；加分号就变成语句返回 `()`。
- `?` 只能用在返回 `Result`/`Option` 的函数里，等价于「Ok 就解包、Err 就 `return Err(From::from(e))`」。
- `Result` / `Option` 都是**普通枚举值**，编译器强制你处理 `Err` / `None`（`unwrap` 会 panic，教学外慎用）。
- `unwrap_or` / `map` / `and_then` 等组合子让错误链路保持函数式、不抛异常。
- 自定义错误通常实现 `std::error::Error` + `Display`，或用 `thiserror` 派生（本例零依赖，用库自带的 `ParseIntError`）。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 错误处理 | 多返回值 `(v, err)` | `throw` / `try..catch` 或 rejected Promise | `try..except` 异常 | `try..catch` 受检/非受检异常 | `Result<T,E>` 值 + `?` |
| 可能为空 | `nil` / 零值 | `null` / `undefined` | `None` | `null` / `Optional` | `Option<T>` |
| 提前返回 | `return` | `return` / `throw` | `return` / `raise` | `return` / `throw` | `?` / `return` |
| 函数作为值 | 支持 | 一等公民 | 一等公民 | 方法引用 / lambda | 函数指针 / `Fn` trait |
| 编译器强制处理 | 否（err 可忽略） | 否 | 否 | 受检异常强制 | 是（`#[must_use]`） |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
