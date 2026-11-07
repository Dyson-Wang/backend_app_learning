# 14 · Rust · 泛型与类型参数

## 这个例子讲什么

泛型函数 `largest<T>`、泛型结构体 `Pair<T>`、trait bound（`PartialOrd + Copy`）、条件实现 `impl<T: ...>`、自定义 trait 作为约束、`impl Trait` 返回闭包，以及 Rust 泛型的实现策略——**单态化**。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/14-generics-and-type-parameters
cargo run --quiet
```

## 关键点

- 约束写在 `<T: PartialOrd + Copy>`，多个约束用 `+`；也可以用 `where` 子句换行。
- `impl<T: Display + PartialOrd + Copy> Pair<T>` 表示只有满足约束的 `T` 才有 `cmp_display`，不满足的类型调用会编译报错。
- `impl Trait` 在返回位置表示「返回某个实现了该 trait 的匿名类型」，适合返回闭包或迭代器链。
- **单态化**：`largest::<i32>` 和 `largest::<char>` 会各生成一份机器码，运行时零开销，代价是二进制体积。
- 泛型 + trait bound 是 Rust 的零成本抽象基础：`Iterator`、`Future`、`Fn` 全是这么工作的。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 泛型语法 | `[T any]` | TS `<T>` / JSDoc | `TypeVar` / 3.12 `def f[T]` | `<T>` | `<T>` |
| 约束 | constraint interface | `extends` | `bound=` / `Protocol` | `extends` / `? super` | trait bound |
| 实现策略 | 单态化 + GC shape | 类型擦除 | 运行时无泛型 | **类型擦除** | **单态化** |
| 运行时可查 T | 否 | 否 | 否 | 否 | 否 |

Java 泛型擦除后 `List<String>` 和 `List<Integer>` 是同一个类，不能 `new T[]`、不能 `instanceof T`；Go 1.18 用 GC shape stenciling 折中；Rust 单态化后每种类型都是独立代码，所以能对 `T` 做任何静态优化，但编译时间和体积会涨。
