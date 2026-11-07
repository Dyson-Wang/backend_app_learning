# TypeScript / Node.js · basics 基础入门

统一 **24 课**，编号在 5 门语言里一致，方便打开相同编号目录直接对照。01-10 打基础，11-24 覆盖语法细节与语言特性。

进度：**24/24**（由 `scripts/refresh_indexes.py` 自动扫描生成）

| 目录 | 主题 | 状态 |
|---|---|---|
| `01-hello-and-toolchain/` | hello 与工具链 / 项目结构 | 已实现 |
| `02-variables-and-types/` | 变量、类型、常量、零值 | 已实现 |
| `03-control-flow/` | 控制流：分支、循环、迭代 | 已实现 |
| `04-functions-and-errors/` | 函数与错误处理 | 已实现 |
| `05-collections-and-iterators/` | 集合与迭代器 | 已实现 |
| `06-structs-classes-and-modules/` | 结构体 / 类、接口、模块 | 已实现 |
| `07-package-management/` | 包管理与依赖 | 已实现 |
| `08-file-io-and-json/` | 文件 I/O 与 JSON | 已实现 |
| `09-concurrency-primitives/` | 并发原语入门 | 已实现 |
| `10-first-http-service-with-db/` | 第一个 HTTP 服务 + 数据层 | 已实现 |
| `11-operators-and-expressions/` | 运算符与表达式 | 已实现 |
| `12-strings-and-encoding/` | 字符串与字符编码 | 已实现 |
| `13-enums-and-pattern-matching/` | 枚举、模式匹配与代数数据类型 | 已实现 |
| `14-generics-and-type-parameters/` | 泛型与类型参数 | 已实现 |
| `15-iterators-and-generators/` | 迭代器与生成器 | 已实现 |
| `16-closures-and-higher-order-functions/` | 闭包与高阶函数 | 已实现 |
| `17-modules-and-visibility/` | 模块进阶：可见性与命名空间 | 已实现 |
| `18-time-regex-and-stdlib/` | 时间、正则与常用标准库 | 已实现 |
| `19-concurrent-communication/` | 并发通信：channel / 队列 / 同步原语 | 已实现 |
| `20-async-and-cancellation/` | 异步编程与取消 / 上下文 | 已实现 |
| `21-memory-model/` | 内存模型：GC / 所有权 / 生命周期 | 已实现 |
| `22-reflection-annotations-and-macros/` | 反射、注解、宏与元编程 | 已实现 |
| `23-testing-benchmarks-and-docs/` | 测试、基准与文档 | 已实现 |
| `24-capstone-http-service-revisited/` | 综合实战：语言特性用回 HTTP 服务 | 已实现 |

## 每课固定结构

- 一个可运行的最小工程（`11`~`24` 全部做成跑完即结束的脚本，不阻塞）
- `README.md` 四段：这个例子讲什么 / 怎么跑 / 关键点 / 和其他语言的差异

## 怎么读

先跑通，再改：把示例里的类型、边界、错误分支故意改错一次，看编译器或运行时怎么反馈。**编译器/运行时的报错信息本身就是这一课的一半内容。**

每课最后都留「和其他语言的差异」小节，读完请把结论补进 `../../99-cross-language/scenario-matrix.md`。
