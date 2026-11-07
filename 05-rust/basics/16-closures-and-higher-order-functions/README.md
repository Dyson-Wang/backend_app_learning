# 16 · Rust · 闭包与高阶函数

## 这个例子讲什么

Rust 闭包如何捕获环境（不可变借用 / 可变借用 / `move` 夺取所有权）、`Fn` / `FnMut` / `FnOnce` 三种 trait 的区别、返回闭包为什么必须 `move`、闭包作为参数与函数指针、以及闭包和迭代器的组合。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/16-closures-and-higher-order-functions
cargo run --quiet
```

## 关键点

- 闭包按「用得最少的权限」自动捕获：只读就 `Fn`，要改就 `FnMut`，要消费就 `FnOnce`。
- `make_counter` 返回 `impl FnMut` 必须 `move`，否则 `count` 是栈上局部变量，返回后失效。
- `move` 闭包会把捕获的变量**所有权**搬进去，之后外部不能再使用该变量。
- 无捕获的闭包可以强转成 `fn` 指针，适合传给 C 风格回调。
- 闭包 + 迭代器是 Rust 最常用的组合，`filter(|&&x| ...)` 的双解引用是因为 `iter()` 产生 `&i32`、`filter` 又传 `&&i32`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 捕获方式 | 按引用（循环变量陷阱） | 词法作用域 | 按引用（`nonlocal`） | 只能捕获 effectively final | 自动推断 / `move` |
| 可变捕获 | 直接改外部变量 | 直接改 | `nonlocal` / `global` | 不允许，需数组或 `Atomic` | `FnMut` |
| 类型表达 | `func(T) R` | 函数值 | `Callable` | 函数式接口 | `Fn` / `FnMut` / `FnOnce` |
| 返回闭包 | 直接返回 | 直接返回 | 直接返回 | 直接返回 | 需 `move` + `impl Fn` |

Go 在 1.22 前循环变量共享导致闭包经典 bug；Java 的 lambda 只能捕获 final 变量，改状态要用 `AtomicInteger` 或单元素数组；Rust 用 `Fn`/`FnMut`/`FnOnce` 在类型层面区分「能不能改、能不能多次调用」。
