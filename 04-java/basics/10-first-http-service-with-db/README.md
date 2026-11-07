# 10 · Java · 第一个 HTTP 服务 + 数据层

## 这个例子讲什么

和其余四门语言同编号功能一致的最小 CRUD，用 JDK 自带 `HttpServer`（零第三方依赖）：线程池、并发 Map、统一错误码、关闭钩子。

## 怎么跑

```
cd backend-roadmap/04-java/basics/10-first-http-service-with-db
javac Main.java && PORT=8083 java Main

# 另开终端
curl -i -H 'X-Request-ID: brm-demo-1' http://localhost:8083/healthz
curl -i -X POST http://localhost:8083/users -d '{"name":"alice"}'
curl -i http://localhost:8083/users
curl -i http://localhost:8083/users/999
# Ctrl-C 观察关闭钩子
```

## 关键点

- **默认 executor 是单线程的**，这里显式换成固定线程池。生产要用带界队列的 `ThreadPoolExecutor`，JDK 21+ 可以直接用 `Executors.newVirtualThreadPerTaskExecutor()`——这就是 Java 对高并发 I/O 的 modern 答案。
- `ConcurrentHashMap` + `AtomicLong` 对应 Go 的 `sync.RWMutex` + 自增、Python 的 `threading.Lock`，是共享状态的另一种解法。
- `addShutdownHook` + `server.stop(2)` 对应 Go 的 `signal.NotifyContext` + `Shutdown(ctx)`；`stop` 的 2 秒就是兜底超时。
- 手写 JSON 只为教学，**生产请用 Jackson**（Spring Boot 默认）+ DTO，别自己拼字符串。
- **接数据库**：Spring Boot + `JdbcTemplate` / JPA + HikariCP 连接池；声明式事务 `@Transactional` 是 `advanced/B01` 的重点，也是 Java 相对其他语言最有辨识度的一块。
- 本机无真实 JDK，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Java | Rust |
|---|---|---|---|
| 数据层抽象 | `Store` 接口 | `interface` + Spring 依赖注入 | `trait` + 泛型 |
| 并发单元 | goroutine | 平台线程 / 虚拟线程 | `tokio::spawn` |
| 共享状态 | `sync.Mutex` | `ConcurrentHashMap` / `Lock` | `Arc<Mutex<T>>` |
| 事务 | 手动 `Tx` | `@Transactional` 声明式 | 手动 `Transaction` |
| 错误 | `error` 返回值 | 异常 | `Result<T,E>` |

Java 的差异化不在于这些接口怎么写，而在于**分层 + IoC + 声明式事务**让大团队能在同一个约束下协作——`advanced/` 会重点展开。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
