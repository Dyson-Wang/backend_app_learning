# A04 优雅关闭、健康检查、上游摘除

## 问题复现

滚动发布直接 `kill -9`，在途请求被切断 → 客户端 502；事务写一半，数据不一致。

## 错误写法

```java
Runtime.getRuntime().addShutdownHook(new Thread(() -> System.exit(0)));
```

`exit(0)` 不会等在途 handler，JVM 直接收工。

## 正确写法

```java
server.stop(2);   // 停止 accept，最多等 2 秒让在途请求跑完
```

配套顺序：Nginx 摘上游 → 等 `proxy_read_timeout` → 发 SIGTERM → `stop(grace)`。

grace 必须**小于**容器编排的 `terminationGracePeriodSeconds`，否则还没等完就被 SIGKILL。

## 怎么验证

```bash
javac Main.java && java Main
```

预期：`body = slow done，关闭耗时 = 120ms 左右`，随后 `新请求 -> 被拒绝 ConnectException`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 触发 | `signal.NotifyContext` | `process.on('SIGTERM')` | `signal.signal` | shutdown hook | 自建 signal |
| 等待在途 | `srv.Shutdown(ctx)` | `server.close()` | `server.shutdown()` | `server.stop(grace)` | 自建计数 |
| 强制兜底 | `Close()` | `closeAllConnections()` | `server_close()` | `stop(0)` | 自建超时 |
