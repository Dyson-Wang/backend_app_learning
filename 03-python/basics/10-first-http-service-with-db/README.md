# 10 · Python · 第一个 HTTP 服务 + 数据层

## 这个例子讲什么

和 Go / Node 版 `basics/10` 同样的最小 CRUD，用标准库 `http.server` 实现：线程模型、显式锁、统一错误码、优雅关闭。

## 怎么跑

```
cd backend-roadmap/03-python
PORT=8082 python3 basics/10-first-http-service-with-db/main.py

# 另开终端
curl -i -H 'X-Request-ID: brm-demo-1' http://localhost:8082/healthz
curl -i -X POST http://localhost:8082/users -d '{"name":"alice"}'
curl -i http://localhost:8082/users
curl -i http://localhost:8082/users/999
# Ctrl-C 观察优雅关闭
```

## 关键点

- **`ThreadingHTTPServer` 每连接一个线程**，因此共享的 `_users` 必须用 `threading.Lock` 保护——这一点和 Go 的 `sync.RWMutex` 完全对应，和 Node 的单线程 `Map` 形成对照。
- `HTTP/1.1` 下必须正确设置 `Content-Length`，否则客户端会一直等；`Content-Length` 要按**字节数**算（`ensure_ascii=False` 时中文占 3 字节）。
- `server.shutdown()` 必须从**另一个线程**调用，在 signal handler 里直接调会死锁——这是 Python 标准库特有的坑。
- 标准库 `http.server` 只适合学习和内部工具，**不要上生产**：没有连接池、没有超时控制、错误处理粗糙。生产用 FastAPI/Starlette + uvicorn/gunicorn，`advanced/` 会引入。
- **接数据库**：换成 SQLAlchemy / asyncpg 的连接池，注意 GIL 下同步驱动会占线程、异步驱动才配得上 `asyncio` 的并发能力，见 `advanced/C01`。

## 和其他语言的差异

| 能力 | Go | Python | Node.js |
|---|---|---|---|
| 并发单元 | goroutine | 每连接一个 OS 线程 | 事件循环 + Promise |
| 共享状态 | `sync.RWMutex` | `threading.Lock` | 单线程无需锁 |
| 内存 | GC | 引用计数 + 分代 GC | GC |
| 取消 | `context` | `asyncio.timeout`（异步版） | `AbortSignal` |
| 优雅关闭 | `Shutdown(ctx)` | `server.shutdown()`（另起线程） | `server.close()` |

Python 版最值得注意的反而是「锁」：同一个功能，Go 需要锁、Node 不需要、Python 也需要——因为后两者走的是两个极端（单线程 vs 多线程），只有 Go 是「轻量协程 + 真并行」。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
