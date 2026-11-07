# 23 · Java · 测试、基准与文档

## 这个例子讲什么

手写断言（`assertEquals` / `assertTrue`）、表驱动测试、带信息的失败、分支覆盖思路、`System.nanoTime` 基准、Javadoc 文档注释格式，以及 JUnit / JMH 的定位。

## 怎么跑

```
cd backend-roadmap/04-java/basics/23-testing-benchmarks-and-docs
javac Main.java && java Main

# 生成文档：javadoc -d doc Main.java
# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **断言**：失败抛 `AssertionError` 并带上用例名与「期望/实际」，便于定位；示例所有断言都通过，故退出码 0，入口可直接 `java Main` 运行。
- **表驱动**：`record Case` + `List<Case>` 把输入/期望列成表，循环断言，是 Java 里组织多用例的常用手法。
- **覆盖率思路**：把正常分支与异常分支（`square(-1)`）都跑到，就是分支覆盖的雏形；真实项目用 JaCoCo。
- **基准**：`System.nanoTime()` 测固定循环，只把耗时当**差值**示意、不当结论（结果依赖机器，输出里校验和才是稳定值）。生产基准必须用 **JMH**，否则 JIT/GC 会严重干扰。
- **文档**：`/** ... */` 配 `@param/@return/@throws`，用 `javadoc` 生成 HTML；单测标准答案是 **JUnit 5 + `@Test`**。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 测试框架 | JUnit 5（本例手写断言） | `testing` 包 | `node:test` / Jest | `unittest`/pytest | `#[cfg(test)]` 内置 |
| 表驱动 | `record` + `List` | `[]struct{}` 用例表 | 数组用例 | `pytest.mark.parametrize` | 用例数组 + 循环 |
| 基准 | JMH（本例 `nanoTime`） | `Benchmark` + `go test -bench` | `perf_hooks` / bench | `timeit` | `#[bench]`/`criterion` |
| 覆盖率 | JaCoCo | `go test -cover` | `c8`/`--experimental-test-coverage` | `coverage.py` | `cargo tarpaulin` |
| 文档 | Javadoc `/** */` | `//` + `go doc` | JSDoc | docstring + `pydoc` | `///` + `rustdoc` |
