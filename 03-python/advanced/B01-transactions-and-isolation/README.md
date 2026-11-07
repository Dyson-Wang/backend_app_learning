# B01 事务与隔离级别、乐观锁

## 问题复现

两个请求同时给同一个账户扣 10 元。没有并发控制时，两边都读到 `100`，各自算出 `90` 写回，账户只少了 10 元而不是 20 元——这叫**丢失更新**（lost update）。

## 错误写法

```python
balance, _ = acc.read()
acc.balance = balance - amount   # 读到旧值直接写回
```

`naive_transfer` 调用两次，`100 -> 90 -> 90`，第二次覆盖了第一次的结果。

## 正确写法

乐观锁：读的时候带上 `version`，写的时候用 `cas(expected_version, new_balance)` 校验。版本不匹配就重试。

```python
def optimistic_debit(acc, amount, max_retry=10):
    for attempt in range(max_retry + 1):
        balance, version = acc.read()
        if acc.cas(version, balance - amount):
            return True, attempt
    return False, max_retry + 1
```

对应 SQL：

```sql
UPDATE accounts SET balance = $1, version = version + 1
WHERE id = $2 AND version = $3;
-- 影响行数 0 表示冲突，需要重读重试
```

## 怎么验证

```bash
python3 main.py
```

预期：

```text
期望 80，实际 90（丢失更新）
T1 ok=True attempts=0，balance=90 version=2
T2 用旧 version=1 写入：False
T2 重试 ok=True attempts=0，balance=80 version=3
期望 80，实际 80
```

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 互斥 | `sync.Mutex` | 单线程无需锁 | `threading.Lock` | `synchronized` / `ReentrantLock` | `Mutex<T>` |
| 乐观锁表达 | 手写 version 比较 | 手写 version 比较 | 手写 version 比较 | JPA `@Version` | 手写 version 比较 |
| 重试 | for 循环 | for 循环 | for 循环 | Spring Retry | for 循环 |
