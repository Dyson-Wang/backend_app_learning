# 06 · Java · 类、接口与模块（结构体）

## 这个例子讲什么

类与封装（private 字段 + getter）、构造器校验、`record` 不可变数据类、`interface`（default / static 方法）、`abstract class` 继承与 `@Override` 多态、可见性规则。

## 怎么跑

```
cd backend-roadmap/04-java/basics/06-structs-classes-and-modules
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **一切皆类**：`record`（Java 16+）是「不可变结构体」，自动生成访问器、`equals/hashCode/toString`，紧凑构造器 `Point { }` 用来校验。
- **接口**方法隐式 `public abstract`，可带 `default` / `static` / `private` 方法；实现类必须把重写方法声明为 `public`。
- 访问修饰符四级：`public > protected > (包私有) > private`；同一顶层类内部可互访 private 成员。
- `abstract class` 可含实现与抽象方法；`@Override` 让编译器帮你检查重写是否正确。
- Java 9 的 **JPMS 模块**（`module-info.java` 的 `exports/requires`）是「模块」的正解，但那是独立文件、需 `--module-path`；本系列单文件示例都在**未命名模块**里，包可见性靠访问修饰符演示。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 数据结构 | `class` / `record` | `struct` | `class` / 对象字面量 | `class` / `dataclass` | `struct` / `enum` |
| 抽象行为 | `interface` | `interface`（隐式实现） | 无（鸭子类型） | ABC / Protocol | `trait` |
| 继承 | 单继承 `extends` | 组合 + 嵌入 | `extends` 原型链 | 多继承 | 无，用 trait / 组合 |
| 不可变数据 | `record` | 无（约定） | `Object.freeze` | `dataclass(frozen=True)` | 字段默认不可变 |
| 可见性 | 四级修饰符 | 首字母大小写导出 | 无（约定 `_`） | 约定 `_` 前缀 | `pub` / `pub(crate)` |
