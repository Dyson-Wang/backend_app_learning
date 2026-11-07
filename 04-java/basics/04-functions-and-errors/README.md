# 04 · Java · 函数与错误处理

## 这个例子讲什么

方法重载、可变参数、递归、`try/catch/finally`、受检异常 vs 非受检异常、自定义异常、`throw/throws`、`try-with-resources`、`Optional`。

## 怎么跑

```
cd backend-roadmap/04-java/basics/04-functions-and-errors
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **受检异常**（继承 `Exception`）编译器强制 `try/catch` 或在签名 `throws` 声明；**非受检异常**（继承 `RuntimeException`）不强制。这是 Java 最有辨识度、也最被其他语言诟病的一点。
- `finally` 一定执行，但现代写法更推荐 `try-with-resources`（`AutoCloseable`）——离开块自动 `close()`，确定性清理。
- `throw` 是抛异常的动作，`throws` 是方法签名上的类型声明。
- `Optional<T>` 表达「可能没有值」，减少 `null` 判断。
- 重载由**编译期参数类型**决定，与返回值无关；可变参数 `T...` 本质就是数组。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 错误类型 | 受检/非受检异常 | `error` 返回值 | `throw` / rejected Promise | `Exception` | `Result<T,E>` |
| 错误传播 | `try/catch` 或 `throws` | `if err != nil { return }` | `try/catch` / `.catch` | `try/except` | `?` 运算符 |
| 资源清理 | `try-with-resources`/`finally` | `defer` | `finally` | `with` / `__exit__` | `Drop` / RAII |
| 可选值 | `Optional<T>` | 多返回值 `(v, ok)` | `undefined` / `null` | `None` | `Option<T>` |
