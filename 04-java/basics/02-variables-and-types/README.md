# 02 · Java · 变量、类型、常量与零值

## 这个例子讲什么

8 种基本类型、`var` 局部变量推断、`final` 常量、字段零值 vs 局部变量无默认值、自动装箱/拆箱、隐式加宽与显式窄化转换。

## 怎么跑

```
cd backend-roadmap/04-java/basics/02-variables-and-types
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **基本类型不是对象**：`int/byte/...` 是值，`String/Integer/...` 是引用；包装类用来把基本类型塞进泛型集合。
- **`var` 只用于局部变量**，编译期类型已确定，不是动态类型。
- **字段有零值**（`0` / `false` / `'\u0000'` / `null`），**局部变量没有默认值**，未初始化就使用会编译报错——比脚本语言严格。
- **整数默认 `int`、浮点默认 `double`**，写 `long` 要加 `L`、`float` 要加 `f`。
- 窄化转换要显式 cast，且会静默截断/溢出，不抛异常。

## 和其他语言的差异

| 维度 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 常量 | `final` / `static final` | `const` | `const` | 约定全大写（无真常量） | `const` / `static` |
| 类型推断 | `var`（局部变量） | `:=` | `let/const`（动态） | 动态 | `let` 推断 |
| 零值 | 字段有零值，局部变量无 | 所有变量有零值 | `undefined` | 名字必须绑定 | 无，必须初始化 |
| 数值溢出 | 静默回绕 | 静默回绕 | Number 双精度 | 任意精度 `int` | release 回绕 / debug panic |
| 装箱 | 自动 `Integer` | 接口 `any` | 无 | 一切皆对象 | `Box<T>` 显式 |
