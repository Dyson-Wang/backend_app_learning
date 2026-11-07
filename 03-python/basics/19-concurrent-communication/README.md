# 19 · Python · 并发通信：channel / 队列 / 同步原语

## 这个例子讲什么

`queue.Queue(maxsize)` 有界队列的生产消费与背压、毒丸哨兵广播关闭信号、多消费者扇出扇入、`task_done`/`join`、`Event`/`Semaphore`/`Barrier` 同步原语、`Lock` 保护共享状态，以及 `ThreadPoolExecutor` + `as_completed` 等待多任务。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/19-concurrent-communication/main.py
```

## 关键点

- **有界队列 `maxsize=4`**：`put` 在满时阻塞，天然形成背压；`Queue` 自身线程安全，但业务共享 `list` 仍需 `Lock`。
- **关闭信号**：给每个消费者投一个 `None` 毒丸，消费者拿到即 `break`；用「N 个哨兵」而非「消费者接力重投」，避免 `task_done()` 计数失衡导致 `join()` 永不返回。
- `task_done()` 必须和每次 `get()` 配对，`q.join()` 等的是未完成计数归零。
- `Event` 广播状态、`Semaphore` 限流、`Barrier` 让 N 个线程互相等待到齐；都要带 `timeout` 防止死锁。
- `as_completed` 按**完成顺序**产出 `Future`，`pool.map` 才保序；二者都属「等待多个任务」。

## 和其他语言的差异

- **Go**：`chan`（可有界/无界）、`select`、`close(ch)` 广播、`sync.WaitGroup`、`errgroup`。
- **Node.js/TS**：`worker_threads` / `MessageChannel` 传消息、`EventEmitter`、异步队列，单线程为主。
- **Java**：`BlockingQueue`、`ExecutorService` / `CompletableFuture`、`CountDownLatch` / `CyclicBarrier`。
- **Rust**：`std::sync::mpsc` channel、`Arc<Mutex<T>>`、`Condvar`、`Barrier`。
- **Python**：`queue.Queue(maxsize)`、`task_done/join`、毒丸哨兵、`Condition`/`Event`/`Semaphore`/`Barrier`、`concurrent.futures`。
