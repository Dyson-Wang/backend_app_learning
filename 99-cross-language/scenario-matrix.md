# 生产场景 × 5 语言对照矩阵

状态列由 `scripts/refresh_indexes.py` 自动扫描生成；「差异结论」列由你每做完一个场景后亲手填，刷新时会保留。

### A · 流量与稳定性

| 编号 | 场景 | Go | Node.js | Python | Java | Rust | 差异结论 |
|---|---|---|---|---|---|---|---|
| A01 | 限流：令牌桶 / 滑动窗口，单机与分布式 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 用 x/time/rate 或手写桶 + 可注入时钟；Node 用回调时钟；Python 用可调用 clock；Java 用内部 long now；Rust 用 &mut self + 原子时钟。共同点：burst 允许突发，rate 卡长期速率 |
| A02 | 超时、重试、指数退避与抖动 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 用 context 传播超时；Node 用 AbortSignal；Python 用 asyncio.timeout；Java 用 orTimeout；Rust 用 Result + 枚举区分成功/永久失败/耗尽。共同点：只重试可重试错误，指数退避 + 抖动 |
| A03 | 熔断与降级 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 用 Mutex + 函数时钟；Node 靠「临界区不 await」+ probeInFlight；Python 用 threading.Lock + 可调用时钟；Java 用 synchronized + 内部 long 时钟；Rust 用 Mutex&lt;Inner&gt; + Arc&lt;AtomicU64&gt;。共同点：open 态绝不发下游请求 |
| A04 | 优雅关闭、健康检查、上游摘除 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 四门语言的标准库 HTTP 服务都自带「停止接收 + 等待在途」；Rust 标准库只有 TcpListener，得自己写非阻塞 accept 循环 + 停止标志。 |
| A05 | 背压与并发上限 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 用 chan 容量 + select default；Node 单线程只能手动判长度或 p-limit；Python 用 Queue(maxsize) + put_nowait 抛 Full；Java 用 ArrayBlockingQueue.offer 返回布尔；Rust 用 sync_channel + try_send 返回 Result。共同点：拒绝在入口，不排队 |

### B · 一致性与并发

| 编号 | 场景 | Go | Node.js | Python | Java | Rust | 差异结论 |
|---|---|---|---|---|---|---|---|
| B01 | 事务与隔离级别、乐观锁 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 只有 Java 有声明式 @Version；其余都要手写 CAS 重试。Node 单线程也要防，因为 await 之间会切走。 |
| B02 | 幂等键与重复请求 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 占位原语差异最大：singleflight / Promise / Event / putIfAbsent / Condvar；落库后都收敛到唯一索引 + ON CONFLICT。 |
| B03 | 分布式锁 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 加锁写法五花八门，解锁逻辑五门完全一致：必须「比较 token + 删除」原子执行。 |
| B04 | 竞态、死锁与调试 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 数据竞争：Go 靠 -race 检测，Python 有 GIL 仍会丢更新，Java 需 synchronized，Rust 编译期禁止，Node 表现为异步交错。死锁：统一加锁顺序或 tryLock/超时检测 |
| B05 | 上下文取消与跨服务传播 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 把 ctx 当首参强制传播；Rust 忘传会编译报类型错；Node/Python/Java 全靠约定，忘了就静默失去取消。 |

### C · 数据与缓存

| 编号 | 场景 | Go | Node.js | Python | Java | Rust | 差异结论 |
|---|---|---|---|---|---|---|---|
| C01 | 连接池与资源耗尽 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 等待超时的表达差异最大：Go 用 ctx，Java 用 Semaphore.tryAcquire，Rust 用 Condvar::wait_timeout，Python/Node 靠库配置。 |
| C02 | N+1 查询与批量写 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 写法都是「循环里查依赖」，差别在 ORM 是否提供预加载；懒加载会让 N+1 藏得很深。 |
| C03 | 游标分页替代 OFFSET | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Python 原生元组比较一行搞定，Java 要 Comparator 链，Go/Rust 手写条件。语义完全一致。 |
| C04 | 缓存穿透 / 击穿 / 雪崩 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go/Python/Java/Rust 都是 per-key 锁 + 双重检查，Java 用 computeIfAbsent 原子建锁；Node 用 in-flight Promise Map 最简洁。共同点：回源必须按 key 单飞，50 并发从 50 次降到 1 次 |
| C05 | 缓存与 DB 双写一致性 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 逻辑一样，差异只在延迟双删的定时器：Timer / setTimeout / ScheduledExecutor / 线程。 |

