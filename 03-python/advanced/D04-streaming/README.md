# D04 流式输出（SSE / chunked）

## 问题复现

5 块数据每块 20ms，全部算完再返回，用户白屏 100ms。

## 错误写法

```python
body = b"".join(render(i) for i in items)
self.wfile.write(body)      # 攒够再发
```

## 正确写法

```python
self.send_header("Content-Type", "text/event-stream")
for item in items:
    self.wfile.write(f"data: {render(item)}\n\n".encode())
    self.wfile.flush()      # 关键
```

## 怎么验证

```bash
python3 main.py
```

预期：缓冲 `首字节 = 100ms`；流式 `首字节 = 20ms`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 原语 | `http.Flusher` | `res.write()` | `wfile.flush()` | `OutputStream.flush()` | `write_all` |
| 背压 | ctx + channel | `drain` 事件 | 阻塞写 | 阻塞写 | 阻塞写 |
| 陷阱 | 忘 Flush | 忘 drain | 忘 flush | 忘 flush | 忘 flush |
