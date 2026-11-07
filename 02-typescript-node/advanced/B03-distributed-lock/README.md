# B03 · Node.js · 分布式锁

## 问题复现

A 拿到锁后业务卡住超过 TTL，锁自动过期，B 抢到同一把锁。A 醒来 `DEL lock`，删掉的是 **B 的锁**，C 又能抢到——两个进程同时改同一份数据。

## 错误写法

```js
await redis.set(key, '1', 'NX', 'PX', ttl); // 固定 value
// ...
await redis.del(key);                        // 不校验持有者
```

## 正确写法

- 加锁 `SET key <random-token> NX PX ttl`；
- 解锁用 Lua 比较 token：

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
end
return 0
```

- 业务可能超时就加看门狗续期；
- 强一致场景用 etcd / ZooKeeper。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/B03-distributed-lock
node index.mjs
```

预期：A 成功、B 失败；150ms 后锁过期；B 成功；`A del` 后 `holder=""`；`delIfMatch(token-A)` 返回 `false` 且 B 的锁还在。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 客户端 | `go-redis` | `ioredis` | `redis-py` | Lettuce | `redis-rs` |
| 原子解锁 | Lua | Lua | Lua | Lua | Lua |
| 现成库 | `redsync` | `redlock` | `python-redis-lock` | Redisson | `redislock` |

Node 的 `ioredis` 支持 `eval` 直接跑 Lua；`redlock` 库封装了多节点 Redlock，但 Redlock 在时钟漂移下有争议，不要把它当成强一致方案。
