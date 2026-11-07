# 04 · Node.js · 函数与错误处理

## 这个例子讲什么

演示函数声明/表达式/箭头函数、默认参数、rest、解构参数、`arguments` 差异；再演示 `throw` / 自定义 `Error` 子类 / `try/catch/finally` / `Error.cause` 错误链，以及 `assert` 断言。核心结论：JS 用异常而非返回值表示错误，`finally` 一定执行。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/04-functions-and-errors/index.mjs
```

只用 `node:assert/strict`，零依赖。

## 关键点

- 箭头函数没有自己的 `this` / `arguments`（16 课细讲），动态 `this` 场景要用普通 `function`。
- `throw` 可以抛任意值，但应始终抛 `Error` 实例；自定义错误靠 `extends Error` + `name`。
- ES2022 的 `new Error(msg, { cause })` 保留根因，避免层层包装时丢失原始错误。
- `finally` 即使 `try` 里 `return` 也会执行，适合释放资源；`catch {}` 可省略绑定变量。
- `err instanceof ValidationError` 做分类处理，认不出的错误 `throw err` 继续上抛，别吞掉。
- Node 里 `throw` 是同步异常；异步 Promise 的失败走 `reject` / `try { await }`（见 09、20）。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 错误模型 | 多返回值 `val, err` | 异常 `raise/except` | 受检/非受检异常 | `Result<T,E>` / `Option` | `throw` + `try/catch` |
| 自定义错误 | 实现 `error` 接口 | 继承 `Exception` | 继承 `Exception` | 实现 `Error` trait | `extends Error` |
| 错误链 | `fmt.Errorf("%w")` | `raise ... from` | `initCause` | `source()` | `{ cause }` |
| 清理 | `defer` | `finally` / `with` | `finally` / try-with-resources | `Drop` 自动 | `finally` |
| 断言 | 无原生 | `assert` 语句 | `assert`（需开启） | `assert!` 宏 | `node:assert` |

Go 和 Rust 把错误当值显式传递、不靠异常控制流；JS / Python / Java 靠抛出与捕获。Rust 的 `Result` 必须被处理，而 JS 漏 `catch` 的异步拒绝会直接触发 `unhandledRejection`。
