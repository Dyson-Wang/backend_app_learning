# Backend Roadmap · 后端多语言示例代码库

一条「场景 × 语言」的后端学习路线：同一套基础概念与同一批生产场景，在 5 门语言里各写一遍，靠横向对照理解差异，而不是孤立地学语法。

规划原稿见 `../.workbuddy-ai/reports/backend-roadmap-plan.md`。

## 语言与顺序

| 顺序 | 语言 | 目录 | 这一门补的认知 | 运行时模型 |
|---|---|---|---|---|
| 01 | Go | `01-go/` | 并发模型 ① | CSP · goroutine / channel · GC · 单二进制 |
| 02 | TypeScript / Node.js | `02-typescript-node/` | 并发模型 ② | 事件循环 · 非阻塞 I/O · GC |
| 03 | Python | `03-python/` | 类型系统 | GIL + asyncio · 动态类型 |
| 04 | Java | `04-java/` | 工程化 | JVM · 虚拟线程 · Spring 分层与事务 |
| 05 | Rust | `05-rust/` | 内存模型 | 所有权 · Tokio · 无 GC · 单二进制 |

选修（体验式，不铺开全量场景）：C# / .NET、Kotlin、Elixir / BEAM，见 `06-elective/`。

## 目录

```
backend-roadmap/
├── README.md                  本文件：路线图、进度、统一约定
├── scripts/verify.sh          一键验证所有语言的示例
├── 00-foundations/            语言无关的公共底座（6 块）
├── 01-go/                     basics/ + advanced/
├── 02-typescript-node/
├── 03-python/
├── 04-java/
├── 05-rust/
├── 06-elective/               选修语言
└── 99-cross-language/         同一场景的 5 语言对照矩阵
```

每门语言固定两层：

- `basics/` —— 基础入门，统一 24 个编号主题（见下表），覆盖语法与语言特性
- `advanced/` —— 生产场景，统一 `A01` ~ `F03` 编号（见 `99-cross-language/scenario-matrix.md`）

## `basics/` 统一 24 课

01-10 打基础，11-24 补齐语法细节与语言特性。编号在 5 门语言里一致，方便打开 `01-go/basics/04-functions-and-errors/` 和 `05-rust/basics/04-functions-and-errors/` 直接对照。

| 编号 | 主题 | 编号 | 主题 |
|---|---|---|---|
| 01 | hello 与工具链 / 项目结构 | 13 | 枚举、模式匹配与代数数据类型 |
| 02 | 变量、类型、常量、零值 | 14 | 泛型与类型参数 |
| 03 | 控制流：分支、循环、迭代 | 15 | 迭代器与生成器 |
| 04 | 函数与错误处理 | 16 | 闭包与高阶函数 |
| 05 | 集合与迭代器 | 17 | 模块进阶：可见性与命名空间 |
| 06 | 结构体 / 类、接口、模块 | 18 | 时间、正则与常用标准库 |
| 07 | 包管理与依赖 | 19 | 并发通信：channel / 队列 / 同步原语 |
| 08 | 文件 I/O 与 JSON | 20 | 异步编程与取消 / 上下文 |
| 09 | 并发原语入门 | 21 | 内存模型：GC / 所有权 / 生命周期 |
| 10 | 第一个 HTTP 服务 + 数据层 | 22 | 反射、注解、宏与元编程 |
| 11 | 运算符与表达式 | 23 | 测试、基准与文档 |
| 12 | 字符串与字符编码 | 24 | 综合实战：语言特性用回 HTTP 服务 |

`11`~`24` 全部做成跑完即结束的脚本（只有 `10` 会起常驻服务），因此可以批量运行验证，见 `scripts/verify.sh`。

## 每个示例的统一结构

- 一个**可运行的最小工程**，不依赖 IDE、不依赖付费服务
- 一份 `README.md`，固定四段：`这个例子讲什么`、`怎么跑`、`关键点`、`和其他语言的差异`
- `advanced/` 额外要求：`问题复现` → `错误写法` → `正确写法` → `怎么验证`

零外部依赖原则：`basics/` 尽量只用标准库；确实需要 PostgreSQL / Redis / MQ 的场景，在 README 给出可选接法与可复现步骤，但不作为运行前提。你的 homelab 已有 PostgreSQL、MariaDB、Redis、Nginx 统一反代，进阶示例可以按 README 直接接上去跑。

## 进度看板

状态由 `python3 scripts/refresh_indexes.py` 自动扫描生成，新增示例后重跑即可。

<!-- PROGRESS:START -->
| 模块 | 状态 |
|---|---|
| 目录骨架 / 总索引 | 完成 |
| `00-foundations/` 六块文档与配置 | 完成 |
| Go `basics/` | 24/24 |
| Go `advanced/` | 28/28 |
| TypeScript / Node.js `basics/` | 24/24 |
| TypeScript / Node.js `advanced/` | 28/28 |
| Python `basics/` | 24/24 |
| Python `advanced/` | 28/28 |
| Java `basics/` | 24/24 |
| Java `advanced/` | 28/28 |
| Rust `basics/` | 24/24 |
| Rust `advanced/` | 28/28 |
| `06-elective/` | 说明已完成，代码未开始 |
| `99-cross-language/` | 矩阵自动生成，28 个场景 + 24 课结论已填满 |
<!-- PROGRESS:END -->

## 本机环境（2026-11-08）

| 工具 | 状态 | 影响 |
|---|---|---|
| Node 22 | 可用 | `basics` 22 个 + `advanced` 3 个全部实际运行通过 |
| Python 3.13 | 可用 | `basics` 22 个 + `advanced` 3 个全部实际运行通过 |
| Go | 缺失 | 32 个 `.go` 文件**未经编译验证**，只做了静态自检 |
| Rust / cargo | 缺失 | 27 个 `.rs` 文件**未经编译验证** |
| javac | macOS 占位 stub，无真实 JDK | 27 个 `.java` 文件**未经编译验证** |
| docker / psql / redis-cli | 缺失 | 依赖容器与数据库的场景未在目标主机验证 |

当前网络到 `go.dev` 不可达（探测超时），所以没能在隔离目录临时装一份 Go SDK 来补验证。装好工具链后跑 `bash scripts/verify.sh` 会一次性补齐：Go 编译 + vet + 全部场景自测、Java `javac`、Rust `cargo check` 与 `cargo run`。

**未经编译验证的部分**：Go 32 个、Java 27 个、Rust 27 个源文件。这些代码已做静态自检（结构、括号平衡、入口函数、README 段落，140 个源文件 0 错误），但静态检查替代不了编译器，装好工具链后请以 `verify.sh` 的结果为准。

**已实测的部分**：`scripts/verify.sh` 当前 106 项通过、0 失败、3 项跳过；Node 与 Python 的每个示例都真跑过，且都在 2 秒内自行退出。

## 怎么开始

1. 先读 `00-foundations/request-lifecycle.md`，建立「一个请求从进入到返回」的链路心智
2. 从 `01-go/basics/01-hello-and-toolchain/` 开始，按顺序做完 24 课；11-24 是语法细节与语言特性，别跳
3. 做 `01-go/advanced/` 的 28 个场景，每个都先自己想「错在哪」再看正确写法
4. 换下一门语言时，打开相同编号的目录对照着写；`basics/13`、`14`、`20`、`21` 的横向差异最大
5. 对照 `99-cross-language/scenario-matrix.md` 里已填的差异结论，检验自己的理解是否一致
