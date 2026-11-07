# 04 · Python · 函数与错误处理

## 这个例子讲什么

参数形态（默认值、位置限定 `/`、关键字限定 `*`、`*args`、`**kwargs`）、可变默认参数陷阱、`try/except/else/finally`、自定义异常、`raise ... from` 异常链、`assert` 与 `lambda`。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/04-functions-and-errors/main.py
```

## 关键点

- **默认值只在 `def` 时求值一次**：`def f(x, box=[])` 是经典坑，所有调用共享同一个 list；应写成 `box=None` 再在函数内新建。
- `/` 之前为位置限定参数，`*` 之后为关键字限定参数，能防止调用方误用、便于后续重命名形参。
- `try/except/else/finally`：`else` **仅在无异常时**执行（和 `finally` 不同），放「成功后」逻辑最清晰。
- `raise New from exc` 保留 `__cause__`，回溯里会打印「直接原因」链；异常是类，可带字段（这里 `ValidationError.code`）。
- `assert` 在 `-O` 优化模式下会被整体去掉，别拿它做输入校验；`lambda` 只适合单行表达式。

## 和其他语言的差异

- **Go**：多返回值 + `error` 值、`defer` 延迟清理、`panic/recover`，错误是普通值不是异常。
- **Node.js/TS**：`try/catch/finally`、`throw` 可抛任意值、异步错误走 `Promise.reject` / `async...await`。
- **Java**：受检异常 `checked exception`、`throws` 声明、`try-with-resources`、`finally`。
- **Rust**：`Result<T, E>` / `Option<T>` + `?` 传播、`panic!`，**根本没有异常机制**，编译期强制处理。
- **Python**：`try/except/else/finally`、`raise ... from` 异常链、`Exception` 子类、`*args/**kwargs`、位置/关键字限定参数。
