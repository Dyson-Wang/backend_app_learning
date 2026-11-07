# A04 · 优雅关闭与健康检查

## 问题复现

发布或扩缩容时进程收到 `SIGTERM`，如果直接退出，正在处理的请求会被砍断，客户端拿到 502 / 半截 JSON，写操作可能停在「已扣款未落单」的中间态。

本例起真实 TCP 服务，发一个 200ms 的慢请求，等它进入 handler 后调用 `Shutdown`，验证：在途请求完整返回、handler 执行完毕、关闭后新连接被拒绝。

## 错误写法

**错法 1 · 收到信号直接退出**：

```go
c := make(chan os.Signal, 1)
signal.Notify(c, syscall.SIGTERM)
<-c
os.Exit(0)   // 在途请求全丢
```

**错法 2 · 用 `srv.Close()` 代替 `Shutdown()`**：`Close` 立即关闭所有连接，包括活跃的；`Shutdown` 才是不再收新连接 + 等在途请求完成。

**错法 3 · 无限等待**：不给 `Shutdown` 传带超时的 ctx，一个卡死的请求能让 Pod 永远退不掉，最后被 `SIGKILL`。必须设兜底上限。

**错法 4 · 顺序反了**：先停进程再从网关摘除，摘除前的新流量会打到已关闭的端口。正确顺序是 **先从上游摘除 → 停止接新流量 → 等在途完成 → 关闭连接池/消费者 → 退出**。

**错法 5 · liveness 里查数据库**：DB 一抖，编排系统把所有健康实例全部重启，故障被放大。`/healthz` 只查进程，`/readyz` 才查依赖。

## 正确写法

`signal.NotifyContext` 拿信号 → `srv.Shutdown(ctx)` 带超时等待 → 超时后强制收尾。本例用 `net.Listen("127.0.0.1:0")` 拿随机端口，避免和本机其他服务冲突；`serveErr` 缓冲通道用于确认 `Serve` 正常返回 `http.ErrServerClosed`。

生产还要补三件事：

- `ReadHeaderTimeout` / `ReadTimeout` / `WriteTimeout` 显式设置（见 `basics/10`）
- 消费者（MQ）、后台任务、连接池也要纳入同一套关闭流程
- K8s 里配 `preStop` sleep 给上游摘除留出时间，`terminationGracePeriodSeconds` 要大于 Shutdown 超时

`basics/10-first-http-service-with-db` 里有可直接运行的服务版实现。

## 怎么验证

```
cd backend-roadmap/01-go
go run ./advanced/A04-graceful-shutdown
# 期望：Shutdown 约 150ms 后返回（等待在途请求），status=200 handler_completed=1，selftest PASS
```

相关底座概念见 `../../../00-foundations/`。做完把差异结论补到 `../../../99-cross-language/scenario-matrix.md`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 触发 | `signal.NotifyContext` | `process.on('SIGTERM')` | `signal.signal` | shutdown hook | `tokio::signal` |
| 等待在途 | `srv.Shutdown(ctx)` | `server.close()` | `server.shutdown()` | `server.stop(grace)` | `with_graceful_shutdown` |
| 强制兜底 | `ctx` 超时后 `Close()` | `closeAllConnections()` | `server_close()` | `stop(0)` | 超时后 abort |
