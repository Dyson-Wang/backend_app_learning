# D03 顺序性保证

## 问题复现

订单的三个事件 `create → pay → ship` 被三个消费者并发处理。`ship` 先到，状态机里 `PAID` 还没出现，直接报错；或者更糟——没有状态机，脏数据就这么写进去了。

## 错误写法

```go
// 多个消费者抢同一个队列，谁先拿到谁处理
for msg := range ch {
    go handle(msg)      // 顺序彻底丢了
}
```

## 正确写法

两条路，选一条：

1. **分区内串行**：`hash(key) % N` 决定分区，每个分区只有一个消费者。
2. **消费端按 key 串行**：内存里维护 `key → 队列`，同 key 排队，不同 key 并行。

```go
h := 0
for _, ch := range k { h = (h*31 + int(ch)) % 4 }
```

## 怎么验证

```bash
go run .
```

预期：乱序时 `ship 在 INIT 状态下不合法`；顺序时 `最终状态 = SHIPPED`；`order-42` 两次都落到同一分区。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 分区 | 手写 hash | 同 | 同 | 同 | 同 |
| 串行手段 | 单 goroutine | 单线程天然串行 | 单线程 | 单线程 executor | 单任务 |
| MQ 支持 | Kafka partition key | 同 | 同 | 同 | 同 |
