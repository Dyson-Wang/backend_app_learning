# E03 · Go · 输入校验与统一错误码

## 问题复现

接口收到非法参数时，如果遇到第一个错误就 `return`，前端只能「改一个、再提交、再报下一个」，体验极差。如果错误直接返回 `err.Error()`，客户端拿到的是人话而不是可编程的错误码，无法做国际化或针对性提示。

## 错误写法

```go
if u.Name == "" {
    return fmt.Errorf("name 不能为空") // 只报第一个，且没有错误码
}
```

## 正确写法

- 定义 `FieldError{Code, Field, Message}`，`Code` 给程序判断（`REQUIRED` / `LENGTH` / `FORMAT` / `RANGE`），`Message` 给人看；
- 一次校验收集**全部**错误，返回数组；
- 空值、长度、格式、范围分开判断，`else if` 避免同一字段重复报错。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/E03-validation-error-code
go run .
```

预期：`全空：3 个错误`（name / email / age 各一个 `REQUIRED`）；`短名 + 坏邮箱 + 超龄：3 个错误`；`合法：0 个错误`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 校验库 | `go-playground/validator` | `zod` / `joi` | `pydantic` | Bean Validation | `validator` crate |
| 错误聚合 | 手写切片 | 数组 | 异常列表 | `BindingResult` | `Vec<FieldError>` |
| 「未提供」表达 | 指针 `*int` | `undefined` | `None` | `Integer` null | `Option<i32>` |

Go 的 validator 用 struct tag，零运行时依赖但错误信息要自己映射；Python 的 pydantic 把校验和类型绑定在一起；Java 的 Bean Validation 是注解式，Spring 直接返回 400；Rust 的 `validator` 也是 derive 宏，但错误类型必须自己定义。
