# 22 · Rust · 反射、注解、宏与元编程

## 这个例子讲什么

Rust **没有传统运行时反射**，但有两条替代路线：运行时用 `std::any::Any` + `TypeId` 做有限的类型判定与 `downcast`；编译期用 `macro_rules!` 和 `derive` 做代码生成。本例把两条路线都跑一遍，并说明边界。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/22-reflection-annotations-and-macros
cargo run --quiet
```

## 关键点

- `Any` 只能回答「是不是某个具体类型」并 `downcast`，**不能按字段名取值、不能遍历字段、不能动态调用方法**。
- `#[derive(Debug, Clone, PartialEq)]` 是过程宏，编译期生成实现，运行时零开销。
- `macro_rules!` 是声明式宏，按模式匹配展开代码；本例用它生成 getter，等价于「编译期反射」。
- `dyn Trait` 是运行时多态，有一次虚表跳转；泛型是编译期单态化，无跳转。这是「动态 vs 静态」的取舍。
- 真正的「注解 + 运行时读取」需要 `serde` / `bevy_reflect` 这类库，标准库不提供。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 运行时反射 | `reflect` 完整支持 | `Reflect` / `Proxy` | `inspect` / `getattr` | `java.lang.reflect` 完整 | 仅 `Any` + `TypeId` |
| 注解/标签 | struct tag | 无（装饰器在 TS） | 装饰器 / `__annotations__` | 注解 + `RUNTIME` 保留 | 无运行时注解 |
| 编译期元编程 | `go:generate` | 无 | 元类 / 装饰器 | 注解处理器 APT | `macro_rules!` / 过程宏 |
| 动态调用 | 反射 + 接口 | 动态属性 | `getattr` | 反射 `invoke` | 不支持 |

Java 的反射最强但也最慢（`setAccessible` 还有模块限制）；Go 的 `reflect` 够用但代码啰嗦；Python 的 `getattr` 最灵活；Rust 故意不给运行时反射，用宏把工作推到编译期，换来零开销和可预测性。
