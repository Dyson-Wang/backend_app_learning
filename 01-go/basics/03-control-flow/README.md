# 03 · 控制流：分支、循环、迭代

## 这个例子讲什么

Go 只有 `for` 一种循环、`switch` 默认不穿透、`defer` 后进先出。语法少，但组合方式需要适应。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/03-control-flow
```

## 关键点

- **没有 while / do-while**，`for condition {}` 就是 while，`for {}` 是死循环。
- `switch` 自带 break，需要穿透必须显式写 `fallthrough`；也支持无表达式的 `switch { case cond: }` 替代 if-else 链。
- `range` 返回 `(index, value)`；map 的迭代顺序**故意随机化**，别依赖顺序。
- `defer` 在函数返回前按 LIFO 执行，典型用途是关文件、解锁、recover。注意 `defer` 的参数在注册时就求值。
- 循环变量在 Go 1.22 之前是复用的，闭包里捕获要传参（本例 `go func(id int){}(i)` 就是这个习惯），1.22 起已修复。

## 和其他语言的差异

| 语言 | 循环 | 分支 | 资源清理 |
|---|---|---|---|
| Go | 只有 `for` | `switch` 默认不穿透 | `defer` |
| Node.js | `for` / `while` / `for-of` | `switch` 需 `break` | `try/finally` |
| Python | `for` / `while` | `match/case`（3.10+） | `with` 上下文管理器 |
| Java | `for` / `while` / do-while | `switch` 需 `break`（新箭头语法不需要） | `try-with-resources` |
| Rust | `loop` / `while` / `for` | `match`，穷尽检查 | `Drop` trait，作用域结束自动析构 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
