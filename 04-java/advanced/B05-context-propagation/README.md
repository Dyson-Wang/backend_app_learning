# B05 上下文取消与跨服务传播

## 问题复现

上游给 50ms 预算，下游要 200ms。中间层不传截止时间，下游老老实实跑满 200ms——上游早断了，做了无用功还占着线程。

## 错误写法

```java
return downstream(null, workMs);      // 忘了传 deadline
```

## 正确写法

截止时间当参数传，不是全局变量：

```java
static String downstream(Deadline d, long workMs) {
    if (d.remainingMs() < workMs) return "deadline exceeded";
    ...
}
```

异步链上用 `orTimeout`，跨服务用请求头：

```text
X-Request-Deadline: 1730000000123
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：第一种约 50ms 返回 `deadline exceeded`；第二种 200ms 才返回，证明取消断在中间层。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 载体 | `context.Context` | `AbortSignal` | 显式参数 | `orTimeout` / `Deadline` | `tokio::time::timeout` |
| 强制程度 | 约定 + lint | 约定 | 约定 | 约定 | 约定 |
| 跨服务 | header | header | header | header | header |
