# B01 · Node.js · 事务与隔离级别、乐观锁

## 问题复现

两个事务都读到 `balance=100`，各扣 10 后写回 90。数据库里应该是 80，实际 90——第二次写入覆盖了第一次。读已提交隔离级别下，这种交错完全合法。

## 错误写法

```js
const { balance } = await read();
await write(balance - 10); // 读和写之间没有版本校验
```

## 正确写法

- **悲观锁**：`SELECT ... FOR UPDATE`；
- **乐观锁**：`UPDATE ... SET balance=$1, version=version+1 WHERE id=$2 AND version=$3`，`rowCount === 0` 就重试；
- Node 里事务必须用同一个连接，`pool.query` 每次可能拿不同连接，事务会失效。

本例用内存 `cas(expectedVersion, newBalance)` 模拟乐观锁。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/B01-transactions-and-isolation
node index.mjs
```

预期：无锁 `实际 90`；乐观锁 `T2 用旧 version 写入失败`，重试后 `balance=80`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 事务 API | `database/sql` `Tx` | `pg` client | SQLAlchemy `Session` | JDBC / `@Transactional` | `sqlx` transaction |
| 乐观锁 | 手写 version | 手写 | `version_id_col` | JPA `@Version` | 手写 |
| 连接绑定 | `Tx` 对象 | 必须用同一 client | Session 绑定 | 事务绑定线程 | `Transaction` 借用 |

Node 最容易踩的坑是「事务里用了 pool.query」：每次可能拿不同连接，`BEGIN` 和 `COMMIT` 落在两个连接上，事务等于没开。必须 `const client = await pool.connect()` 再全程用 `client`。
