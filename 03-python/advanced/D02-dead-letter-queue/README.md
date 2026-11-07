# D02 死信队列与重试策略

## 问题复现

永久错误的消息被反复重投，卡住队头；或重试无上限，把下游打挂。

## 错误写法

```python
while True:
    if handle(m):
        break
```

## 正确写法

```python
if msg.permanent:
    return to_dlq(msg)
if attempt >= max_attempts:
    return to_dlq(msg)
time.sleep(0.01 * 2 ** (attempt - 1))
```

## 怎么验证

```bash
python3 main.py
```

预期：`成功 = 3，死信 = 2，总投递 = 9`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| DLQ | 另建 topic | 同 | 同 | RabbitMQ DLX | 同 |
| 退避 | `time.Sleep` | `setTimeout` | `time.sleep` | `Thread.sleep` | `thread::sleep` |
| 幂等 | 必须 | 必须 | 必须 | 必须 | 必须 |
