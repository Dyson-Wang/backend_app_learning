# A03 · Node.js · 熔断与降级

## 问题复现

下游超时后，Node 的事件循环仍会不断发起新的 `fetch`，每个失败请求都占着 socket 和 Promise，最终 `MaxListenersExceededWarning`、内存上涨、事件循环被回调塞满。熔断器让「下游已挂」这件事变成一个**本地状态判断**，不再产生网络调用。

## 错误写法

```js
// 没有状态机，失败就重试
for (let i = 0; i < 3; i += 1) {
  try { return await fetch(url); } catch { /* 继续重试 */ }
}
```

Node 是单线程，但 `await` 会让出执行权，多个请求仍会交错进入这段代码，重试叠加后下游压力是乘法级增长。

## 正确写法

`Breaker` 三态 + `probeInFlight` 标记：

- `open` 时 `allow()` 直接返回拒绝，不创建 Promise、不发请求；
- 冷却结束转 `half-open`，用 `probeInFlight` 保证**只有一个**探测请求；
- 探测成功 `success()` 清零失败计数并回 `closed`，失败 `failure()` 立刻回 `open`。

Node 没有锁，但 `allow()` 到 `success()/failure()` 之间不能有 `await`，否则状态判断会失去原子性——这是 JS 版熔断器最容易写错的地方。

## 怎么验证

```bash
cd backend-roadmap/02-typescript-node/advanced/A03-circuit-breaker
node index.mjs
```

预期看到 `state=closed` → 5 次失败后 `state=open` → `rejected: circuit open: fast fail` → 时钟推进 200 后 `state=half-open` → 探测成功回 `closed`。用注入的 `now` 推进时间，输出稳定。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 并发保护 | `sync.Mutex` | 单线程 + 无 `await` 临界区 | `threading.Lock` | `synchronized` | `Mutex<T>` |
| 时间注入 | 函数字段 | `now` 回调 | `clock` 可调用 | `Clock` | `fn() -> Instant` |
| 半开探测 | 状态位 | `probeInFlight` | 布尔标记 | `AtomicBoolean` | `Mutex` 内状态 |

Go 靠 mutex 保证原子性，Node 靠「临界区里不 await」，Java 靠锁或原子变量，Rust 靠所有权让状态只能通过 `Mutex` 访问。
