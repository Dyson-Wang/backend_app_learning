# 16 · Java · 闭包与高阶函数

## 这个例子讲什么

lambda / 方法引用、函数作为值、返回函数（闭包捕获）、`effectively final` 约束、捕获可变状态的坑（final 数组）、闭包生命周期、`this` 绑定、常见高阶函数（`filter` 等）。

## 怎么跑

```
cd backend-roadmap/04-java/basics/16-closures-and-higher-order-functions
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **函数式接口**：`Function`/`Predicate`/`IntUnaryOperator`/`Runnable` 等只含一个抽象方法，lambda 就是它们的实例；`java.util.function` 是标准库自带的函数类型。
- **捕获必须 `effectively final`**：被 lambda 捕获的局部变量要么 `final`，要么「赋值后不再改」。想改状态就把引用保持 final、改它指向的对象（示例用 `int[] holder`），或用 `AtomicInteger`——这是 Java 特有的坑，也是循环里捕获 `i` 报错的根因。
- **闭包捕获的是值/引用副本**，能脱离原栈帧存活（`adder(100)` 返回的函数稍后仍可用）。
- **lambda 的 `this` 是词法捕获的外层对象**，与 JavaScript 普通函数按调用者动态绑定 `this` 不同；箭头函数才和 Java lambda 类似。
- `Stream.filter/map` 接收 `Predicate`/`Function`，就是最常见的高阶函数用法。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 函数类型 | 函数式接口 / `Function` | `func` 一等值 | 一等函数 | 一等函数 | `Fn`/`FnMut`/`FnOnce` trait |
| 捕获 | 必须 effectively final | 按引用捕获变量 | 按引用（词法） | 按引用，`nonlocal` 改外层 | `move` 闭包 / 借用 |
| 可变状态坑 | final 引用 + 可变对象 | 循环变量捕获（1.22 前坑） | `var` 循环变量坑 | late binding 闭包坑 | 所有权/借用检查 |
| this/self | 词法捕获外层 `this` | 无（显式接收者） | 普通函数动态 / 箭头词法 | 显式 `self` | 显式 `self` |
| 高阶函数 | `Stream.filter/map` | 手写 | `Array.map/filter` | `map/filter`/`functools` | `Iterator` 适配器 |
