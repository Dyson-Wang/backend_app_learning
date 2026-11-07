# 24 · Rust · 综合实战：语言特性用回 HTTP 服务

## 这个例子讲什么

把前面 23 课的东西拼成一个能跑的最小 HTTP 服务：`TcpListener` 绑随机端口、子线程 `accept`、用 `enum Route` + `match` 做路由、`String`/`format!` 拼响应、`Result`/`Option` 处理 I/O 与解析、主线程发三个真实请求验证 200 和 404，最后 `join` 退出。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/24-capstone-http-service-revisited
cargo run --quiet
```

## 关键点

- 绑 `127.0.0.1:0` 让内核分配端口，`local_addr()` 拿回来再连接，避免端口冲突。
- 手写 HTTP/1.1 响应必须带 `Content-Length` 和 `Connection: close`，否则客户端会一直等。
- 用 `enum` 表达路由，`match` 保证新增路由时不会漏处理——这是第 13 课的实战用法。
- 所有 I/O 都返回 `Result`，本例用 `unwrap_or` / `if let Ok` 处理，生产代码应映射成 500 并记录日志。
- 生产不会手写 TCP + HTTP：用 `axum` / `actix-web` / `hyper`，本例只为展示协议层。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| HTTP 服务 | `net/http` | `node:http` | `http.server` | `com.sun.net.httpserver` | 标准库只有 TCP，需手写或 `hyper` |
| 路由 | `ServeMux` | 手写 if/switch | 手写 | `HttpContext` 匹配 | 手写 `match` |
| 并发模型 | 每请求 goroutine | 事件循环 | 每请求线程 | 线程池 | 每连接线程 / tokio |
| 生产框架 | gin / echo | express / fastify | FastAPI / Flask | Spring Boot | axum / actix-web |

Go 和 Node 的标准库自带 HTTP 服务器，Python/Java 也有，Rust 标准库刻意只给 TCP——HTTP 属于生态层，交给 `hyper`/`axum`。这个差异正好说明「标准库边界」的设计哲学。
