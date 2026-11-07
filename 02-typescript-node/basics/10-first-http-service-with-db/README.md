# 10 · Node.js · 第一个 HTTP 服务 + 数据层

## 这个例子讲什么

和 Go 版 `basics/10` 功能完全一致的最小 CRUD 服务：路由、JSON、统一错误码、`X-Request-ID`、显式超时、优雅关闭。**同一编号横向对照是这套库的正确读法。**

## 怎么跑

```
cd backend-roadmap/02-typescript-node
PORT=8081 node basics/10-first-http-service-with-db/index.mjs

# 另开终端
curl -i -H 'X-Request-ID: brm-demo-1' http://localhost:8081/healthz
curl -i -X POST http://localhost:8081/users -d '{"name":"alice"}'
curl -i http://localhost:8081/users
curl -i http://localhost:8081/users/999
# Ctrl-C 观察优雅关闭日志
```

## 关键点

- **零依赖**：只用 `node:http`。真实项目会用 Fastify / Express / Hono，但先把标准形态看清。
- `Map` 在单线程里天然不会被并发写坏——这是 Node 少有的「不需要锁」场景，但多实例部署时状态依然不共享。
- `for await (const chunk of req)` 消费请求流，同时检查大小上限；**别用同步的 `fs.readFileSync` 或一次性大 JSON 解析**，那会阻塞事件循环。
- `try/catch/finally` 里统一打访问日志并带 `request_id`；生产要把 `request_id` 放进 `AsyncLocalStorage`，避免一层层手动传参。
- `server.close()` 停止接收新连接并等在途请求结束；用 `setTimeout(...).unref()` 做强制退出兜底，且 `unref()` 让这个定时器不阻止进程正常退出。
- **接数据库**：把 `Map` 换成 `pg` / `postgres` 的连接池（`Pool`），注意连接池大小要和事件循环的并发模型一起评估，见 `advanced/C01`。

## 和其他语言的差异

| 能力 | Go | Node.js |
|---|---|---|
| 路由 | `http.ServeMux` | 手工 `new URL()` + 分支（框架才提供 router） |
| 中间件 | `func(Handler) Handler` | 包装函数 / `AsyncLocalStorage` |
| 数据层抽象 | `Store` 接口 | 无接口，靠对象形状（TS 里用 `interface`） |
| 共享状态 | `sync.RWMutex` | 单线程 `Map`，通常不需要锁 |
| 优雅关闭 | `signal.NotifyContext` + `Shutdown` | `process.on` + `server.close` |
| 取消 | `context.Context` | `AbortSignal` |

最大的差别不是语法，是**一个同步阻塞在 Go 里只影响一个 goroutine，在 Node 里会拖垮整个进程**。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
