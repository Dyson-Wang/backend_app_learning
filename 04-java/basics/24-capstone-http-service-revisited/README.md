# 24 · Java · 综合实战：语言特性用回 HTTP 服务

## 这个例子讲什么

把前 23 课的特性（泛型 `ApiResult<T>`、枚举 `Code`、`record`、闭包、`instanceof` 模式匹配、错误处理与错误码、并发 Map/原子类、`java.time` 时间格式化）重新用回一个**内存版** HTTP 服务：路由、JSON、校验、错误码、日志；随机端口 + 自请求 + `stop(0)`。

## 怎么跑

```
cd backend-roadmap/04-java/basics/24-capstone-http-service-revisited
javac Main.java && java Main

# 进程自己请求自己并在 3 秒内退出，退出码 0
# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **随机端口**：`HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0)` 绑端口 0 由系统分配，`getAddress().getPort()` 拿到实际端口；不打印端口值以保证输出可复现。
- **自请求自关闭**：`java.net.http.HttpClient`（JDK 内置）带 `connectTimeout` / request `timeout`，依次请求 `/health`、`POST /users`（成功 + 空名校验失败）、`GET /users`、`GET /users/999`；随后 `server.stop(0)`、`pool.shutdownNow()`，并 `System.exit(0)` 兜底，确保 3 秒内退出。
- **特性回用**：`record ApiResult<T>` 泛型统一响应；`enum Code` 映射状态码；`server.createContext` 的 lambda **闭包捕获** `service`；`instanceof User u` 做解构；`ConcurrentHashMap` + `AtomicLong` 处理并发；`DateTimeFormatter` 用**固定时刻**格式化日志时间，不打印当前绝对时间。
- **日志顺序稳定**：handler 先记日志、再写响应体，客户端收完响应才继续下一个请求，因此输出顺序确定。
- 手写 JSON 仅教学，生产用 Jackson + DTO；真实服务见 `10-` 与 `advanced/`。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| HTTP 服务 | `com.sun.net.httpserver`/`java.net.http` | `net/http` + `httptest` | `node:http` | `ThreadingHTTPServer` | `TcpListener` 手写 |
| 随机端口 | 绑 `:0` + `getPort()` | `httptest.NewServer` | `server.listen(0)` | 绑 `("", 0)` | 绑 `127.0.0.1:0` |
| 关闭 | `stop(0)` / `close()` | `Shutdown(ctx)` | `server.close()` | `shutdown()` | drop listener |
| 综合特性 | 泛型/枚举/record/闭包 | iota/interface/闭包 | 判别联合/闭包 | dataclass/match/闭包 | enum/trait/闭包 |
| 退出保证 | `stop(0)` + `System.exit(0)` | 进程自然结束 | `close()` 回调 | `shutdown()` | 主线程返回 |
