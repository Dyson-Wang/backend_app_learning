# B04 · Node.js · 竞态、死锁与调试

## 问题复现

Node 是单线程，所以没有传统数据竞争。但 `await` 会让出执行权，两个请求都读到旧值、再各自写回，结果只加了一次——这是 Node 里最隐蔽的「竞态」。

死锁在 Node 里表现为 Promise 永远不 settle：两个任务互相 `await` 对方，事件循环里再也没有东西能推进它们。

## 错误写法

```js
const v = counter;
await something(); // 让出
counter = v + 1;   // 覆盖了别人写的结果
```

## 正确写法

- 临界区里不要 `await`；必须异步时，用一条 Promise 链把「读-改-写」串起来；
- 死锁用 `Promise.race` + 超时兜底，永远不要无限等；
- 排查用 `--inspect` + async stack traces，或 `diagnostics_channel` 打点。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/B04-race-and-deadlock
node index.mjs
```

预期：`期望 2，实际 1`、`期望 1000，实际 1000`、`结果: timeout`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 竞态本质 | 真并行数据竞争 | 异步交错 | 线程 + GIL | 真并行 | 编译期禁止 |
| 检测工具 | `-race` | async stack / inspector | `faulthandler` | JFR / thread dump | 编译器 |
| 互斥 | `sync.Mutex` | Promise 链 / 队列 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 死锁表现 | goroutine 全阻塞 | Promise 永不 settle | 线程阻塞 | 线程 BLOCKED | 线程阻塞 |

Node 的「锁」其实就是串行化队列；这也是为什么很多 Node 代码用单个 worker 或 `p-queue` 来避免交错，而不是加锁。
