# Java · 04-java

**运行时模型**：JVM · 线程池 / 虚拟线程 · GC · 静态强类型 · 字节码 + 运行时

## 环境准备

本机 `java` 是 macOS 的占位 stub，**没有真实 JDK**，示例尚未编译验证。装 JDK 后：

```bash
brew install openjdk     # 建议 JDK 21+，虚拟线程已正式可用
cd 04-java/basics/01-hello-and-toolchain
javac Main.java && java Main
```

`basics/` 全部只用 JDK 标准库（`com.sun.net.httpserver` / `java.sql` 等），不需要 Maven。`advanced/` 引入 Spring Boot 时会单独给 `pom.xml`。

## 目录

- `basics/` —— 基础入门 10 课（统一编号）
- `advanced/` —— 生产场景（统一 A01 ~ F03 编号）

## 重点补什么

最重的一套工程化：分层架构、依赖注入、声明式事务、异常体系、可维护性。学它的价值不在于语法，而在于理解「一个要活十年的系统靠什么约束复杂度」。

## 与其他语言的对照

Go 用组合和显式 error，Java 用继承和异常；Go 的 goroutine 是语言级轻量协程，Java 21 的虚拟线程是 JVM 级解法——这组对照最能说明「同一个问题，运行时和语言层各自能做到哪一步」。

## 跨语言对照

每课 / 每个场景做完，把「这门语言和其他语言的差异」补到 `../99-cross-language/scenario-matrix.md`。同一编号横向打开，是这套库唯一正确的读法。
