# B05 上下文取消与跨服务传播

## 问题复现

上游给了 50ms 预算，下游接口要 200ms。如果中间层不把截止时间传下去，下游会老老实实跑满 200ms，上游早就断了——**做了无用功，还占着连接和线程**。

## 错误写法

```python
async def handler(work_ms):
    return await downstream(None, work_ms)   # 忘了传 deadline
```

耗时 200ms，取消完全失效。

## 正确写法

截止时间是一个**显式参数**，不是全局变量：

```python
async def downstream(deadline, work_ms):
    if deadline.remaining() < work_ms / 1000:
        return "deadline exceeded"
    await asyncio.sleep(work_ms / 1000)
```

跨服务时把剩余预算放进请求头：

```text
X-Request-Deadline: 1730000000123
```

下游读到后换算成本地 deadline，继续往下传。

## 怎么验证

```bash
python3 main.py
```

预期：第一种 52ms 左右返回 `deadline exceeded`；第二种 200ms 才返回，证明取消断在中间层。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 载体 | `context.Context` 首参 | `AbortSignal` | 显式参数 / `asyncio.timeout` | `CompletableFuture.orTimeout` | `tokio::time::timeout` |
| 强制程度 | 约定 + lint | 约定 | 约定 | 约定 | 约定 |
| 跨服务 | 中间件解 header | 中间件解 header | 中间件解 header | 拦截器 | tower layer |
