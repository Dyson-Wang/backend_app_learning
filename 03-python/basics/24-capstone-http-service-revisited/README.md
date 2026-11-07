# 24 · Python · 综合实战：语言特性用回 HTTP 服务

## 这个例子讲什么

把前 23 课的特性重写进一个**内存版** HTTP 服务：泛型 `Result[T]`、`enum` 错误码、`dataclass`、`match case` 路由、闭包日志、`try/except` 校验、`ThreadingHTTPServer` + `Lock` 并发、`datetime` 时间格式化、JSON、错误码与访问日志。子线程起服务、主线程 `urllib` 自请求，最后 `shutdown()` 自关闭，3 秒内退出。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/24-capstone-http-service-revisited/main.py
```

## 关键点

- **自请求自关闭**：`ThreadingHTTPServer(("127.0.0.1", 0), ...)` 由内核分配随机端口；`serve_forever` 放子线程，主线程发完 7 个请求后调用 `server.shutdown()` + `server.server_close()` + `join()`。
- `shutdown()` **必须从 `serve_forever` 之外的线程调用**，否则死锁——这是 10 课同样的坑。
- 语言特性回炉：`class Result[T]`（14）、`ErrCode(Enum)`（13）、`@dataclass`+`field`（06）、`match (method, path)` 路由（13）、闭包 `make_logger`（16）、`ValueError`+异常捕获（04）、`Lock`（09/19）、固定基准 `datetime` 格式化（18）。
- `HTTP/1.1` 必须设 `Content-Length`（按字节数，`ensure_ascii=False` 时中文占 3 字节）；非 2xx 用 `urllib.error.HTTPError` 仍可读出 JSON 错误体。
- 日志用「固定基准时间 + 序号」演示格式化与差值，**不打印当前绝对时间**；输出按请求序号稳定可复现。

## 和其他语言的差异

- **Go**：`net/http` + `httptest.NewServer` 起内存服务，goroutine + `context` 关闭，泛型 handler。
- **Node.js/TS**：`node:http` + 随机端口 + `fetch` 自请求 + `server.close()` + `AbortController`。
- **Java**：`com.sun.net.httpserver.HttpServer` + `record`/`sealed`/`enum`，`stop(0)` 关闭。
- **Rust**：`TcpListener` 绑 `127.0.0.1:0` 手写 HTTP，`match` 路由、零第三方依赖。
- **Python**：`ThreadingHTTPServer` + `urllib.request` 自请求，`shutdown()` 需另起线程调用，标准库零依赖。
