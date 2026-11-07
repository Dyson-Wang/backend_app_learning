# A05 · Node.js · 背压与并发上限

## 问题复现

Node 单线程不会因为「线程太多」而崩，但无界数组会把内存吃光；更常见的是 `Promise.all(tasks)` 一次性发起几万个请求，把下游和自身的事件循环一起拖死。本例用容量 5 的数组队列提交 20 个任务，输出「接收 5、拒绝 15」。

## 错误写法

```js
const queue = [];
for (const task of tasks) queue.push(task); // 无界
```

或者：

```js
await Promise.all(tasks.map(handle)); // 一次性打满
```

两种写法都没有上限，上游感受不到下游变慢。

## 正确写法

入口判断长度，满则拒绝：

```js
if (queue.length < capacity) { queue.push(task); accepted += 1; }
else { rejected += 1; }
```

- 容量是硬上限，内存可预测；
- 拒绝立即返回，不排队、不占 Promise；
- 生产上配合 `p-limit` / `fastq` 做并发上限，或直接返回 429。

本例用 `gate` 先挡住消费者，保证提交阶段队列必满，输出确定。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/A05-backpressure
node index.mjs
```

预期：`accepted=5 rejected=15 processed=5`，最后一行 `无界队列长度=20`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 有界队列 | `chan` + `cap` | 数组 + 长度判断 | `queue.Queue(maxsize)` | `ArrayBlockingQueue` | `sync_channel` |
| 满时行为 | `select default` | 手动 `if` | `put_nowait` 抛 `Full` | `offer` 返回 `false` | `try_send` 返回 `Full` |
| 并发上限库 | `semaphore` | `p-limit` | `Semaphore` | `Semaphore` | `tokio::sync::Semaphore` |

Node 没有内建阻塞队列，所以「背压」通常体现为「限制并发 Promise 数量」而不是「阻塞生产者」；这也是 `p-limit` 这类库存在的原因。
