# B03 分布式锁

## 问题复现

A 的业务超过锁 TTL，锁自动过期，B 抢到锁；A 执行完直接 `DEL`，把 B 的锁删了，第三个请求又能进来，互斥失效。

## 错误写法

```java
try {
    doWork();
} finally {
    redis.del(key);        // 不看 token，可能删掉别人的锁
}
```

## 正确写法

```java
String token = UUID.randomUUID().toString();
if (!redis.set(key, token, "NX", "PX", ttl)) return false;
try {
    doWork();
} finally {
    redis.eval(DEL_IF_MATCH_LUA, List.of(key), List.of(token));
}
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
javac Main.java && java Main
```

预期：`A delIfMatch(token-A) = false`（B 的锁保住了），`B delIfMatch(token-B) = true`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 客户端 | `go-redis` | `ioredis` | `redis-py` | Lettuce / Redisson | `redis-rs` |
| 原子解锁 | `Eval` | `eval` | `register_script` | `RedisScript` | `Script::invoke` |
| 续期 | 自建 goroutine | `setInterval` | 自建线程 | Redisson watchdog | 自建线程 |
