# E03 · Java · 输入校验与统一错误码

## 问题复现

接口收到非法参数时抛 `IllegalArgumentException("参数错误")`，客户端只知道「错了」：不知道哪个字段、不知道是必填还是格式问题，也无法做多语言。

## 错误写法

```java
if (body.getName() == null) throw new IllegalArgumentException("参数错误");
```

只报第一个、没有字段、没有错误码。

## 正确写法

- `record FieldError(Code code, String field, String message)`，`Code` 用 enum；
- 一次收集全部错误，返回 `List<FieldError>`；
- `Integer age` 用 null 表示未提供，避免 0 被误判；
- 生产用 Bean Validation（`@NotBlank` / `@Size` / `@Email` / `@Min` / `@Max`），Spring 自动转 400。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/E03-validation-error-code
javac Main.java && java Main
```

预期：`全空：3 个错误`；`短名 + 坏邮箱 + 超龄：3 个错误`；`合法：0 个错误`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 校验库 | `validator` tag | `zod` / `joi` | `pydantic` | Bean Validation | `validator` crate |
| 错误聚合 | 切片 | 数组 | 列表 | `BindingResult` | `Vec<FieldError>` |
| 「未提供」 | `*int` | `undefined` | `None` | `Integer` null | `Option<i32>` |

Java 的 Bean Validation 是声明式：注解写在 DTO 上，`@Valid` 触发，`BindingResult` 收集全部错误；比手写 if 简洁，但错误码映射要自己配 `MessageSource`。
