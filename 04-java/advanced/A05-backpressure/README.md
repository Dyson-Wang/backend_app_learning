# A05 · Java · 背压与并发上限

## 问题复现

`LinkedBlockingQueue` 默认容量是 `Integer.MAX_VALUE`，等于无界；线程池的 `LinkedBlockingQueue` 工作队列也是无界。上游突增时任务全部入队，内存上涨、延迟飙升，最终 OOM。本例用 `ArrayBlockingQueue(5)` 提交 20 个任务，输出「接收 5、拒绝 15」。

## 错误写法

```java
BlockingQueue<Integer> q = new LinkedBlockingQueue<>(); // 无界
q.put(task); // 永远不阻塞、永远不拒绝
```

`put` 在无界队列上永不阻塞，背压完全失效。

## 正确写法

```java
if (queue.offer(task)) accepted++;
else rejected++; // 满则拒绝
```

- `ArrayBlockingQueue(capacity)` 是硬上限；
- `offer` 满时返回 `false`，`put` 会阻塞，`add` 会抛异常——三种语义必须明确选；
- 生产上把 `rejected` 映射成 429，并给队列长度打点。

本例用 `CountDownLatch` 先挡住消费者，保证提交阶段队列必满，输出确定。

## 怎么验证

```bash
cd backend-roadmap/04-java/advanced/A05-backpressure
javac Main.java && java Main
```

预期：`accepted=5 rejected=15 processed=5`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 有界队列 | `chan` + `cap` | 数组 + 长度判断 | `queue.Queue(maxsize)` | `ArrayBlockingQueue` | `sync_channel` |
| 满时行为 | `select default` | 手动 `if` | `put_nowait` 抛 `Full` | `offer` 返回 `false` | `try_send` 返回 `Full` |
| 三套语义 | 阻塞 / 非阻塞 | 只有手动判断 | `put` / `put_nowait` | `put` / `offer` / `add` | `send` / `try_send` |

Java 的 `ThreadPoolExecutor` 默认用无界队列，这是线上最常见的 OOM 来源之一；显式传 `ArrayBlockingQueue` + 拒绝策略（`AbortPolicy` / `CallerRunsPolicy`）才算真正做了背压。
