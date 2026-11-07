# E03 · Rust · 输入校验与统一错误码

## 问题复现

接口收到非法参数时返回 `Err("参数错误".to_string())`，客户端只知道「错了」：不知道哪个字段、不知道是必填还是格式问题，也无法做多语言。

## 错误写法

```rust
if req.name.is_empty() {
    return Err("参数错误".into()); // 只报第一个，且没有错误码
}
```

## 正确写法

- `FieldError { code: Code, field, message }`，`Code` 用 enum，`as_str()` 给外部；
- 一次收集全部错误，返回 `Vec<FieldError>`；
- `Option<T>` 表达「未提供」，避免 `0` 被误判；
- `match` 三分支处理 age，`Some(_) => {}` 显式表示合法。

## 怎么验证

```bash
cd backend-roadmap/05-rust/advanced/E03-validation-error-code
cargo run --quiet
```

预期：`全空：3 个错误`；`短名 + 坏邮箱 + 超龄：3 个错误`；`合法：0 个错误`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 校验库 | `validator` tag | `zod` / `joi` | `pydantic` | Bean Validation | `validator` crate |
| 错误聚合 | 切片 | 数组 | 列表 | `BindingResult` | `Vec<FieldError>` |
| 「未提供」 | `*int` | `undefined` | `None` | `Integer` null | `Option<i32>` |

Rust 的 `Option` 强制你区分「没传」和「传了 0」，这是它比动态语言更安全的地方；`validator` crate 用 derive 宏生成校验，但错误类型要自己定义。
