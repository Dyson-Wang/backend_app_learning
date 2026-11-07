# 09 · Java · 并发原语入门

## 这个例子讲什么

`Thread` + `Runnable`、`ExecutorService` 线程池、`AtomicInteger`（CAS）、`synchronized` 监视器、`volatile` 可见性、`CountDownLatch` 等待多任务。

## 怎么跑

```
cd backend-roadmap/04-java/basics/09-concurrency-primitives
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **`synchronized`** 基于对象监视器，修饰实例方法时锁的是 `this`；`AtomicInteger` 用 CAS 无锁，粒度更细。
- **`volatile`** 只保证可见性与禁止指令重排，**不保证 `i++` 原子性**（那是读-改-写三步），复合操作仍需锁或原子类。
- `CountDownLatch.await()` 让主线程等 N 个任务归零，比 `sleep` 精确且可复现；`ExecutorService` 用完必须 `shutdown()`，否则非守护线程会拖住 JVM 不退出。
- JDK 21+ 可用 `Executors.newVirtualThreadPerTaskExecutor()` 获得廉价虚拟线程（`advanced` 重点）。
- 示例把线程结果汇总成最终计数再打印，避免交错日志导致输出不可复现。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 并发单元 | 平台线程 / 虚拟线程 | goroutine | 事件循环 + worker | `threading` / asyncio | `std::thread` / tokio task |
| 共享计数 | `AtomicInteger`/`synchronized` | `sync.Mutex`/`atomic` | 单线程无需锁 | `threading.Lock`（GIL） | `Arc<Mutex<T>>`/`AtomicUsize` |
| 等待多任务 | `CountDownLatch`/`join` | `sync.WaitGroup` | `Promise.all` | `ThreadPoolExecutor`/`join` | `join()` / `Barrier` |
| 线程池 | `ExecutorService` | 无（goroutine 自带调度） | `worker_threads` 池 | `concurrent.futures` | rayon / tokio |
| 内存可见性 | `volatile` / JMM happens-before | Go memory model | 单线程 | GIL + 锁 | `Atomic*` / `Ordering` |
