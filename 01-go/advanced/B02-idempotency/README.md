# B02 · 幂等键与重复请求

## 问题复现

客户端网络抖动重发、网关重试、MQ 至少一次投递——重复请求是**常态**不是异常。没有幂等，重复 POST 就会创建两笔订单、扣两次款。

本例 20 个并发请求带同一个 `Idempotency-Key`，放大 5ms 业务耗时制造竞态，验证最终只创建 1 笔订单，其余请求回放同一个结果。

## 错误写法

**错法 1 · 先查再写（check-then-act）**：

```go
if !exists(key) {      // 两个并发请求都查到「不存在」
    createOrder(key)   // 于是都创建了
}
```

进程内加锁能修单实例，但多实例部署下毫无作用。**唯一可靠的兜底是数据库唯一约束**。

**错法 2 · 用业务字段当幂等键但不加唯一索引**：靠应用层保证，一旦有并发或重放就漏。

**错法 3 · 幂等记录不存响应结果**：重复请求只返回「已处理」，客户端拿不到订单号，只能再查一次；正确做法是把首次响应一起存下来并回放。

**错法 4 · 永不过期**：幂等表无限膨胀，需要按业务周期清理（例如保留 24h~7d），且清理窗口要大于客户端最大重试周期。

**错法 5 · 把 in-flight 当成失败**：第一个请求还在处理，重复请求直接返回 500，客户端又重试，形成新的风暴。应返回 409 + Retry-After 或等待后回放。

## 正确写法

核心是把「查 + 占位」变成**一个原子操作**：本例 `Claim()` 在锁内完成，对应 PostgreSQL 的

```sql
INSERT INTO orders (order_no, account_id, amount, idempotency_key)
VALUES ('ORD-1', 1, 100, 'client-key')
ON CONFLICT (idempotency_key) DO NOTHING
RETURNING id, status;
```

`RETURNING` 为 0 行说明已被别人占位，再去查已有结果回放。表结构见 `../../../00-foundations/sql-and-postgres/schema.sql`。

三种状态要分开处理：**首次执行**（201）、**已完成回放**（200 + 原响应）、**执行中**（409 + Retry-After 或等待）。幂等键应由客户端生成（UUID / 业务单号），服务端校验格式和长度，并把它纳入签名防篡改。

数据库写入与业务操作要放在**同一个事务**里，否则「订单创建成功、幂等记录没写上去」仍然会重复。

## 怎么验证

```
cd backend-roadmap/01-go
go run ./advanced/B02-idempotency
# 期望：201:1 200:19 distinct orders:1，缺 key 返回 400，selftest PASS

go run -race ./advanced/B02-idempotency   # 验证没有数据竞争
```

相关底座概念见 `../../../00-foundations/`。做完把差异结论补到 `../../../99-cross-language/scenario-matrix.md`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 占位原语 | `singleflight` | Promise Map | `threading.Event` | `putIfAbsent` | `Mutex` + `Condvar` |
| 落库保证 | 唯一索引 + `ON CONFLICT` | 同 | 同 | 同 | 同 |
| 常见坑 | 忘存响应体 | 忘存响应体 | 忘存响应体 | 事务里吞异常 | 忘存响应体 |
