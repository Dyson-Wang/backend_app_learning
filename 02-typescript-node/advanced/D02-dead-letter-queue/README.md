# D02 死信队列与重试策略

## 问题复现

一条永久错误的消息被反复重投，卡住队头；或重试次数无上限，把下游打挂。

## 错误写法

```js
while (true) { if (await handle(m)) break; }   // 无上限
```

## 正确写法

```js
if (isPermanent(err)) return toDLQ(m);
if (attempt >= maxAttempts) return toDLQ(m);
await sleep(10 * 2 ** (attempt - 1));
```

| 错误类型 | 策略 |
|---|---|
| 永久（400 / 校验） | 直接 DLQ |
| 可重试（超时 / 503） | 指数退避，最多 N 次 |
| 未知 | 有限重试后 DLQ + 告警 |

## 怎么验证

```bash
node index.mjs
```

预期：`成功 = 3，死信 = 2，总投递 = 9`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| DLQ | 另建 topic | 同 | 同 | RabbitMQ DLX | 同 |
| 退避 | `time.Sleep` | `setTimeout` | `time.sleep` | `Thread.sleep` | `thread::sleep` |
| 幂等 | 必须 | 必须 | 必须 | 必须 | 必须 |
