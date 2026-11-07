# 16 · Node.js · 闭包与高阶函数

## 这个例子讲什么

演示闭包捕获、循环里 `var`/`let` 捕获可变状态的经典坑、被捕获变量的生命周期；再把函数当值传递/返回，实现 `compose`/`once`/`memoize` 等常见高阶函数，最后讲清回调里普通函数丢 `this` 与箭头函数词法绑定 `this` 的区别（`bind`/`call`/`apply`）。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/16-closures-and-higher-order-functions/index.mjs
```

## 关键点

- 闭包 = 函数 + 其词法环境；外层函数返回后，被捕获的变量仍存活，直到该闭包不再被引用（GC 才能回收，21 课）。
- `var` 是函数/模块作用域，循环里所有闭包共享同一个 `i`，结果都是终值；`let` 每轮迭代创建新绑定，结果才是 `0,1,2`。
- 每次调用 `makeCounter()` 生成**独立**的闭包环境，互不干扰；这正是工厂函数与私有状态的实现基础。
- 高阶函数：接收函数（`map/filter/reduce`）或返回函数（`compose`/`once`/`memoize`）；`memoize` 用闭包里的 `Map` 缓存结果。
- `this` 绑定：普通 `function` 的 `this` 由**调用方式**决定，把 `obj.scale` 当裸回调传会丢 `this`（严格模式下为 `undefined` 而抛 `TypeError`）。
- `箭头函数`没有自己的 `this`，从定义处词法捕获；`obj.makeMultiplier()` 里返回的箭头就固定绑到 `obj`。
- `fn.bind(obj)` / `fn.call(obj, ...)` / `fn.apply(obj, [])` 可显式指定 `this`；事件回调、数组方法回调优先用箭头。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 闭包 | 匿名函数捕获变量 | `lambda`/嵌套函数 | lambda / 匿名类 | 闭包 `Fn`/`FnMut`/`FnOnce` | 箭头/普通函数 |
| 捕获可变状态坑 | 循环变量捕获（1.22 前） | 后期绑定（late binding） | 需 effectively final | `move` 显式转移 | `var` 共享绑定 |
| 函数作为值 | `func` 类型 | 一等对象 | 函数式接口 / 方法引用 | `Fn` trait / 函数指针 | 一等对象 |
| 常见高阶函数 | 手写 | `map`/`filter`/`functools` | Stream API | `Iterator` 适配器 | 数组 `map/filter/reduce` |
| this/self 绑定 | 无 this，显式接收者 | 显式 `self` | `this` 隐式 | 显式 `self` | 动态 `this` + 词法箭头 |

Python 也有「循环里闭包后期绑定」的同类坑（默认参数 `x=x` 或 `functools.partial` 解决），与 JS 的 `var` 坑几乎一样。Rust 用 `Fn/FnMut/FnOnce` 三 trait 在编译期约束闭包能否改捕获状态、`move` 显式转移所有权，是这里最严格的。Java lambda 要求捕获变量 effectively final，从根上避开了共享可变绑定的坑；JS 的 `this` 动态绑定则是其余四门都没有的独特复杂度。
