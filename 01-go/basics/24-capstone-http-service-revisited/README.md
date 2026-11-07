# 24 · 综合实战：语言特性用回 HTTP 服务

## 这个例子讲什么

把前 23 课的语言特性拼成一个**内存版 HTTP 服务**：泛型 `writeJSON[T]`/`Store[T]`、`iota` 错误码枚举、闭包日志中间件、`RWMutex`+`atomic` 并发保护、`context` 超时取消、`time` 格式化、`encoding/json` 与输入校验。用 `net/http/httptest` 起服务，自己请求 6 次，然后自行关闭，3 秒内退出。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/24-capstone-http-service-revisited
```

## 关键点

- **自请求自关闭**：`httptest.NewServer(handler)` 在随机回环端口监听，`srv.Client()` 配 `srv.URL` 自己请求自己；`defer` 里 `srv.Close()` 关闭监听并等待在途请求结束，进程退出码 0。地址含随机端口，**不打印具体 URL** 以保持输出可复现。
- **路由**：`http.ServeMux` + `/health`、`/tasks`、`/tasks/`。注意 `go.mod` 是 **go 1.21**，不支持 `"GET /tasks"` 这类带方法的模式（1.22 才有），所以方法判断放在 handler 内的 `switch r.Method`。
- **泛型落地**：`writeJSON[T any]` 统一编码，`Store[T]` 泛型存储配 `setID func(*T,int)` 闭包把自增 id 写回 `Task`；`Store` 用 `sync.RWMutex` 保护，`List()` 按 id 顺序输出保证稳定。
- **枚举 + 错误处理**：`ErrCode`（`iota`）自带 `String()` 与 `HTTPStatus()`，把领域错误码映射成 400/404/409；校验失败（`title` 为空）走 `BAD_REQUEST`。
- **闭包中间件**：`withLogging` 包装 handler，用 `atomic.Int64` 生成请求号、`statusRecorder` 捕获真实状态码，把日志送进 channel；主协程每发一个请求就收回一条日志，顺序稳定。
- **取消/时间**：`context.WithTimeout(..., 2*time.Second)` 给所有请求兜底；时间只演示固定时刻的 `Format`，不打印当前绝对时间。
- 全部只用标准库、零第三方依赖、无外部端口暴露、无常驻。

## 和其他语言的差异

| 语言 | 本课内存服务做法 | 关闭方式 |
|---|---|---|
| Go | `net/http` + `httptest.NewServer` 随机端口 + `srv.Client()` 自请求 | `srv.Close()`（defer） |
| Node.js | `node:http` 绑 `0` 端口 + `fetch` 自请求 | `server.close()` |
| Python | `ThreadingHTTPServer` 起在子线程，主线程 `urllib` 自请求 | `shutdown()` + `server_close()` |
| Java | `com.sun.net.httpserver.HttpServer` 随机端口 + `HttpClient` 自请求 | `stop(0)` |
| Rust | `std::net::TcpListener` 绑 `127.0.0.1:0`，子线程手写 HTTP，主线程自请求 | `drop(listener)` + `join` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
