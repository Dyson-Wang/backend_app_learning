# 10 · 第一个 HTTP 服务 + 数据层

## 这个例子讲什么

把前 9 课拼起来：路由、JSON、统一错误码、请求上下文中间件、数据层抽象、显式超时、优雅关闭。这是后面所有 `advanced/` 场景的起点。

## 怎么跑

```
cd backend-roadmap/01-go
PORT=8080 go run ./basics/10-first-http-service-with-db

# 另开一个终端
curl -i -H 'X-Request-ID: brm-demo-1' http://localhost:8080/healthz
curl -i -X POST http://localhost:8080/users -d '{"name":"alice"}'
curl -i http://localhost:8080/users
curl -i http://localhost:8080/users/1
curl -i http://localhost:8080/users/999      # 404 NOT_FOUND
curl -i -X POST http://localhost:8080/users -d '{"name":""}'  # 400 VALIDATION
curl -i -X DELETE http://localhost:8080/users                  # 405
# Ctrl-C 观察优雅关闭日志
```

## 关键点

- **`Store` 接口是这一课的核心**：handler 只依赖抽象，换 PostgreSQL 时不动一行 handler 代码。这就是 Java 分层架构在 Go 里的极简版。
- `http.Server` 的四个超时全部显式设置；尤其 `ReadHeaderTimeout`，不设会暴露 Slowloris 类攻击面。
- 中间件就是 `func(http.Handler) http.Handler`，链式包装；`X-Request-ID` 从网关透传，没有就自己生成，并回写到响应头。
- `/healthz` 只表示进程活着，`/readyz` 才检查依赖——别把 DB 检查放进 liveness，否则 DB 一抖会把所有实例重启一遍。
- `signal.NotifyContext` + `srv.Shutdown` 实现优雅关闭；直接 `os.Exit` 会丢掉正在处理的请求。
- **接真实数据库**：实现 `PostgresStore` 满足 `Store` 接口，用 `database/sql` + `pgx`（或你 `linnet-server` 里用的 `sqlx`），把 `BRM_PG_DSN` 从环境变量读入，先执行 `00-foundations/sql-and-postgres/schema.sql`。`advanced/C01` 会专门讲连接池配置与耗尽复现。

这个例子刻意保持零第三方依赖，所以数据层用 `sync.RWMutex` 保护的内存 map。它已经覆盖了生产服务的骨架，但**不是**生产可用的数据层：进程一重启数据就没了，且多实例之间不共享状态。

## 和其他语言的差异

| 语言 | 最小 HTTP 服务 | 中间件形态 | 优雅关闭 |
|---|---|---|---|
| Go | `net/http` 标准库 | `func(Handler) Handler` | `signal.NotifyContext` + `Shutdown` |
| Node.js | `node:http` | `req/res` + 包装函数 | `server.close()` + signal |
| Python | `http.server` / FastAPI | ASGI middleware | uvicorn 处理信号 |
| Java | `HttpServer` / Spring Boot | Filter / Interceptor | `ContextClosedEvent` + Tomcat 停机 |
| Rust | 手写 TcpListener / axum | tower `Layer` / `middleware::from_fn` | `with_graceful_shutdown` |

对照你已有的 `linnet-server`：`tower-http` 的 `request-id` / `trace` / `cors` Layer 就对应这里的 `withRequestID`，`axum` 的 `State` 对应 `API.store`。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
