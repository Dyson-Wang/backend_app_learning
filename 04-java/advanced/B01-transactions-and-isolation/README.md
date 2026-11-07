# B01 事务与隔离级别、乐观锁

## 问题复现

两个请求同时给同一账户扣 10 元，都读到 `100`，各自算出 `90` 写回，账户只少了 10 元——丢失更新。

## 错误写法

```java
int b = acc.balance();
acc.balance = b - 10;      // 读到旧值直接写回
```

## 正确写法

```java
synchronized boolean cas(long expectedVersion, int newBalance) {
    if (version != expectedVersion) return false;
    balance = newBalance;
    version++;
    return true;
}
```

对应 SQL：

```sql
UPDATE accounts SET balance = $1, version = version + 1
WHERE id = $2 AND version = $3;
-- 影响 0 行 = 冲突，重读重试
```

JPA 里直接 `@Version private long version;`，冲突时抛 `OptimisticLockException`，用 `@Retryable` 重试。

## 怎么验证

```bash
javac Main.java && java Main
```

预期：`期望 80，实际 90（丢失更新）`，乐观锁那段最终 `期望 80，实际 80`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 互斥 | `sync.Mutex` | 单线程 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 乐观锁 | 手写 version | 手写 version | 手写 version | `@Version` | 手写 version |
| 隔离级别 | `sql.TxOptions` | 驱动参数 | `isolation_level` | `@Transactional(isolation=)` | `SET TRANSACTION` |
