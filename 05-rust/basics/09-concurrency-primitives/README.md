# 09 · Rust · 并发原语入门

## 这个例子讲什么

三块最小并发原语：`thread::spawn` + `join`、`Arc<Mutex<T>>` 共享可变状态、`std::sync::mpsc` 通道传消息。重点体会 `move` 闭包移交所有权，以及 `Mutex` 锁作用域要手动收窄。更深的 channel / 扇入扇出放到 19 课。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/09-concurrency-primitives
cargo run
```

## 关键点

- `thread::spawn` 要求闭包 `'static`，所以要 `move` 把用到的值搬进去；`join()` 返回 `Result<T, _>`，子线程 panic 会变成 `Err`。
- `Arc<T>`（多所有权）+ `Mutex<T>`（互斥）是 Rust 共享可变状态的标准答案；**编译器强制你显式表达共享**，Go 里随手一个包级变量在这里会被借用检查器拦下。
- `c.lock().unwrap()` 返回 `MutexGuard`，离开作用域自动解锁；本例用内层 `{}` 把临界区收窄，就是提醒别持锁干别的事。
- 锁被 poison（持锁线程 panic）时 `lock()` 返回 `Err`，这里用 `unwrap` 仅作教学。
- `mpsc::channel` 可 `tx.clone()` 多生产者；`drop(tx)` 后 `rx.iter()` 才会在所有发送者消失时结束，否则会一直等。
- 为输出稳定，收到的消息先 `collect` 再 `sort` 后打印；线程调度顺序本是不确定的。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 并发单元 | goroutine | 事件循环 / worker | `threading`（受 GIL 限） | `Thread` / 虚拟线程 | OS 线程 / 异步任务 |
| 启动 | `go f()` | `Worker` / 异步 | `Thread.start()` | `new Thread().start()` | `thread::spawn` |
| 共享状态 | `sync.Mutex` | 单线程为主 / `SharedArrayBuffer` | `threading.Lock` | `synchronized` / `Lock` | `Arc<Mutex<T>>` |
| 消息传递 | `chan` | `postMessage` | `queue.Queue` | `BlockingQueue` | `std::sync::mpsc` |
| 数据竞争 | 运行时 race detector | 单线程少有 | GIL 保护字节码 | 靠 synchronized | **编译期**借用检查杜绝 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
