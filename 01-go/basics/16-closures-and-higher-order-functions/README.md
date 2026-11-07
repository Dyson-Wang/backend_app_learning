# 16 · 闭包与高阶函数

## 这个例子讲什么

函数作为一等值、返回函数的闭包、捕获可变状态的坑与生命周期、`Apply`/`Keep`/`Reduce`/`Compose` 等常见高阶函数，以及「回调与 self 绑定」在 Go 里的对应物——方法值与方法表达式。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/16-closures-and-higher-order-functions
```

## 关键点

- **函数是值**：可以赋给变量、当参数传、当返回值。`func(...)` 字面量即闭包。
- **捕获变量生命周期**：`Adder(10)` 返回后 `base` 仍存活——编译器通过**逃逸分析**把它放到堆上；闭包不死，捕获的变量就不被回收。
- **捕获可变状态**：`newCounter` 每次调用产生独立的 `n`；而多个闭包捕获**同一个**外层变量时（本例 `shared`）会互相看到对方的修改，这是共享可变状态的典型坑。
- **循环变量陷阱**：Go 1.21 及以前的 `for` 循环变量是**同一个**，闭包里直接引用会全部拿到最终值；惯用修复是 `i := i` 显式遮蔽。Go 1.22 起语言改为每次迭代新建变量，但本仓库 `go.mod` 是 1.21，所以保留这行。
- **self 绑定**：Go 没有 `this`。`g.Greet` 是**方法值**（已绑定接收者，类型 `func(string) string`），`Greeter.Greet` 是**方法表达式**（类型 `func(Greeter, string) string`，调用时显式传接收者）。这正是「回调要不要绑 self」的 Go 答案。
- 没有装饰器、没有 `nonlocal`（直接改捕获变量即可），函数没有 `.bind/.call/.apply`。

## 和其他语言的差异

| 语言 | 闭包 / 捕获 | self / this 绑定 | 高阶函数 |
|---|---|---|---|
| Go | 闭包捕获变量，逃逸到堆；直接改外层变量 | 方法值 `g.M`（已绑定）/ 方法表达式 `T.M` | 手写 `Apply`/`Reduce`；无内建 `map/filter` |
| Node.js | 闭包 + 作用域链；`let`/`const` 块级作用域 | `this` 动态绑定；箭头函数继承 `this`；`.bind/.call/.apply` | `Array.map/filter/reduce` |
| Python | 闭包；改外层需 `nonlocal` | 实例方法自动传 `self`；描述符/`MethodType` | `map/filter`、`functools.reduce`、装饰器 |
| Java | lambda 只能捕获 **effectively final** 变量 | `this` 指向外层实例；方法引用 `Class::method` | `Stream` / `Function` / `Consumer` |
| Rust | `Fn`/`FnMut`/`FnOnce` 三 trait，`move` 转移所有权 | 方法第一个参数是 `&self`/`self` | 迭代器适配器 `map/filter/fold` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
