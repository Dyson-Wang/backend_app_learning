# 06 · Rust · 结构体 / 类、接口、模块

## 这个例子讲什么

Rust 没有「类」：`struct` 只放数据，`impl` 块挂方法，`trait` 就是接口（可带默认实现），`mod` 划分命名空间。还演示具名 / 元组 / 单元三种结构体，以及 `dyn Trait` 的运行时多态。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/06-structs-classes-and-modules
cargo run
```

## 关键点

- 数据与行为分离：`struct` 定义字段，方法写在 `impl Struct { }` 里；实现接口写在 `impl Trait for Struct { }`。
- `trait` 即接口：`fn area(&self)` 是实现者必须提供的方法，`describe` 带默认实现可直接复用。
- `&self` / `&mut self` / `self` 对应「只读借用 / 可变借用 / 拿走所有权」，由编译器检查。
- `Circle::new(...)` 是无 `self` 的关联函数（类似静态方法 / 构造器）；Rust 没有 `new` 关键字。
- `Vec<&dyn Shape>` 是 **trait 对象**，通过虚函数表在运行时分发，类似 Java 接口引用；代价是一次间接跳转。
- `pub` 控制可见性；跨 crate 访问必须逐层 `pub`。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 数据载体 | `struct` | `class` / 对象字面量 | `class` / `dataclass` | `class` | `struct` |
| 接口 | `interface`（隐式实现） | 无（鸭子类型） | `ABC` / `Protocol` | `interface` | `trait`（显式 `impl`） |
| 默认实现 | 无（嵌入模拟） | 原型链 | 基类方法 | `default` 方法 | trait 默认方法 |
| 构造器 | 手写 `NewXxx` | `constructor` | `__init__` | 构造方法 | 关联函数 `new`（约定） |
| 多态 | 接口隐式 | 原型 / duck typing | duck typing | 接口 / 继承 | `dyn Trait` / 泛型 |
| 命名空间 | package | 模块 (ESM) | package / 模块 | package | `mod` / crate |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
