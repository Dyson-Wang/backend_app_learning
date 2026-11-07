# D02 死信队列与重试策略

## 问题复现

永久错误的消息被反复重投，卡住队头；或重试无上限，把下游打挂。

## 错误写法

```java
while (true) {
    if (handle(m)) break;      // 无上限
}
```

## 正确写法

```java
if (m.permanent) return toDlq(m);
if (attempt >= maxAttempts) return toDlq(m);
Thread.sleep(10L << (attempt - 1));
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：`成功 = 3，死信 = 2，总投递 = 9`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| DLQ | 另建 topic | 同 | 同 | RabbitMQ DLX | 同 |
| 退避 | `time.Sleep` | `setTimeout` | `time.sleep` | `Thread.sleep` | `thread::sleep` |
| 幂等 | 必须 | 必须 | 必须 | 必须 | 必须 |
