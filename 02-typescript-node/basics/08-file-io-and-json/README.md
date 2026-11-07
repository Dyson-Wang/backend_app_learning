# 08 · Node.js · 文件 I/O 与 JSON

## 这个例子讲什么

用 `node:fs` 同步 / promises / 流三种读法，`node:path`、`node:os`、`node:url` 定位文件，演示 `JSON.parse` 的 reviver、`JSON.stringify` 的 replacer 与缩进、JSON 的能力边界（BigInt/undefined），以及 `Buffer` 与 UTF-8 字节/字符的关系。全部在系统临时目录操作并自动清理。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/08-file-io-and-json/index.mjs
```

## 关键点

- `readFileSync`/`writeFileSync` 会阻塞事件循环；服务里用 `fs/promises` 或 `createReadStream` + `for await`。
- `for await (const chunk of stream)` 按块消费，避免把大文件一次性读进内存；块是 `Buffer`，要注意多字节字符跨块问题。
- `JSON.parse(text, reviver)` 可在解析时改写值；`JSON.stringify(obj, replacer, indent)` 可过滤字段并美化。
- JSON 序列化不了 `bigint`、`undefined`（被忽略）、`Symbol`、循环引用；跨层传值用 `structuredClone` 或自己转字符串。
- `Buffer.byteLength('中文')` 是 6 字节，而 `'中文'.length` 是 2 个码元——字节、码元、码点是三件事（12 课细讲）。
- `import.meta.url` 配 `fileURLToPath` 得到模块自身路径，比 `__dirname` 更通用（ESM 里没有 `__dirname`）。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 同步读 | `os.ReadFile` | `open().read()` | `Files.readString` | `fs::read_to_string` | `readFileSync` |
| 异步读 | goroutine + `os.Open` | `aiofiles`（三方） | NIO / `CompletableFuture` | tokio（三方） | `fs/promises` |
| JSON 库 | `encoding/json` + struct tag | `json` | Jackson/Gson（三方） | serde（三方） | 全局 `JSON` 对象 |
| 字节缓冲 | `[]byte` | `bytes` | `byte[]`/`ByteBuffer` | `Vec<u8>` | `Buffer` |
| 流 | `io.Reader` | 文件迭代器 | `InputStream` | `Read` trait | `Readable` + AsyncIterator |

Node 与 Python 把 JSON 塞进标准库的全局对象/内置模块；Go 靠 struct tag 做字段映射；Java/Rust 的真实项目几乎都必须引入 Jackson / serde 这类第三方库。Node 的 `Buffer` 是独立于字符串的字节类型，和 Go 的 `[]byte` 最像。
