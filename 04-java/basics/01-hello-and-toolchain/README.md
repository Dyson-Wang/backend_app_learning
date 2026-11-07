# 01 · Java · hello 与工具链

## 这个例子讲什么

确认 JDK 版本，看 `var` 局部变量推断、`final`、不可变 `String`、受检异常。Java 这一门的分量不在语法，在后面的分层与事务。

## 怎么跑

```
cd backend-roadmap/04-java/basics/01-hello-and-toolchain
javac Main.java && java Main

# 建议 JDK 21+（虚拟线程正式可用）
# macOS 装：brew install openjdk
```

## 关键点

- **编译成字节码，跑在 JVM 上**：一次编译到处运行，代价是部署必须带 JVM 并调 GC / 堆参数。对照 Go / Rust 的单二进制。
- `var` 只是局部变量类型推断，编译期类型已确定，**不是** Python 那种动态类型。
- **受检异常**：编译器强制处理 `IOException` 一类错误。好处是错误路径不会漏，代价是大量 `try/catch` 与异常声明污染签名——这正是 Go 用返回值 error、Rust 用 `Result` 想避开的东西。
- `String` 不可变，循环拼接要用 `StringBuilder`。
- 本机 `java` 目前是 macOS 占位 stub，没有真实 JDK，示例**未经编译验证**。

## 和其他语言的差异

| 语言 | 类型推断 | 错误机制 | 交付 |
|---|---|---|---|
| Java | `var`（编译期确定） | 受检 / 非受检异常 | 字节码 + JVM |
| Go | `:=` | `error` 返回值 | 单二进制 |
| Node.js | 动态（TS 静态） | `throw` / rejected Promise | 源码 + Node |
| Python | 动态，注解可选 | `Exception` | 源码 + CPython |
| Rust | `let` 推断 | `Result<T,E>` | 单二进制 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
