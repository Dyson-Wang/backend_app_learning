# B03 · Go · 分布式锁

## 问题复现

A 拿到锁后业务卡住超过 TTL，锁自动过期，B 抢到了同一把锁。A 醒来执行 `DEL lock`，把 **B 的锁**删了，C 又能抢到——同一时刻两个进程在改同一份数据。

## 错误写法

```go
ok := redis.SetNX(ctx, key, "1", ttl) // value 是固定字符串
// ...
redis.Del(ctx, key)                   // 解锁不校验持有者
```

## 正确写法

- 加锁：`SET key <random-token> NX PX ttl`，token 必须随机；
- 解锁：**比较 token 再删除**，用 Lua 保证原子：

```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
  return redis.call("del", KEYS[1])
else
  return 0
end
```

- 业务可能超过 TTL 时加看门狗续期（Redisson 的 `watchdog`）；
- 强一致场景用 etcd / ZooKeeper，它们的租约 + revision 比 Redis 更可靠。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/B03-distributed-lock
go run .
```

预期：A 加锁成功、B 失败；推进 150ms 后锁过期；B 加锁成功；`A Del` 后 `holder=""`（误删）；`DelIfMatch(token-A)` 返回 `false` 且 B 的锁还在。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 客户端 | `go-redis` | `ioredis` | `redis-py` | Lettuce / Jedis | `redis-rs` |
| 解锁原子性 | Lua 脚本 | Lua 脚本 | Lua 脚本 | Lua 脚本 | Lua 脚本 |
| 现成方案 | `redsync` | `redlock` | `python-redis-lock` | Redisson | `redislock` |

所有语言的正确做法都一样：随机 token + Lua 比较删除。差别只在客户端 API。Redlock 算法在时钟漂移下有争议，强一致场景优先 etcd。
