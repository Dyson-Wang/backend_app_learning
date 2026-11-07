# 04 · 函数与错误处理

## 这个例子讲什么

**整套库最值得横向对照的一课。** Go 没有异常：`error` 就是一个普通返回值，调用方必须显式处理。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/04-functions-and-errors
```

## 关键点

- `error` 是内建接口：只有一个 `Error() string` 方法，任何类型实现它就能当错误用。
- **错误是值**：`errors.Is` 判断哨兵错误，`errors.As` 取出具体错误类型；用 `fmt.Errorf("...: %w", err)` 包装才能保留错误链。
- 多返回值让「结果 + 错误」成为惯例，`if err != nil { return }` 满屏但调用链完全显式。
- `panic` 只用于真正不可能恢复的情况（如程序初始化失败），**不要当异常用**；`recover` 只在 `defer` 中生效，常用来做请求级兜底。
- 闭包捕获外部变量，是中间件、函数选项模式的基础。

## 和其他语言的差异

| 语言 | 错误表达 | 调用方是否被迫处理 |
|---|---|---|
| Go | `error` 返回值，显式 `if err != nil` | 是（编译器不强制，但 `_` 忽略会被 lint 抓） |
| Node.js | `throw` / rejected Promise | 否，可能冒泡成未捕获异常 |
| Python | `Exception` | 否 |
| Java | 受检 / 非受检异常 | 受检异常强制声明 |
| Rust | `Result<T, E>` + `?` | 是，`Result` 必须被消费（有 `#[must_use]`） |

Rust 的 `Result` 和 Go 的 `error` 是同一哲学的两种强度：Go 靠约定，Rust 靠类型系统。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
