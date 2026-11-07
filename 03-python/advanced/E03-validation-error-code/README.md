# E03 · Python · 输入校验与统一错误码

## 问题复现

接口收到非法参数时抛 `ValueError("参数错误")`，客户端只知道「错了」：不知道哪个字段、不知道是必填还是格式问题，也没法做多语言。

## 错误写法

```python
if not body.get("name"):
    raise ValueError("参数错误")
```

只报第一个、没有字段、没有错误码。

## 正确写法

- 统一结构 `{code, field, message}`，`code` 用 `enum.Enum` 定义；
- 一次收集全部字段错误；
- `age is None` 单独判必填，避免 `0` 岁被误判；
- 注意 `bool` 是 `int` 的子类，`True` 会被当成 1，所以显式排除。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/E03-validation-error-code
python3 main.py
```

预期：`全空：3 个错误`；`短名 + 坏邮箱 + 超龄：3 个错误`；`合法：0 个错误`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 校验库 | `validator` tag | `zod` / `joi` | `pydantic` | Bean Validation | `validator` crate |
| 错误聚合 | 切片 | 数组 | 列表 | `BindingResult` | `Vec<FieldError>` |
| 「未提供」 | `*int` | `undefined` | `None` | `Integer` null | `Option<i32>` |

Python 的 `pydantic` 把校验和类型绑定，FastAPI 里几乎零成本；但要注意 `bool`/`int` 的坑，以及 `None` 与「缺字段」在 `dict.get` 下是同一个值。
