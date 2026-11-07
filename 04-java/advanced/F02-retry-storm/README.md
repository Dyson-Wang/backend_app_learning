# F02 重试风暴

## 问题复现

下游超时，100 个客户端同时重试 3 次，100 QPS 变 400 QPS，打在一个已过载的服务上。

## 错误写法

```java
for (int i = 0; i < 3; i++) {
    if (call()) return;      // 立即重试，无退避无抖动
}
```

## 正确写法

```java
Retry.decorateSupplier(Retry.backoff(3, Duration.ofMillis(100))
        .jitter(0.5), this::call).get();
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：`400 QPS` → 退避后 `133 QPS` → 熔断后 `0 QPS`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 退避库 | `backoff` | `p-retry` | `tenacity` | **Resilience4j** | `backoff` crate |
| 熔断 | `gobreaker` | `opossum` | `pybreaker` | Resilience4j | `failsafe` |
| 陷阱 | 重试非幂等写 | 同 | 同 | `@Retryable` 加在事务内 | 同 |
