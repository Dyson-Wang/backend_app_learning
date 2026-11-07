# A05 · Python · 背压与并发上限

## 问题复现

`queue.Queue()` 默认无界，生产者来多少收多少，内存持续上涨；`ThreadPoolExecutor` 提交几万个任务同理。本例用 `maxsize=5` 的队列提交 20 个任务，输出「接收 5、拒绝 15」。

## 错误写法

```python
q = queue.Queue()  # 无界
for task in tasks:
    q.put(task)    # 来多少收多少
```

没有上限，下游变慢时上游毫无感知。

## 正确写法

```python
try:
    q.put_nowait(task)
    accepted += 1
except queue.Full:
    rejected += 1  # 满则拒绝
```

- `maxsize` 就是并发上限；
- `put_nowait` 满时立刻抛 `queue.Full`，不阻塞、不堆积；
- 想要「阻塞生产者」就用 `put()`，想要「拒绝」就用 `put_nowait()`，语义要明确选一种。

本例用 `threading.Event` 先挡住消费者，保证提交阶段队列必满，输出确定。

## 怎么验证

```bash
cd backend-roadmap/03-python/advanced/A05-backpressure
python3 main.py
```

预期：`accepted=5 rejected=15 processed=5`，最后一行 `无界队列长度=20`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 有界队列 | `chan` + `cap` | 数组 + 长度判断 | `queue.Queue(maxsize)` | `ArrayBlockingQueue` | `sync_channel` |
| 满时行为 | `select default` | 手动 `if` | `put_nowait` 抛 `Full` | `offer` 返回 `false` | `try_send` 返回 `Full` |
| 阻塞式入队 | `q <- x` | 无 | `put()` | `put()` | `send()` |

Python 的 `queue.Queue` 把阻塞和非阻塞拆成 `put` / `put_nowait` 两个方法，比 Java 的 `put` / `offer` / `add` 三套语义更好记；Go 用 `select` 把两种语义写在同一处。
