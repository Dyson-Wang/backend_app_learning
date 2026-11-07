# B02 幂等键与重复请求

## 问题复现

用户连点两次提交，或者网关超时自动重试，同一个请求到达两次 → 两张订单、两次扣款。

## 错误写法

```python
if key in done:
    return done[key]
order = create(key)      # 查和写之间有窗口
```

20 个线程全部穿过判断，建出 20 张单。

## 正确写法

「占位」和「判断」放在同一把锁里：

```python
with self._mu:
    if key in self._done: ...
    ev = self._inflight.get(key)
    owner = ev is None
    if owner:
        ev = threading.Event()
        self._inflight[key] = ev
```

落库版靠唯一索引 + `ON CONFLICT`，并**把响应体一起存**，重放时原样返回。

## 怎么验证

```bash
python3 main.py
```

预期：`真实建单 = 1，回放 = 19，订单 id 去重后 = 1`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 占位原语 | `singleflight` | Promise Map | `threading.Event` | `putIfAbsent` | `Mutex<HashMap>` |
| 落库保证 | 唯一索引 | 唯一索引 | 唯一索引 | 唯一索引 | 唯一索引 |
| 常见坑 | 忘存响应体 | 忘存响应体 | 忘存响应体 | 事务里吞异常 | 忘存响应体 |
