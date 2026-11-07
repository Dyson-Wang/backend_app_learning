# F02 重试风暴

## 问题复现

下游超时，100 个客户端同时重试 3 次，100 QPS 变 400 QPS，打在已过载的服务上，恢复时间从 30 秒变 5 分钟。

## 错误写法

```python
for _ in range(3):
    if call():
        return          # 立即重试，无退避无抖动
```

## 正确写法

```python
for i in range(max_attempts):
    if call():
        return
    time.sleep(2 ** i * 0.1 + random.random() * 0.1)   # 退避 + 抖动
```

## 怎么验证

```bash
python3 main.py
```

预期：`400 QPS` → 退避后 `133 QPS` → 熔断后 `0 QPS`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 退避库 | `backoff` | `p-retry` | **`tenacity`** | Resilience4j | `backoff` crate |
| 熔断 | `gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |
| 陷阱 | 重试非幂等写 | 同 | 同 | 同 | 同 |
