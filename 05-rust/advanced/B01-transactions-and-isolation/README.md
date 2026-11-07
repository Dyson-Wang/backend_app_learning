# B01 事务与隔离级别、乐观锁

## 问题复现

两个请求同时扣 10 元，都读到 `100`，各自算 `90` 写回，账户只少 10 元。

Rust 的借用检查器**不会**帮你防这个：它保证同一时刻没有数据竞争，但两次「读-改-写」在时间上串行，逻辑上仍然是丢失更新。

## 错误写法

```rust
let read1 = acc.balance;
acc.balance = read1 - 10;      // 用过期的读视图
```

## 正确写法

```rust
fn cas(&mut self, expected_version: u64, new_balance: i64) -> bool {
    if self.version != expected_version { return false; }
    self.balance = new_balance;
    self.version += 1;
    true
}
```

```sql
UPDATE accounts SET balance = $1, version = version + 1
WHERE id = $2 AND version = $3;
```

## 怎么验证

```bash
cargo run --quiet
```

预期：`期望 80，实际 90（丢失更新）`，乐观锁那段最终 `期望 80，实际 80`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 互斥 | `sync.Mutex` | 单线程 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 编译期防线 | 无 | 无 | 无 | 无 | 借用检查（只防数据竞争） |
| 乐观锁 | 手写 version | 手写 version | 手写 version | `@Version` | 手写 version |
