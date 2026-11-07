# 01 · Rust · hello 与工具链

## 这个例子讲什么

Rust 的第一课不是语法，是三条规则：**默认不可变、值有唯一所有者、借用要么多个不可变要么一个可变**。编译期就把数据竞争和悬垂引用挡掉。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/01-hello-and-toolchain
cargo run

# 建议把注释掉的两行依次打开，亲眼看看编译器的报错——报错本身就是教材
cargo check
```

## 关键点

- `let` 绑定默认不可变，`mut` 是显式 opt-in；这和 Go/Java/Python 的默认可变完全相反。
- `String` 在堆上，赋值会**移动所有权**而不是拷贝；想拷贝要 `.clone()`，想临时用要 `&` 借用。
- 借用检查规则：任意时刻只能有「多个 `&T`」或「一个 `&mut T`」，两者不可共存。这是 Rust 无 GC 还能内存安全的核心。
- **没有 null、没有 GC、没有异常**：`Option<T>` 表达可能为空，`Result<T,E>` 表达可能失败，值离开作用域自动 `Drop`。
- `Cargo.toml` 声明包与依赖，`Cargo.lock` 锁版本；`cargo check` 只做类型检查（快），`cargo build` 才生成二进制。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 语言 | 可变性 | 空值 | 内存 |
|---|---|---|---|
| Rust | 默认不可变 | `Option<T>` | 所有权，编译期检查，无 GC |
| Go | 默认可变 | `nil` | GC |
| Node.js | 默认可变 | `null` / `undefined` | GC |
| Python | 默认可变 | `None` | 引用计数 + GC |
| Java | 默认可变 | `null` | JVM GC |

Rust 的学习曲线主要来自「编译器会拒绝你写惯了的写法」，而不是语法本身。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
