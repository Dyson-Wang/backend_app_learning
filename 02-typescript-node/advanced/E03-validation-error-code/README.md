# E03 · Node.js · 输入校验与统一错误码

## 问题复现

接口收到非法参数时，如果抛一个 `Error('参数错误')`，客户端只能拿到一句人话：不知道是哪个字段、不知道是必填还是格式问题，也没法做多语言。

## 错误写法

```js
if (!body.name) throw new Error('参数错误');
```

只报第一个、没有字段、没有错误码。

## 正确写法

- 统一结构 `{ code, field, message }`；
- `code` 用 `Object.freeze` 的常量（生产用 TypeScript 联合类型或 `as const`），客户端可 switch；
- 一次收集全部字段错误，返回数组；
- 空值、长度、格式、范围分开判断。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/E03-validation-error-code
node index.mjs
```

预期：`全空：3 个错误`（name / email / age 各一个 `REQUIRED`）；`短名 + 坏邮箱 + 超龄：3 个错误`；`合法：0 个错误`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 校验库 | `validator` tag | `zod` / `joi` | `pydantic` | Bean Validation | `validator` crate |
| 错误聚合 | 切片 | 数组 | 异常列表 | `BindingResult` | `Vec<FieldError>` |
| 「未提供」表达 | `*int` | `undefined` | `None` | `Integer` null | `Option<i32>` |

Node 生态里 `zod` 是事实标准：schema 即类型，校验和 TS 类型一次定义；Java 的 Bean Validation 靠注解，Spring 自动转 400；Python 的 pydantic 在 FastAPI 里几乎是默认选择。
