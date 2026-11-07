# 14 · Python · 泛型与类型参数

## 这个例子讲什么

对照旧写法 `typing.TypeVar` / `Generic` 与 3.12（PEP 695）新语法 `def first[T](...)` / `class Box[T]` / 内联 bound `T: HasLen`，并用 `Protocol`、`get_type_hints`、`__args__` 实证类型擦除。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/14-generics-and-type-parameters/main.py
```

## 关键点

- **PEP 695（3.12+）**：`def first[T](items: list[T]) -> T:` 无需先声明 `TypeVar`；`class Box[T]:` 直接声明泛型类；`[T: HasLen]` 把约束写在方括号里。本课在 3.13 上**实际跑通**了这些语法。
- 旧写法仍可用：`T = TypeVar("T")` + `class Stack(Generic[T])`；约束可用 `TypeVar("N", int, float)`（值约束）或 `bound=...`（上界）。
- **类型擦除**：`Box[int]` 是 `GenericAlias`，`type(Box[int](1)) is Box` 为 `True`——运行时不为每种 `T` 生成新类，注解只给静态检查器看。
- `Protocol` 是结构化约束，`@runtime_checkable` 后才允许 `isinstance(x, Proto)`；泛型集合直接写 `list[int]` / `dict[str,int]`。
- 对比：Rust 泛型是**单态化**（编译期为每个 T 生成代码），Java/TS/Python 都是**擦除**。

## 和其他语言的差异

- **Go**：`[T any]` / `[T comparable]` + 自定义 constraint interface，1.18 才引入，采用「GC shape stenciling」字典式实例化。
- **Node.js/TS**：`function f<T>(...)`、泛型擦除、运行时完全没有 T，JSDoc `@template` 是给 JS 的替代。
- **Java**：`<T extends ...>`、通配符 `? extends` / `? super`、**类型擦除**（`List<String>.class == List<Integer>.class`）。
- **Rust**：`<T: Trait>`、`impl<T>`、**单态化**、`dyn Trait` 才走动态分发。
- **Python**：`TypeVar`/`Generic`/`Protocol` + PEP 695 `def f[T]`，运行时**类型擦除**、无单态化。
