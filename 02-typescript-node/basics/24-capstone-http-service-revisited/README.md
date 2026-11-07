# 24 · Node.js · 综合实战：语言特性用回 HTTP 服务

## 这个例子讲什么

把前 23 课的语言特性拼回一个**内存版** HTTP 服务：`node:http` 监听随机端口（`listen(0)`）、自己用 `http.request` 请求自己、覆盖路由 / JSON / 校验 / 统一错误码 / 访问日志，最后 `server.closeAllConnections()` + `server.close()` 自关闭退出。用到的特性包括：枚举常量（`Object.freeze`）、自定义错误类、泛型（JSDoc `@template`）、闭包工厂捕获 `Map`、`for await` 消费请求流、正则路由、`Promise.allSettled` 并发、`performance.now()` 时间差值。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/24-capstone-http-service-revisited/index.mjs
```

只用 `node:http` / `node:perf_hooks`；自请求自关闭，约 0.1 秒结束，退出码 0。

## 关键点

- `listen(0, '127.0.0.1')` 让系统分配随机端口，`server.address().port` 取回；只绑本机，不对外暴露。
- 自请求用 `node:http` 的 `http.request` 包成 Promise；请求头带 `Connection: close`，避免 Node 19+ 默认 keep-alive 让 `server.close()` 一直等在途连接。
- 顺序跑 9 个用例覆盖全部路由与错误分支（建号顺序确定 → 输出可复现），再用 `Promise.allSettled` 并发打 2 个只读请求演示扇入。
- 收尾：`server.closeAllConnections()` 断掉残留连接 → `server.close(cb)` 停止接收 → 回调里清掉看门狗；`setTimeout(...).unref()` 是「不阻止正常退出、仅在卡住时兜底」的标准写法（10 课）。
- 统一错误：自定义 `HttpError extends Error` 带 `status`/`code`，`try/catch/finally` 里把错误转成 `{error:{code,message}}`，`finally` 打访问日志。
- 各课特性落点：枚举/ADT（13）、泛型擦除（14）、闭包捕获 `Map`（16/05）、异步迭代流（09/08）、正则（18）、并发等待（20）、单调时钟差值（23）、鸭子类型响应（06）。
- 输出稳定：request id 用确定性自增序号、不打印随机端口与波动耗时；时间只演示「差值已测量」而非具体数值。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 内存服务 | `net/http/httptest` | `ThreadingHTTPServer` | `HttpServer` | `TcpListener` 手写 | `node:http` + `listen(0)` |
| 随机端口 | `httptest.NewServer` 自动 | 绑 `(':0')` 取端口 | `new HttpServer(..., 0)` | `127.0.0.1:0` | `listen(0)` |
| 自请求 | `server.Client().Get` | 主线程 `urllib` 请求子线程服务 | `HttpClient` | 主线程手写 HTTP 请求 | `http.request` 自请求 |
| 关闭 | `server.Close()` | `shutdown()` | `stop(0)` | drop listener / 进程退出 | `server.close()` |
| 路由 | `ServeMux` | 手写 `do_GET` | 手写 `HttpHandler`/context | 手写解析 | 手写 `URL` + 分支 |
| 错误模型 | `error` 返回值 | 异常 | 异常 | `Result` | `throw` + `try/catch` |

Go 的 `httptest.NewServer` 把「随机端口 + 自请求 + 自动关闭」封装得最完整；Python 需要子线程跑服务、主线程 `shutdown()`；Java 的 `com.sun.net.httpserver.HttpServer` + `stop(0)`；Rust 标准库没有 HTTP 服务端，得手写 `TcpListener` 解析请求。Node 与 Go 一样同进程内即可自请求，但要注意 keep-alive 连接必须先 `closeAllConnections()` 才能干净退出——这是 Node 19+ 默认开启 keep-alive 后最容易踩的坑。
