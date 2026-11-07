# D04 流式输出（SSE / chunked）

## 问题复现

5 块数据每块 20ms，全部算完再返回，用户白屏 100ms。

## 错误写法

```java
ByteArrayOutputStream buf = new ByteArrayOutputStream();
for (Item i : items) buf.write(render(i));
out.write(buf.toByteArray());       // 攒够再发
```

## 正确写法

```java
ex.getResponseHeaders().set("Content-Type", "text/event-stream");
OutputStream os = ex.getResponseBody();
for (Item i : items) {
    os.write(("data: " + render(i) + "\n\n").getBytes());
    os.flush();                     // 关键
}
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：缓冲 `首字节 = 100ms`；流式 `首字节 = 20ms`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 原语 | `http.Flusher` | `res.write()` | `wfile.flush()` | `OutputStream.flush()` | `write_all` |
| 背压 | ctx + channel | `drain` | 阻塞写 | 阻塞写 | 阻塞写 |
| 陷阱 | 忘 Flush | 忘 drain | 忘 flush | 忘 flush | 忘 flush |
