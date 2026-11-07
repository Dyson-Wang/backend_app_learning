# B02 幂等键与重复请求

## 问题复现

用户点「提交订单」没反应，又点了一次；或者网关超时后自动重试。同一个请求到达两次，创建了两张订单、扣了两次款。**至少一次投递 + 客户端重试是常态，不是异常。**

## 错误写法

```js
if (orders.has(key)) return orders.get(key);
const order = await createOrder(key);   // 查和写之间有窗口
```

20 个并发请求会全部穿过 `has()` 判断，建出 20 张单。

## 正确写法

「查」和「占位」必须是同一步。内存版用 `inflight` Promise 占位；落库版靠唯一索引：

```sql
INSERT INTO orders (order_no, idempotency_key, amount)
VALUES ($1, $2, $3)
ON CONFLICT (idempotency_key) DO NOTHING
RETURNING id;
-- 没有返回行 = 已存在，回查原结果返回
```

还要存**响应体**，这样重放时返回一模一样的结果（包括状态码）。

## 怎么验证

```bash
node index.mjs
```

预期：`真实建单 = 1，回放 = 19，订单 id 去重后 = 1`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 占位原语 | `sync.Map` + `singleflight` | `inflight` Map | `threading.Event` | `ConcurrentHashMap.putIfAbsent` | `Mutex<HashMap>` |
| 落库保证 | `ON CONFLICT` | `ON CONFLICT` | `ON CONFLICT` | 唯一索引 + 捕获异常 | `ON CONFLICT` |
| 常见坑 | 忘记存响应体 | 忘记存响应体 | 忘记存响应体 | 事务里吞异常 | 忘记存响应体 |
