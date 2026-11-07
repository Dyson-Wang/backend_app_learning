# A04 优雅关闭、健康检查、上游摘除

## 问题复现

滚动发布时直接 `kill -9`，正在处理的请求全部被切断，客户端看到 502。更隐蔽的是：进程已退出但连接池里的 SQL 还在跑，数据可能写一半。

## 错误写法

```js
process.on('SIGTERM', () => process.exit(0));   // 立刻死，在途请求全丢
```

## 正确写法

三步：停止接新连接 → 等在途请求跑完（带超时）→ 退出。

```js
const closed = new Promise((resolve) => server.close(resolve));
const force = setTimeout(() => {
  server.closeAllConnections();
  process.exit(1);
}, 10_000).unref();
await closed;
clearTimeout(force);
```

再配合 Nginx 先摘上游（`down` + `reload`），等一个 `proxy_read_timeout` 再发 SIGTERM。

## 怎么验证

```bash
node index.mjs
```

预期：在途请求返回 `slow done`，`关闭耗时` 约 90-120ms（说明确实等了），随后新请求报 `ECONNREFUSED`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 触发 | `signal.NotifyContext` | `process.on('SIGTERM')` | `signal.signal` | shutdown hook | `ctrlc` crate / 自建 |
| 等待在途 | `srv.Shutdown(ctx)` | `server.close()` | `server.shutdown()` | `server.stop(grace)` | 自建计数 + Condvar |
| 强制兜底 | `ctx` 超时后 `Close()` | `closeAllConnections()` | `server_close()` | `stop(0)` | 自建超时 |
