# 14 · Java · 泛型与类型参数

## 这个例子讲什么

泛型方法 / 泛型类、上界约束 `<T extends ...>`、类型推断、通配符 `? extends` / `? super`（PECS）、泛型集合，以及**类型擦除**的运行期实证。

## 怎么跑

```
cd backend-roadmap/04-java/basics/14-generics-and-type-parameters
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **类型擦除**：Java 泛型只在编译期做类型检查，运行期 `List<String>` 与 `List<Integer>` 都擦成同一个 `List.class`。实证：`strings.getClass() == ints.getClass()` 为 `true`。因此 `List<String>.class` 非法、不能 `new T[]`、不能用 `instanceof List<String>`。
- 这与 **Rust 的单态化**（`impl<T>` 为每个 T 生成一份机器码）正好相反：Java 省代码体积、无特化开销，但丢掉了运行期类型信息。
- **PECS**：`? extends Number` 是生产者（能安全地按 `Number` 读），`? super Integer` 是消费者（能安全地往里写 `Integer`）。
- 上界 `<T extends Comparable<T>>` 让泛型内部能调用 `compareTo`；`<T extends Number>` 能调 `doubleValue()`。
- 泛型方法靠实参推断 `T`（`first(List.of("a"))` 返回 `String`），属编译期推断，不是运行期反射。

## 和其他语言的差异

| 能力 | Java | Go | Node.js (TS) | Python | Rust |
|---|---|---|---|---|---|
| 泛型实现 | **类型擦除** | 单态化（GC shape/stencil） | **类型擦除**（编译后无类型） | 运行期无约束（typing 仅提示） | **单态化** |
| 约束/界 | `<T extends X>` | `[T Constraint]` | `extends` / `keyof` | `Protocol` / `Bound` | `T: Trait` |
| 通配符 | `? extends` / `? super` | 无 | 条件类型 / `in`/`out` | 无 | `impl Trait` |
| 运行期泛型信息 | 擦除，无 | 保留部分 | 无 | `Generic` 可保留 `__orig_class__` | 特化，无 |
| 泛型集合 | `List<T>` | `[]T` | `Array<T>` | `list[T]`（注解） | `Vec<T>` |
