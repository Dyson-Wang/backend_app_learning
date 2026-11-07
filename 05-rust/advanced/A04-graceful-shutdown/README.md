# A04 优雅关闭、健康检查、上游摘除

## 问题复现

滚动发布直接 `kill -9`，在途请求被切断 → 502；后台任务跑到一半，状态不一致。

## 错误写法

```rust
std::process::exit(0);      // 立刻退，在途 handler 全丢
```

## 正确写法

两件事分开：**停止接收**（关 listener）与**等待在途**（join 工作线程 / 等计数归零）。

```rust
stop.store(true, Ordering::SeqCst);
acceptor.join()?;           // listener 被 drop，端口关闭
```

生产用 tokio 时：

```rust
axum::serve(listener, app)
    .with_graceful_shutdown(async { tokio::signal::ctrl_c().await.ok(); })
    .await?;
```

## 怎么验证

```bash
cargo run --quiet
```

预期：`body = slow done，关闭耗时 = 5ms 左右`（只是停止 accept 的耗时），随后 `新请求 -> 被拒绝 ConnectionRefused`。

注意「关闭耗时」量的是停止接收的耗时，不是等在途请求的耗时——在途请求是主线程自己读回来的。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 触发 | `signal.NotifyContext` | `process.on('SIGTERM')` | `signal.signal` | shutdown hook | `tokio::signal` |
| 等待在途 | `srv.Shutdown(ctx)` | `server.close()` | `server.shutdown()` | `server.stop(grace)` | `with_graceful_shutdown` |
| 易错点 | 忘传 ctx | 忘 `unref()` | 主线程死锁 | 忘设 executor | 阻塞 accept 无法中断 |
