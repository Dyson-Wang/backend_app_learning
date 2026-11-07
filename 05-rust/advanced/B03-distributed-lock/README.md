# B03 分布式锁

## 问题复现

A 的业务超过锁 TTL，锁过期，B 抢到锁；A 执行完直接 `DEL`，删掉了 B 的锁，第三个请求又能进来。

## 错误写法

```rust
let _guard = lock(&r, key, ttl)?;
do_work();
r.delete(key);      // 不看 token
```

## 正确写法

```rust
if !r.set_nx_px(key, &token, ttl_ms) { return Err(Busy); }
do_work();
r.del_if_match(key, &token);      // 生产用 Lua 保证原子
```

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
else
  return 0
end
```

## 怎么验证

```bash
cargo run --quiet
```

预期：`A del_if_match(token-A) = false`（B 的锁保住），`B del_if_match(token-B) = true`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 客户端 | `go-redis` | `ioredis` | `redis-py` | Lettuce | `redis-rs` |
| 原子解锁 | `Eval` | `eval` | `register_script` | `RedisScript` | `Script::invoke` |
| 续期 | goroutine | `setInterval` | 线程 | watchdog | 线程 |
