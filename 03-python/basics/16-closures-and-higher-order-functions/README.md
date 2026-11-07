# 16 · Python · 闭包与高阶函数

## 这个例子讲什么

闭包捕获与 `__closure__` 单元、`nonlocal` 修改可变状态、循环里建闭包的**迟绑定坑**及默认参数修复、函数作为值/返回函数、`map`/`filter`/`reduce`/`sorted(key=)`、装饰器、`functools.partial` 与绑定方法自带 `self`。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/16-closures-and-higher-order-functions/main.py
```

## 关键点

- **闭包**把外层变量装进 `__closure__` 的 cell，变量生命周期随闭包延长；`__closure__[i].cell_contents` 可直接看到捕获值。
- 改捕获变量要 `nonlocal`（不是 `global`）；否则 `count += 1` 会被当成局部赋值而报 `UnboundLocalError`。
- **迟绑定坑**：`[lambda: i for i in range(3)]` 三个闭包共享同一个 `i`，调用时都得到 2；用默认参数 `lambda i=i: i` 在定义时绑定，或用工厂函数生成。
- 装饰器就是「接收函数、返回函数」的高阶函数；`@functools.wraps` 保留 `__name__`/`__doc__`。
- **绑定方法**自带 `self`（`__self__`），回调传 `obj.method` 即可；JS 里要靠箭头函数或 `.bind(this)` 解决 `this` 丢失。

## 和其他语言的差异

- **Go**：闭包按变量捕获，1.22 起 `for` 循环每轮新建变量（旧版同样有迟绑定坑）、`defer` 常配闭包。
- **Node.js/TS**：`var` 的循环捕获坑 → 用 `let` 修复；`this` 绑定靠箭头函数 / `.bind()`。
- **Java**：lambda 只能捕获 `effectively final` 变量、方法引用 `Class::method`、`this` 指向外层实例。
- **Rust**：`Fn` / `FnMut` / `FnOnce` 三类闭包、`move` 关键字转移所有权、生命周期由编译器检查。
- **Python**：闭包 + `nonlocal`、迟绑定坑（默认参数修复）、`__closure__`、绑定方法自带 `self`、`functools.partial`。
