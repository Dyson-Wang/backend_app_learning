# 22 · Java · 反射、注解与元编程

## 这个例子讲什么

`java.lang.reflect` 运行时类型信息、按名字访问字段 / 方法、动态调用、自定义 `@interface` 注解（`RetentionPolicy.RUNTIME`）与反射读取，以及 APT / 宏的边界。

## 怎么跑

```
cd backend-roadmap/04-java/basics/22-reflection-annotations-and-macros
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **RTTI**：`Class<?>`、`getDeclaredMethods()`、`getSuperclass()` 在运行期暴露类型结构；这是 Spring / Jackson 等框架的底座。
- **按名字访问**：`getDeclaredField` / `getDeclaredMethod` + `setAccessible(true)` 可突破 `private`；`Field.getInt/setInt`、`Method.invoke` 做动态读写/调用。
- **注解**：`@Retention(RUNTIME)` 才会保留到运行期被 `getAnnotation` 读到；`@Target` 限定作用位置，`default` 给默认值。CLASS/SOURCE 级注解运行期读不到。
- **代价与边界**：反射绕过编译期类型检查、有性能开销，JPMS 下还受 `opens` 限制；不宜在热路径滥用。
- **Java 没有宏**：运行期元编程靠反射，编译期代码生成靠 **APT（注解处理器）**；Lombok 就是 APT 的典型，不是运行期魔法。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 运行时类型信息 | `java.lang.reflect` + `Class` | `reflect` 包 | `Reflect`/`Proxy` | `inspect`/`__annotations__` | 基本无（`Any`/`TypeId`） |
| 注解/标签 | `@interface`（RUNTIME 可读） | struct tag（字符串） | 装饰器（TS）/ JSDoc | 装饰器 / `__annotations__` | `#[attr]` 属性宏 |
| 编译期生成 | APT（Lombok/MapStruct） | `go:generate` | 无（Babel/TS 转换） | 装饰器/元类（运行期） | `macro_rules!`/proc-macro |
| 宏 | 无 | 无 | 无 | 无（有装饰器） | **有，编译期宏** |
| 动态调用代价 | `Method.invoke`（反射开销） | `reflect.Value.Call` | `Proxy` trap | `getattr` | 静态分发，无反射 |
