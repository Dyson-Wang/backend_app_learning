# 09 · Python · 并发原语入门

## 这个例子讲什么

`threading.Thread` + `join` 起线程，`threading.Lock` 保护共享计数（4 线程各加 1000 次得到确定的 4000），`daemon` 线程，`ThreadPoolExecutor` 的 `map` / `submit`，并点明 GIL 对并发模型的影响。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/09-concurrency-primitives/main.py
```

## 关键点

- **`counter["value"] += 1` 不是原子操作**，多线程下会丢更新；用 `with lock:` 包住，结果才可复现。
- `Thread.start()` 后必须 `join()`，否则主线程可能先退出；`join(timeout=...)` 避免无限等待。
- `ThreadPoolExecutor` 用 `with` 管理，`pool.map` 保序、`submit().result()` 从 `Future` 取值。
- **GIL**：CPython 同一时刻只有一个线程执行字节码，线程适合 I/O 密集；CPU 密集要用 `multiprocessing` 或释放 GIL 的 C 扩展。
- 队列、条件变量、扇出扇入等通信原语留到 19 课；异步与取消留到 20 课。

## 和其他语言的差异

- **Go**：`go` 关键字起 goroutine、`sync.Mutex` / `sync.WaitGroup`、`channel` 通信、真多核并行。
- **Node.js/TS**：单线程事件循环 + `Promise`、`worker_threads` 才有多线程、`Atomics` 共享内存。
- **Java**：`Thread` / `ExecutorService`、`synchronized` / `ReentrantLock`、虚拟线程（Java 21）。
- **Rust**：`std::thread`、`Arc<Mutex<T>>`、`mpsc`，`Send`/`Sync` 在**编译期**挡住数据竞争。
- **Python**：`threading.Thread` + `Lock`、`ThreadPoolExecutor`、GIL 限制真并行、无编译期数据竞争检查。
