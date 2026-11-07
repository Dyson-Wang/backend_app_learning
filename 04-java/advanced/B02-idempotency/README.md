# B02 幂等键与重复请求

## 问题复现

用户连点两次提交，或网关超时重试，同一请求到达两次 → 两张订单、两次扣款。

## 错误写法

```java
if (done.containsKey(key)) return done.get(key);
Order o = create(key);      // 查和写之间有窗口
```

## 正确写法

用 `putIfAbsent` 做占位，后到的请求 `await` 第一个的结果：

```java
CountDownLatch prev = inflight.putIfAbsent(key, mine);
if (prev != null) {
    prev.await();
    return done.get(key);
}
```

落库版：

```sql
INSERT INTO orders (order_no, idempotency_key, amount)
VALUES (?, ?, ?)
ON CONFLICT (idempotency_key) DO NOTHING;
```

并发下捕获唯一索引冲突，回查原结果返回——**连响应体一起存**。

## 怎么验证

```bash
javac Main.java && java Main
```

预期：`真实建单 = 1，回放 = 19，订单 id 去重后 = 1`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 占位原语 | `singleflight` | Promise Map | `threading.Event` | `putIfAbsent` + `CountDownLatch` | `Mutex<HashMap>` |
| 落库保证 | 唯一索引 | 唯一索引 | 唯一索引 | 唯一索引 | 唯一索引 |
| 常见坑 | 忘存响应体 | 忘存响应体 | 忘存响应体 | 事务里吞异常 | 忘存响应体 |
