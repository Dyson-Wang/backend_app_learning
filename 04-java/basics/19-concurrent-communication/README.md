# 19 · Java · 并发通信：队列 / 同步原语

## 这个例子讲什么

`BlockingQueue` 有界队列与阻塞、生产消费模型、毒丸（poison pill）关闭信号传播、扇出扇入（`invokeAll`）、等待多个任务（`Future.get` / `CountDownLatch`）、共享状态保护。

## 怎么跑

```
cd backend-roadmap/04-java/basics/19-concurrent-communication
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **`ArrayBlockingQueue(4)` 有界**：`put` 满时阻塞、`poll` 空时按超时等待，天然实现生产者-消费者背压；示例用 `poll(1s)` + 显式循环上限，杜绝死循环。
- **关闭信号传播**：生产者完成后往队列放一个 `POISON` 特殊值，消费者读到就退出；这对应 Go 里 `close(channel)` 的语义。
- **扇出扇入**：`invokeAll` 一次提交多个 `Callable` 并汇总每个 `Future.get()`；`CountDownLatch.await()` 等待一组任务归零。
- **共享状态**：`Collections.synchronizedList` 与 `AtomicInteger` 保护共享数据；示例把结果汇总成计数/求和再打印，输出稳定可复现。
- 用完 `ExecutorService.shutdown()`，否则非守护线程会拖住 JVM 不退出。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 消息通道 | `BlockingQueue` | `chan`（可关闭） | `MessageChannel`/EventEmitter | `queue.Queue` | `std::sync::mpsc` |
| 有界/背压 | `ArrayBlockingQueue(容量)` | 带缓冲 channel | 手写队列 | `Queue(maxsize)` | `sync_channel` |
| 关闭信号 | 毒丸 / `shutdown` | `close(ch)` + `range` | 自定义事件 | `Sentinel`/`task_done` | `drop(tx)` 通道关闭 |
| 等待多任务 | `invokeAll`/`CountDownLatch` | `sync.WaitGroup` | `Promise.all` | `concurrent.futures` | `join`/`Barrier` |
| 共享状态 | `AtomicInteger`/锁 | `sync.Mutex` | 单线程无需锁 | `Lock`（GIL） | `Arc<Mutex<T>>` |
