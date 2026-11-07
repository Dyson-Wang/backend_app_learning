# 07 · Rust · 包管理与依赖

## 这个例子讲什么

讲清 Cargo 的三层概念：**package（一个 `Cargo.toml`）→ crate（编译单元）→ module（`mod` 命名空间）**，以及 `crate::` / `self::` / `super::` 路径、`pub use` 重导出、`cfg!` 与 `env!` 编译期注入。本课 `[dependencies]` 故意留空，只依赖 `std`。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/07-package-management
cargo run

# 常用命令：cargo check（只类型检查）/ cargo build / cargo run / cargo test
# 加依赖是在 [dependencies] 写 `serde = "1"`，Cargo 自动解析并写 Cargo.lock
```

## 关键点

- `Cargo.toml` 描述 package 与依赖；`Cargo.lock` 锁定精确版本（binary 应提交，library 不必）。
- `std` 是一个隐式链接的外部 crate；`use std::collections::BTreeSet;` 就是一次「依赖引入」。
- 路径：`crate::` 从 crate 根开始，`self::` 当前模块，`super::` 父模块；`pub use` 可把深层项重导出，缩短调用方路径。
- `cfg!(...)` 是编译期布尔（按 target / feature 注入），`env!("CARGO_PKG_NAME")` 在编译期读取 Cargo 注入的环境变量。
- 依赖解析是**编译期**完成、静态链接进二进制，和 Node/Python 的运行时加载完全不同。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 清单文件 | `go.mod` | `package.json` | `pyproject.toml` / `requirements.txt` | `pom.xml` / `build.gradle` | `Cargo.toml` |
| 锁文件 | `go.sum` | `package-lock.json` | 无（或 `poetry.lock`） | 无标准 | `Cargo.lock` |
| 包 / 模块 | `package` + `import` | ESM `import` / npm 包 | `import` / 包 | `package` / Maven 坐标 | `crate` / `mod` |
| 工具 | `go` | `npm` / `pnpm` | `pip` / `uv` | Maven / Gradle | `cargo` |
| 依赖链接 | 静态（编译进二进制） | 运行时 `node_modules` | 运行时 `import` | 运行时 classpath / jar | 编译期静态链接 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
