# A05 · Go · 背压与并发上限

## 问题复现

入口不限流、队列无界时，突发流量会被全部接收，内存持续上涨，GC 压力变大，最终 OOM 或延迟雪崩。本例用容量 5 的 channel 提交 20 个任务，把「接收 5 个、拒绝 15 个」直接打出来。

## 错误写法

```go
queue := make([]int, 0) // 无界
for _, task := range tasks {
    queue = append(queue, task) // 来多少收多少
}
```

没有上限就没有背压：上游永远感受不到下游变慢，问题只会以 OOM 的形式在更晚、更难排查的时候爆发。

## 正确写法

用带缓冲的 channel + `select ... default` 做「非阻塞入队」：

```go
select {
case q <- task:
    accepted++
default:
    rejected++ // 满则快速拒绝
}
```

- 容量就是并发上限，内存占用可预测；
- 拒绝发生在入口，响应时间是常数，不会拖垮调用方；
- 生产上把 `rejected` 映射成 `429 Too Many Requests` 或 `503`，并打点监控。

本例用 `gate` 先挡住消费者，保证提交阶段队列必满，输出完全确定。

## 怎么验证

```bash
cd backend-roadmap/01-go/advanced/A05-backpressure
go run .
```

预期：`accepted=5 rejected=15 processed=5`，最后一行 `无界队列长度=20`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 有界队列 | `chan` + `cap` | 数组 + 长度判断 | `queue.Queue(maxsize)` | `ArrayBlockingQueue` | `sync_channel` |
| 满时行为 | `select default` | 手动 `if` | `put_nowait` 抛 `Full` | `offer` 返回 `false` | `try_send` 返回 `Full` |
| 阻塞式入队 | `q <- x` | 无（单线程） | `put()` 阻塞 | `put()` 阻塞 | `send()` 阻塞 |

Go 的 channel 把「容量 + 阻塞/非阻塞」做成了语言原语；Java 的 `BlockingQueue` 有 `put` / `offer` / `add` 三种语义；Python 的 `Queue` 满时 `put_nowait` 抛异常；Rust 的 `sync_channel` 有 `send` 和 `try_send`；Node 单线程没有内建阻塞队列，要自己判长度。
