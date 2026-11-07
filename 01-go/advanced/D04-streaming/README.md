# D04 流式输出（SSE / chunked）

## 问题复现

接口要返回 5 块数据，每块算 20ms。全部算完再返回，用户盯着白屏 100ms；如果每块是 200ms，就是 1 秒白屏。

## 错误写法

```go
var buf bytes.Buffer
for _, item := range items { buf.Write(render(item)) }   // 攒够再发
w.Write(buf.Bytes())
```

## 正确写法

```go
w.Header().Set("Content-Type", "text/event-stream")
flusher := w.(http.Flusher)
for _, item := range items {
    fmt.Fprintf(w, "data: %s\n\n", render(item))
    flusher.Flush()          // 关键：每块立刻推出去
}
```

Nginx 要关缓冲：`proxy_buffering off;`，否则又会攒起来。

## 怎么验证

```bash
go run .
```

预期：缓冲模式 `首字节 = 100ms`；流式模式 `首字节 = 20ms`，两者总耗时都是 100ms。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 原语 | `http.Flusher` | `res.write()` | `wfile.flush()` | `OutputStream.flush()` | `write_all` |
| SSE 头 | `text/event-stream` | 同 | 同 | 同 | 同 |
| 背压 | `select` + ctx | `write()` 返回 false | 阻塞写 | 阻塞写 | 阻塞写 |
