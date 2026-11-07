# 17 · Rust · 模块进阶：可见性与命名空间

## 这个例子讲什么

Rust 的模块系统：`mod` 定义命名空间、`pub` / `pub(crate)` / 默认私有三级可见性、`use` 导入与别名、嵌套模块与 `super` 路径、以及「私有边界」如何成为 API 设计工具。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/17-modules-and-visibility
cargo run --quiet
```

## 关键点

- 默认一切私有；`pub` 才对 crate 外可见，`pub(crate)` 只在本 crate 内可见——后者适合「内部工具函数」。
- 私有项对**子模块**可见（子模块能看父模块私有项），但父模块看不到子模块的私有项。
- `use ... as ...` 解决命名冲突；`use geometry::Point;` 后可直接写 `Point`。
- `super::` 指向父模块，`self::` 指向当前模块，`crate::` 指向 crate 根。
- 模块可以写在单独文件（`mod foo;` 对应 `foo.rs`），也可以内联（`mod foo { ... }`），本例用内联便于单文件运行。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 可见性单位 | 首字母大小写 | `export` | `_` 前缀约定 | `public` / `private` / `protected` / 包私有 | `pub` / `pub(crate)` / 私有 |
| 包/模块 | package + 目录 | ESM 文件 | 模块 + 包 | package + class | crate + mod |
| 内部可见 | `internal/` 目录约定 | 无 | `_` 约定 | 包私有（默认） | `pub(crate)` |
| 循环依赖 | 编译报错 | ESM 支持但有陷阱 | 运行时报错 | 允许（有坑） | 模块间允许，类型间不允许 |

Go 用目录 `internal/` 做「只能本仓库用」的硬约束，Rust 用 `pub(crate)` 在语言层面表达同一件事；Java 的默认包私有其实和 `pub(crate)` 很像，但没有 crate 这一层；Python 的 `_` 只是约定，运行时完全拦不住。
