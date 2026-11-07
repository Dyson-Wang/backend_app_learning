# D04 流式输出（SSE / chunked）

## 问题复现

5 块数据每块 20ms，全部算完再返回，用户白屏 100ms；每块 200ms 就是 1 秒。

## 错误写法

```js
const parts = [];
for (const item of items) parts.push(render(item));
res.end(parts.join(''));      // 攒够再发
```

## 正确写法

```js
res.writeHead(200, { 'Content-Type': 'text/event-stream', 'Cache-Control': 'no-cache' });
for (const item of items) {
  if (!res.write(`data: ${render(item)}\n\n`)) await once(res, 'drain');   // 背压
}
res.end();
```

## 怎么验证

```bash
node index.mjs
```

预期：缓冲 `首字节 = 100ms`；流式 `首字节 = 20ms`，总耗时都是 100ms。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 原语 | `http.Flusher` | `res.write()` | `wfile.flush()` | `OutputStream.flush()` | `write_all` |
| 背压 | ctx + channel | `write()` 返回 false + `drain` | 阻塞写 | 阻塞写 | 阻塞写 |
| 陷阱 | 忘 Flush | 忘 `drain` | 忘 `flush` | 忘 `flush` | 忘 `flush` |
