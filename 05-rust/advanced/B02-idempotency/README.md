# B02 幂等键与重复请求

## 问题复现

用户连点两次，或网关超时重试，同一请求到达两次 → 两张订单、两次扣款。

## 错误写法

```rust
if let Some(o) = store.get(key) { return o; }
let o = create(key);      // 查和写之间有窗口
```

## 正确写法

「占位」和「判断」在同一把锁里，后到的线程 `Condvar::wait`：

```rust
while g.busy.contains(key) {
    g = self.cv.wait(g).unwrap();
    if let Some(o) = g.done.get(key) { return (o.clone(), true); }
}
g.busy.insert(key.to_string());
drop(g);
```

落库版靠唯一索引 + `ON CONFLICT DO NOTHING`，并把响应体一起存。

## 怎么验证

```bash
cargo run --quiet
```

预期：`真实建单 = 1，回放 = 19，订单 id 去重后 = 1`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 占位原语 | `singleflight` | Promise Map | `threading.Event` | `putIfAbsent` | `Mutex` + `Condvar` |
| 锁粒度 | 整表 | 单线程 | 整表 | 整表 | 整表 |
| 常见坑 | 忘存响应体 | 忘存响应体 | 忘存响应体 | 事务吞异常 | `drop(g)` 漏写导致死锁 |
