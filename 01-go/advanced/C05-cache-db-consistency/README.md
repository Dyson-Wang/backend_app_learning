# C05 缓存与 DB 双写一致性

## 问题复现

更新用户信息时「删缓存 + 写库」两步之间，来了一个读请求：缓存已删 → 回源读到**旧值** → 写回缓存。写请求随后把库更新成新值，但缓存里永远留着旧值。

## 错误写法

```go
cache.Del(key)      // 先删
db.Update(key, v)   // 中间有读请求回填旧值
```

## 正确写法

Cache-Aside 的标准顺序是**先写库、再删缓存**，并且延迟再删一次：

```go
db.Update(key, v)
cache.Del(key)
go func() {
    time.Sleep(500 * time.Millisecond)
    cache.Del(key)          // 双删，兜住迟到的回填
}()
```

更强的做法：订阅 binlog（Canal / Debezium）异步删缓存，把「删除」变成可重试的可靠事件。

## 怎么验证

```bash
go run .
```

预期：前两种都是 `是否脏 = true`，双删后 `缓存命中 = false`，下次回源拿到 `v2`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 缓存客户端 | `go-redis` | `ioredis` | `redis-py` | Lettuce | `redis-rs` |
| 延迟双删 | `time.AfterFunc` | `setTimeout` | `threading.Timer` | `ScheduledExecutor` | 线程 + sleep |
| 强一致替代 | 读走主库 | 同 | 同 | 同 | 同 |
