# B04 · Go · 竞态、死锁与调试

## 问题复现

两个 goroutine 同时执行 `v := counter; counter = v + 1`，两次读都拿到旧值，最终只加了一次。并发越高丢得越多，而且不加 `-race` 根本看不出来。

## 错误写法

```go
v := counter
counter = v + 1 // 读和写之间没有保护
```

`counter++` 不是原子操作，它由读、加、写三步组成。

## 正确写法

- 互斥：`sync.Mutex` 包住整个读改写；
- 原子：`atomic.AddInt64` 一步完成，适合计数器；
- 死锁：两把锁必须**全局统一顺序**获取，或者用 `TryLock` + 超时。

本例用 `WaitGroup` 强制两个 goroutine 都读完再写，让丢更新必然发生；死锁部分用 `select + time.After` 检测，而不是让程序真的挂住。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/B04-race-and-deadlock
go run .
go run -race .   # 有工具链时，第 1 段会报 DATA RACE
```

预期：`期望 2，实际 1`、`期望 4000，实际 4000`（两次）、`检测到死锁`、两行「完成」。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 竞态检测 | `go run -race` | 单线程无 | GIL 不防读改写 | 无内建，靠 JUC | 编译期禁止 |
| 互斥 | `sync.Mutex` | 无锁，靠串行化 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 原子 | `sync/atomic` | 无 | 无（需锁） | `AtomicInteger` | `AtomicUsize` |
| 死锁检测 | `select + timeout` | Promise.race + timeout | `join(timeout)` | `tryLock(timeout)` | `recv_timeout` |

Go 有官方 race detector，是排查并发问题最实用的工具；Node 单线程没有数据竞争但有异步交错；Python 的 GIL 只保证字节码原子，`x += 1` 仍会丢更新；Rust 在编译期就拒绝数据竞争，但死锁仍然可能。
