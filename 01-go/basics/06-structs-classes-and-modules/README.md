# 06 · 结构体、接口、模块

## 这个例子讲什么

Go 没有类、没有继承、没有 `implements`。用 struct 存数据、用方法绑行为、用组合复用、用隐式接口解耦。

## 怎么跑

```
cd backend-roadmap/01-go && go run ./basics/06-structs-classes-and-modules
```

## 关键点

- **值接收者 vs 指针接收者**：要改接收者、或 struct 很大、或要保持一致，用指针接收者。混用会让方法集不一致，容易踩接口实现的坑。
- **接口隐式实现**：只要方法集匹配就算实现，包之间不需要互相依赖。因此 Go 倾向于在使用方定义**小接口**（`Storer` 只有两个方法）。
- **嵌入是组合不是继承**：`Service` 获得 `Logger.Log`，但没有子类型关系，不能把 `*Service` 当成 `Logger` 传。
- 类型断言 / type switch 用 `x.(T)` 和 `switch v := x.(type)`；`interface{}`（`any`）应少用。
- 包可见性只有一条规则：首字母大写导出；`internal/` 目录额外限制只能被父目录内的包 import。

## 和其他语言的差异

| 语言 | 抽象机制 | 接口实现 |
|---|---|---|
| Go | struct + 方法 + 嵌入 | 隐式，使用方定义小接口 |
| Node.js | class（原型语法糖）/ TS interface | 结构化类型，编译期擦除 |
| Python | class + 多继承 / Protocol | 鸭子类型 |
| Java | class + 单继承 + interface | 显式 `implements` |
| Rust | struct + impl + trait | 显式 `impl Trait for Type` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
