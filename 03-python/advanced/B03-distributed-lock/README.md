# B03 分布式锁

## 问题复现

单机 `threading.Lock` 管不了跨进程。多副本部署时，两个实例同时处理同一个订单，必须用外部存储做互斥——通常是 Redis。

坑在**解锁**：A 的业务执行超过锁的 TTL，锁自动过期，B 抢到了锁；A 执行完直接 `DEL`，把 B 的锁删了。第三个请求 C 又能进来，互斥彻底失效。

## 错误写法

```python
r.set(key, token, nx=True, px=ttl)
try:
    do_work()
finally:
    r.delete(key)      # 不看 token，可能删掉别人的锁
```

## 正确写法

加锁必须带随机 token，解锁必须「先比较再删除」，且这两步要原子（Lua）：

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
else
  return 0
end
```

```python
r.set(key, token, nx=True, px=ttl)
try:
    do_work()
finally:
    r.del_if_match(key, token)
```

业务可能超过 TTL 时用看门狗续期（Redlock 的续租，或 `SET` 重新 PX）。

## 怎么验证

```bash
python3 main.py
```

预期看到 `A delIfMatch(token-A) = False`（保护住了 B 的锁），`B delIfMatch(token-B) = True`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 客户端 | `go-redis` | `ioredis` | `redis-py` | Lettuce / Jedis | `redis-rs` |
| 原子解锁 | `Eval` 跑 Lua | `eval` 跑 Lua | `register_script` | `RedisScript` | `Script::invoke` |
| 续期 | 自建 goroutine | `setInterval` | 自建线程 | Redisson watchdog | 自建线程 |