### D · 消息与异步

| 编号 | 场景 | Go | Node.js | Python | Java | Rust | 差异结论 |
|---|---|---|---|---|---|---|---|
| D01 | 至少一次投递与重复消费 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 去重载体不同，核心约束相同：去重表必须和业务写在同一事务，否则还是会重复。 |
| D02 | 死信队列与重试策略 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 用 enum 让「可重试/永久」成为类型的一部分，其余靠异常类型或哨兵错误，漏分类编译器不管。 |
| D03 | 顺序性保证 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Node/Python 单线程天然串行，是最容易做对的；Go/Java/Rust 一不小心就多消费者抢同一分区。 |
| D04 | SSE / WebSocket / gRPC 流式 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 都有 flush 原语，差异在背压：只有 Node 的 write() 会返回 false 让你等 drain，其余是阻塞写。 |

### E · 安全与工程化

| 编号 | 场景 | Go | Node.js | Python | Java | Rust | 差异结论 |
|---|---|---|---|---|---|---|---|
| E01 | JWT 刷新与吊销 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 标准库不含 HMAC/SHA-256，必须引 crate；其余四门内建，代码几乎一样。 |
| E02 | 密码哈希与 RBAC | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | KDF 内建程度：Python/Node/Java 一行，Go 要手写 PBKDF2，Rust 全靠 argon2 crate。 |
| E03 | 输入校验与统一错误码 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 统一结构 {code, field, message}；「未提供」表达：Go 用 *int、Node 用 undefined、Python 用 None、Java 用 Integer null、Rust 用 Option<i32>。共同点：一次收集全部错误，code 与 message 分离 |
| E04 | 结构化日志 / trace id / 指标 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | trace 透传：Go 显式 context、Node AsyncLocalStorage、Python contextvars、Java MDC、Rust tracing::Span；日志统一一行 JSON，指标统一 RED 三件套。Java MDC 在线程池复用时要手动清理 |
| E05 | 配置、密钥与 DB migration | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 的 sqlx 把 migration 编译进二进制最特别；其余靠外部工具 + 版本表，本质一致。 |
| E06 | 单测、集成测试与容器化 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 测试框架五门都内建；差异在 mock 生态：Java 的 Mockito、Python 的 unittest.mock 最成熟，Go/Rust 靠接口 + 手写 fake。 |

### F · 故障复盘

| 编号 | 场景 | Go | Node.js | Python | Java | Rust | 差异结论 |
|---|---|---|---|---|---|---|---|
| F01 | 慢查询拖垮连接池 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 超时的表达：Go 是上下文、Rust 是组合子、Java 是注解、Node/Python 是库参数。慢查询本身跟语言无关。 |
| F02 | 重试风暴 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 退避库生态各异，结论一致：上限 + 退避 + 抖动 + 熔断，缺一个就放大故障。 |
| F03 | 消息积压 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 并发模型差异最大：Python 要多进程绕 GIL，Node 要 worker_threads，Go/Rust 开 task 最轻。 |

## 基础知识对照区

