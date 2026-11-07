# A04 优雅关闭、健康检查、上游摘除

## 问题复现

滚动发布直接 `kill -9`，正在跑的请求被切断，客户端 502；数据库事务可能写一半。进程退出和「请求处理完」之间没有任何协调。

## 错误写法

```python
signal.signal(signal.SIGTERM, lambda *_: sys.exit(0))   # 立刻退
```

## 正确写法

```python
def stop(*_):
    threading.Thread(target=server.shutdown, daemon=True).start()

signal.signal(signal.SIGTERM, stop)
```

注意两点：

1. `shutdown()` **必须**在另一个线程调用，否则会和 `serve_forever` 死锁。
2. `shutdown()` 只是停止 accept，还要 `server_close()` 释放监听 socket。

生产顺序：Nginx 摘上游 → 等一个 `proxy_read_timeout` → 发 SIGTERM → 等 `shutdown()` 返回 → 关连接池。

## 怎么验证

```bash
python3 main.py
```

预期：`body = slow done，关闭耗时 = 90ms 左右`，随后新请求报 `被拒绝 ConnectionRefusedError`（监听 socket 已关）。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 触发 | `signal.NotifyContext` | `process.on('SIGTERM')` | `signal.signal` | shutdown hook | 自建 signal 处理 |
| 等待在途 | `srv.Shutdown(ctx)` | `server.close()` | `server.shutdown()` | `server.stop(grace)` | 自建计数 |
| 易错点 | 忘了传 ctx | 忘了 `unref()` 定时器 | 在主线程调 `shutdown()` 死锁 | 忘了设 executor | 忘了唤醒等待线程 |
