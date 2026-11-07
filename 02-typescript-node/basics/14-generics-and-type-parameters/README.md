# 14 · Node.js · 泛型与类型参数（类型擦除）

## 这个例子讲什么

用 JSDoc 的 `@template` 表达泛型函数、带约束的泛型（`{string|number} T`）、多类型参数、泛型容器 `Stack<T>` 与泛型集合 `Map<K,V>`，并给出 TypeScript 的 `<T>` / `<T extends ...>` 等价签名。重点实证：泛型在 JS 运行时被**完全擦除**，一个函数体服务所有类型，零特化、零检查。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/14-generics-and-type-parameters/index.mjs
```

## 关键点

- `node` 直接执行 `.mjs` 时，`@template` / `@param {T}` 只是注释；配 `tsc --checkJs --noEmit` 或 IDE 才做静态类型检查。
- 约束写法是 `@template {string | number} T`，对应 TS 的 `<T extends string | number>`、Java 的 `<T extends Number>`、Rust 的 `T: Trait`。
- 泛型函数只有**一份**函数体：`identity.length` 固定、`toString()` 里没有 `<T>`，传错类型照样运行——这就是类型擦除。
- 泛型容器 `Stack<T>`、`Map<K,V>` 运行时就是普通对象，`T` 不参与任何校验；想运行时校验得自己写 `if`。
- 类型推断是编译期行为：TS 从实参推出 `T`；运行时只能靠 `typeof` / `Array.isArray` 看值本身。
- **单态化 vs 类型擦除**：C++/Rust 为每种 `T` 生成专门代码（单态化，体积换速度）；Java/TS/JS 擦除成同一份（Java 还保留 `List<String>.class == List<Integer>.class` 的实证）。JS 连注解都不保留，最彻底。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js / TS |
|---|---|---|---|---|---|
| 泛型语法 | `[T any]` / `[T comparable]` | `TypeVar` / 3.12 `def f[T]` | `<T>` | `<T: Trait>` | `<T>`（TS）/ `@template`（JSDoc） |
| 约束 | constraint interface | `Protocol` / `bound` | `extends` / `? super` | trait bound | `extends` / `@template {..}` |
| 实现策略 | 单态化 + 字典（GC shape） | 运行时无泛型 | **类型擦除** | **单态化** | **类型擦除** |
| 运行时可查 T | 否 | 否（`__class_getitem__` 仅提示） | 擦除，`instanceof` 不可查 | 否 | 否 |
| 泛型集合 | 手写 / 第三方 | `list[T]` 仅标注 | `List<T>` | `Vec<T>` | `Array<T>` / `Map<K,V>` 标注 |

Java 与 TS 同为类型擦除，但 Java 还有泛型不变性、`? extends/? super` 通配符与桥方法；Rust/Go 走单态化（Go 1.18 用 GC shape stenciling 折中）。JS 因为根本没有运行时类型系统，泛型完全退化为编译期/编辑器提示，连「擦除后的 Class 对象」都不存在。