| 编号 | 主题 | Go | Node.js | Python | Java | Rust | 结论 |
|---|---|---|---|---|---|---|---|
| 01 | hello 与工具链 / 项目结构 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 出单二进制、Rust 出原生码、Java 出字节码靠 JVM、Node/Python 靠解释器——部署形态从这里就分叉了。 |
| 02 | 变量、类型、常量、零值 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 静态强类型：Go/Java/Rust；静态但运行时擦除：TS；动态：Python。零值/默认值规则五门都不同。 |
| 03 | 控制流：分支、循环、迭代 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 的 match 是表达式且强制穷尽；Go 的 switch 默认不穿透；Python 靠缩进。 |
| 04 | 函数与错误处理 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 返回 error、Rust 返回 Result、Java/Python 抛异常、Node 两者混用——这是错误处理风格的分水岭。 |
| 05 | 集合与迭代器 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 的 slice 共享底层数组最容易踩坑；Rust 的 Vec 独占所有权；Java 集合有装箱开销。 |
| 06 | 结构体 / 类、接口、模块 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 接口是隐式实现，Rust 是 trait，Java 必须显式 implements，TS 是结构化类型——抽象成本差异很大。 |
| 07 | 包管理与依赖 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | go mod / npm / pip / Maven / cargo：依赖解析策略和锁文件机制各不相同，供应链安全要从这里管。 |
| 08 | 文件 I/O 与 JSON | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 靠 struct tag、Rust 需 serde、Java 需 Jackson；标准库直出 JSON 的只有 Go 和 Node。 |
| 09 | 并发原语入门 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | goroutine 最轻（KB 级）、Java 平台线程最重（MB 级）、Python 受 GIL 限制、Rust 靠 Send+Sync 在编译期挡数据竞争。 |
| 10 | 第一个 HTTP 服务 + 数据层 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 四门语言标准库自带 HTTP 服务；Rust 标准库只有 TcpStream，要么手写协议要么引 tokio/axum。 |
| 11 | 运算符与表达式 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 整数溢出：Go 静默回绕、Java 回绕、Rust debug 下 panic、Python 任意精度。这个差异会直接变成 bug。 |
| 12 | 字符串与字符编码 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 字符串模型是重灾区：Go/Rust 是 UTF-8 字节序列（按字节切片会切坏汉字），Java/JS 是 UTF-16，Python 是码点序列。 |
| 13 | 枚举、模式匹配与代数数据类型 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 只有 Rust 的 enum 能带数据（真正的代数数据类型）；Java 用 sealed + 模式匹配勉强接近，Go/Python/TS 靠联合类型或约定。 |
| 14 | 泛型与类型参数 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Java/TS 泛型运行时擦除，Go 用字典 + 单态化，Rust 纯单态化——性能与二进制体积的取舍不同。 |
| 15 | 迭代器与生成器 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 迭代器是零成本抽象（编译期展开），Python/JS 生成器是协程实现，Go 1.23 才有 range-over-func。 |
| 16 | 闭包与高阶函数 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 闭包分 Fn/FnMut/FnOnce 三种 trait，其余语言统一成函数对象，捕获规则更宽松。 |
| 17 | 模块进阶：可见性与命名空间 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 靠首字母大小写决定导出，Rust 用 pub，Java 用修饰符，Python 靠下划线约定——只有 Go 是语法级强制。 |
| 18 | 时间、正则与常用标准库 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 时间格式化差异最大：Go 用「参考时间」2006-01-02，Java 用 DateTimeFormatter，其余用 strftime 风格。 |
| 19 | 并发通信：channel / 队列 / 同步原语 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Go 的 channel 是语言内建，Rust 的 mpsc 在标准库，Java/Python/Node 要靠并发队列库。 |
| 20 | 异步编程与取消 / 上下文 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 必须显式选 runtime（tokio），Go 用 goroutine 把异步藏起来，JS/Python 是事件循环 + async/await。 |
| 21 | 内存模型：GC / 所有权 / 生命周期 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 所有权在编译期释放，无 GC；Go/Java 是并发 GC；Python 是引用计数 + 循环 GC——停顿特征完全不同。 |
| 22 | 反射、注解、宏与元编程 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | Rust 宏是编译期代码生成，Java 是注解 + 运行时反射，Go 几乎没有元编程能力，Python/JS 靠动态特性。 |
| 23 | 测试、基准与文档 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 五门都内建测试框架；差异在基准测试（Go 内建 Benchmark、Rust 需 criterion）和文档测试（只有 Rust 内建）。 |
| 24 | 综合实战：语言特性用回 HTTP 服务 | 已实现 | 已实现 | 已实现 | 已实现 | 已实现 | 同一个 HTTP 服务，代码量差异主要来自错误处理和并发样板：Go 最短，Java 最啰嗦，Rust 编译期检查最严。 |

## 五个机制速查（其余场景对照时反复用到）

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 取消 / 超时 | `context.Context` | `AbortSignal` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 请求上下文 | `context.WithValue` | `AsyncLocalStorage` | `contextvars` | `ThreadLocal` + MDC | `tower` Extension |
| 并发单元 | goroutine | 事件循环 + Promise | Task / 进程 | 平台线程 / 虚拟线程 | `tokio::spawn` |
| 共享状态 | `sync.Mutex` / channel | 单线程无锁 | `asyncio.Lock` / GIL | `synchronized` / Lock | `Arc<Mutex<T>>` |
| 错误表达 | `error` 返回值 | `throw` / rejected Promise | `Exception` | 受检 / 非受检异常 | `Result<T, E>` |
| 内存管理 | GC | GC | 引用计数 + GC | JVM GC | 所有权，编译期 |
| 交付 | 单二进制 | Node 运行时 | CPython 运行时 | JVM 字节码 | 单二进制 |
