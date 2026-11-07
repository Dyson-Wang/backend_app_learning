# E06 单测、集成测试与容器化

## 问题复现

测试连真库，CI 一半时间等容器；或者没测试，改一行靠手工点。

## 错误写法

```rust
#[test]
fn create_user() {
    let db = connect_real_db();      // 单测连真库
}
```

## 正确写法

用 trait 抽象依赖，测试传 fake：

```rust
trait Store { fn create(&self, u: &User) -> Result<(), Error>; }
// #[cfg(test)] struct FakeStore;
```

内建测试框架：

```rust
#[test]
fn add_works() { assert_eq!(add(1, 2), 3); }
```

集成测试放 `tests/`，`cargo test` 会自动把它们当独立 crate 编译。

容器化：`cargo build --release` + distroless，注意 `musl` 静态链接。

## 怎么验证

```bash
cargo run --quiet
```

预期：前两条 `PASS`，第三条 `FAIL`（故意），汇总 `2 通过，1 失败`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 测试框架 | `testing` | `node:test` | `unittest` | JUnit | `#[test]` 内建 |
| 集成测试 | `_test` 包 / tag | `node:test` | pytest 标记 | Testcontainers | `tests/` 目录 |
| 容器镜像 | distroless | `node:slim` | `python:slim` | JRE | distroless（最小） |
