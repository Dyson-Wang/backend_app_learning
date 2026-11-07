# 10 · Rust · 第一个 HTTP 服务 + 数据层

## 这个例子讲什么

零依赖手写 TCP + HTTP/1.1，实现和其余四门同编号一致的 CRUD。目的不是造框架，是看清 axum / tower 替你做了哪些事。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/10-first-http-service-with-db
cargo run

# 另开终端
curl -i -H 'X-Request-ID: brm-demo-1' http://localhost:8084/healthz
curl -i -X POST http://localhost:8084/users -d '{"name":"alice"}'
curl -i http://localhost:8084/users
curl -i http://localhost:8084/users/999
```

## 关键点

- `Arc<Mutex<HashMap<...>>>`：`Arc` 让多个线程共享所有权，`Mutex` 保证互斥。**编译器强制你显式表达共享**，Go 里随手一个包级 map 在这里会被借用检查器拦下。
- 每连接 `std::thread::spawn` 一个 OS 线程，是最朴素的模型；生产用 `tokio::spawn` 异步任务，成本才和 goroutine 一个量级。
- `Content-Length` 必须用 `body.as_bytes().len()`（字节数），中文 UTF-8 下 `chars().count()` 会算错，客户端会一直等或截断。
- `Connection: close` 简化连接管理；keep-alive、chunked、超时、背压全部没实现——**这就是「为什么需要框架」的答案**。
- `extract_name` / `json_escape` 手写 JSON 只为教学，生产用 `serde` + `serde_json`（编译期生成，零反射开销）。
- **接数据库**：`sqlx` 或 `diesel` + 连接池；你已有的 `linnet-server` 用的正是 `sqlx` + PostgreSQL。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Rust（本例） | Rust（生产） |
|---|---|---|---|
| HTTP | `net/http` 标准库 | 手写 `TcpListener` | axum / actix-web |
| 中间件 | `func(Handler) Handler` | 无 | tower `Layer` |
| 共享状态 | `sync.Mutex` | `Arc<Mutex<T>>` | `Arc` + `State` |
| JSON | `encoding/json` | 手写 | serde_json |
| 优雅关闭 | `Shutdown(ctx)` | 无 | `with_graceful_shutdown` |

**下一步直接看 `~/workwork2/linnet-server`**：axum 0.8 + tokio + sqlx + tower-http（cors / trace / request-id / compression / limit）+ JWT/bcrypt + tracing，几乎是 `advanced/` 里 `A01`、`D04`、`E01`、`E02`、`E04` 的现成生产答案。这门语言你不用从零学，拿它当对照案例效率最高。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
