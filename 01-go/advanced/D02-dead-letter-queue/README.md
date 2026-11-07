# D02 死信队列与重试策略

## 问题复现

一条消息因为字段缺失（永久错误）被反复重投，卡在队头，后面的消息全部堵住；或者重试次数无上限，把下游打挂。

## 错误写法

```go
for {                       // 无上限重试
    if err := handle(m); err == nil { break }
}
```

## 正确写法

先分类，再定策略：

```go
if isPermanent(err) { return toDLQ(m) }        // 400 / 校验失败，不重试
if attempt >= maxAttempts { return toDLQ(m) }  // 重试耗尽
sleep(backoff(attempt))                        // 10ms, 20ms, 40ms
```

| 错误类型 | 例子 | 策略 |
|---|---|---|
| 永久 | 字段缺失、400 | 直接 DLQ |
| 可重试 | 超时、503 | 指数退避，最多 N 次 |
| 未知 | 空指针 | 有限重试后 DLQ + 告警 |

## 怎么验证

```bash
go run .
```

预期：`成功 = 3，死信 = 2，总投递 = 9`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| DLQ 实现 | 另建 topic | 同 | 同 | RabbitMQ DLX | 同 |
| 退避 | `time.Sleep` | `setTimeout` | `time.sleep` | `Thread.sleep` | `thread::sleep` |
| 幂等要求 | 必须 | 必须 | 必须 | 必须 | 必须 |
