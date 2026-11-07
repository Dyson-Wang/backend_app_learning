# D04 流式输出（SSE / chunked）

## 问题复现

5 块数据每块 20ms，全部算完再返回，用户白屏 100ms。

## 错误写法

```rust
let mut buf = Vec::new();
for item in &items { buf.extend(render(item)); }
stream.write_all(&buf)?;      // 攒够再发
```

## 正确写法

```rust
stream.write_all(b"HTTP/1.1 200 OK\r\nContent-Type: text/event-stream\r\n\r\n")?;
for item in &items {
    stream.write_all(format!("data: {}\n\n", render(item)).as_bytes())?;
    stream.flush()?;          // 关键
}
```

## 怎么验证

```bash
cargo run --quiet
```

预期：缓冲 `首字节 = 100ms`；流式 `首字节 = 20ms`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 原语 | `http.Flusher` | `res.write()` | `wfile.flush()` | `OutputStream.flush()` | `Write::flush` |
| 背压 | ctx + channel | `drain` | 阻塞写 | 阻塞写 | 阻塞写 |
| 陷阱 | 忘 Flush | 忘 drain | 忘 flush | 忘 flush | 忘 flush |
