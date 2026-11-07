# 14 · 泛型与类型参数

## 这个例子讲什么

Go 1.18 起的泛型：泛型函数、泛型类型、自定义 constraint、`comparable`/`any`、类型推断、显式类型实参，以及泛型集合（`Stack[T]`、`Set[T]`、`Pair[A,B]`）。并说明 Go 的「字典 + 形状模板」实现与单态化/类型擦除的区别。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/14-generics-and-type-parameters
```

## 关键点

- **约束即接口**：`any` 是 `interface{}` 别名；`comparable` 允许 `==`；自定义约束写 `[T Number]`，`Number interface { ~int | ~int64 | ~float64 }`，`|` 是类型集并、`~T` 匹配底层类型。
- 约束限制了能对 `T` 做什么：只有 `Number` 里的数值类型支持 `+`，所以 `Sum` 里 `total += x` 才合法。
- **类型推断**：`Sum([]int{...})` 自动推出 `T=int`；推断不出时用 `Sum[int64](...)` 显式指定。
- **泛型类型**：`type Stack[T any] struct{}`，方法接收者要写 `func (s *Stack[T]) Push(v T)`；`Pop` 用 `var zero T` 返回零值。
- 泛型函数不能直接当值传递、不能对类型参数做类型断言、方法不能新增类型参数——这些是 Go 泛型的边界。
- **实现策略**：Go 用「GC shape stenciling + 字典」，按类型的内存形状共享一份机器码、靠运行时字典查方法，介于 Rust 单态化与 Java 类型擦除之间；代价是二进制比擦除大、比全量单态化小。

## 和其他语言的差异

| 语言 | 约束写法 | 运行时实现 | 泛型集合 |
|---|---|---|---|
| Go | `[T any]` / `[T comparable]` / 自定义 `interface{ ~int \| ~float64 }` | GC shape stenciling + 字典（部分单态化） | `Stack[T]` / `Set[T comparable]` |
| Node.js | TS `<T extends X>` / JSDoc `@template` | **类型擦除**，运行时无泛型信息 | `Array<T>` 仅类型标注 |
| Python | `TypeVar` / `Generic` / `Protocol`，3.12 `def f[T]` | 类型提示，运行时不强制（鸭子类型） | `list[T]` 仅标注 |
| Java | `<T extends Number>`、通配符 `? extends`/`? super` | **类型擦除**，`List<String>.class == List<Integer>.class` | `List<T>` |
| Rust | trait bound `T: Ord`、`impl<T: Display>` | **单态化**，每个具体类型生成独立机器码（零成本抽象） | `Vec<T>` / `HashMap<K,V>` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
