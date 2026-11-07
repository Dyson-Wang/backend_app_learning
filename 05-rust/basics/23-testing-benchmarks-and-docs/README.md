# 23 · Rust · 测试、基准与文档

## 这个例子讲什么

Rust 的工程化三件套：`#[test]` 单元测试（放在 `#[cfg(test)] mod tests` 里）、表驱动测试与断言消息、`Instant` 手写基准，以及 `///` 文档注释（`cargo doc` 生成 HTML、`cargo test` 会跑文档示例）。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/23-testing-benchmarks-and-docs
cargo run --quiet      # 入口演示
cargo test --quiet     # 跑 #[test] 与文档示例
cargo doc --no-deps --open   # 生成文档
```

## 关键点

- `#[cfg(test)]` 让测试模块只在 `cargo test` 时编译，不进入发布二进制。
- `assert_eq!` 失败会打印左右值；自定义断言可以补上「哪个用例」的上下文。
- 表驱动测试用一个数组驱动多组输入，比复制粘贴多个 `#[test]` 好维护。
- `cargo bench` 需要 nightly 或 `criterion` crate；标准库的 `Instant` 只能做粗粒度计时。
- `///` 是文档注释，代码块会被 `cargo test` 当作 doc test 执行——文档里的示例不会腐烂。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 测试框架 | `testing` + `_test.go` | `node:test` | `unittest` / `pytest` | JUnit | 内建 `#[test]` |
| 断言 | `t.Errorf` | `assert` | `assert` / `assertEqual` | `assertEquals` | `assert_eq!` |
| 基准 | `Benchmark` | `benchmark` / `perf_hooks` | `timeit` | JMH | `cargo bench` / criterion |
| 文档 | `go doc` | JSDoc | docstring / Sphinx | Javadoc | `///` + `cargo doc` |

Rust 把测试和文档都做进了 `cargo`：`cargo test` 同时跑单元测试、集成测试和文档示例，`cargo doc` 直接从注释生成站点，这是它工程体验最好的部分之一。
