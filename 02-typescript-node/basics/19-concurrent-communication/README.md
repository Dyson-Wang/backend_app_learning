# 19 · Node.js · 并发通信：线程 / 队列 / 同步原语

## 这个例子讲什么

两部分覆盖并发通信的全部要点：第一部分用 `worker_threads` 起真线程跑扇出任务、`postMessage` 传消息、`SharedArrayBuffer` + `Atomics.add` 保护共享状态、`Promise.all` 扇入汇总；第二部分用 `EventEmitter` 实现一个有界异步队列，演示生产消费、容量背压、关闭信号传播与等待多个任务。worker 源码用 `data:` URL 内联，避免多文件。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/19-concurrent-communication/index.mjs
```

只用 `node:worker_threads` / `node:events`，零依赖、无长睡、瞬间结束。

## 关键点

- `worker_threads` 每个 worker 是**独立 V8 线程 + 独立事件循环**，靠 `parentPort.postMessage` / `workerData` 结构化克隆传消息，不共享普通内存。
- 用 `new Worker(new URL(dataUrl))` 内联 worker 源码，worker 内照样 `import` `node:` 内置模块；监听 `exit` 等它自然退出（退出码 0）。
- **共享状态保护**：`SharedArrayBuffer` 是唯一跨线程共享的内存，必须用 `Atomics.add/load` 做原子操作，避免数据竞争；这是 JS 版的「锁-free 同步原语」。
- 扇出 = 把数据切片交给 N 个 worker；扇入 = 主线程 `Promise.all` 收齐结果；`message` 带回部分和，`exit` 确认收尾。
- 单线程里的 `Map`/数组 `push` 天然无并发写坏，**不需要锁**；只有跨线程共享内存才需要 `Atomics`（与 Go 的 `sync.Mutex` 形成对照）。
- 有界队列：`push` 超过 `capacity` 时返回 pending Promise，生产者 `await` 即被「背压/阻塞」；`take` 在空队列时等待；`close()` 置关闭标志并唤醒所有等待的消费者（返回 `undefined` 作为结束哨兵）。
- `close()` 是**关闭信号传播**：所有 `take()` 等待者被唤醒、消费者循环退出，再由 `Promise.all` 统一收尾。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 真线程 | goroutine（M:N） | `threading`（受 GIL 限制） | `Thread` / 虚拟线程 | `std::thread` | `worker_threads` |
| 消息传递 | channel | `queue.Queue` | `BlockingQueue` | `std::sync::mpsc` | `postMessage` / `MessageChannel` |
| 共享内存 | `sync.Mutex`/atomic | `Lock` / `multiprocessing.Value` | `synchronized`/`AtomicInteger` | `Arc<Mutex<T>>`/`AtomicUsize` | `SharedArrayBuffer` + `Atomics` |
| 等待多个 | `WaitGroup`/`errgroup` | `concurrent.futures` | `CompletableFuture.allOf` | `join!`/线程 join | `Promise.all` |
| 有界队列背压 | 带缓冲 channel 满则阻塞 | `queue.Queue(maxsize)` | `ArrayBlockingQueue` | 无标准（用 crossbeam） | 手写 `BoundedQueue` |
| 关闭信号 | `close(channel)` + `range` | `task_done`/sentinel | `poison`/中断 | `drop(sender)` | `queue.close()` + `undefined` 哨兵 |

Go 的 channel 天然有界且 `close` 可被 `range` 感知，是最接近本课的 19 课模型；Python 的 `queue.Queue(maxsize)` 同样内置背压。Node 单线程事件循环通常不需要同步原语，只有用了 `worker_threads` + `SharedArrayBuffer` 才需要 `Atomics`——这是它与 Rust/Java「默认就上锁」最大的区别。Rust 靠 `Arc<Mutex<T>>` + mpsc 在编译期保证线程安全，JS 则完全靠程序员自觉，没有借用检查兜底。
