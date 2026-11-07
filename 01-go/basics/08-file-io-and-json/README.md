# 08 · 文件 I/O 与 JSON

## 这个例子讲什么

后端天天在做的事：把 struct 变成 JSON、把请求体读回 struct、大文件流式处理、临时文件清理。

## 怎么跑

```
cd backend-roadmap/01-go && go run ./basics/08-file-io-and-json
```

## 关键点

- **`json.Marshal` / `Unmarshal` 靠 struct tag 控制字段名**：`omitempty` 省略零值、`-` 完全忽略。字段名不加 tag 时用导出字段名。
- 未导出的字段（小写）**不会**被序列化。
- `json.Decoder` 从 `io.Reader` 流式读，处理大请求体/NDJSON 时不用把整个内容读进内存；对应 `json.Encoder` 直接写 `io.Writer`。
- `bufio.Scanner` 默认单行上限 64KB，处理超长行要调 `scanner.Buffer`。
- `os.MkdirTemp` + `defer os.RemoveAll` 是示例和测试里安全用临时文件的标准写法。
- 文件句柄记得 `Close`，并检查错误；生产代码里用 `defer` 紧跟在 `Open` 成功后。

## 和其他语言的差异

| 语言 | JSON 序列化 | 流式 | 注意点 |
|---|---|---|---|
| Go | `encoding/json` + struct tag | `Decoder` / `Encoder` | 字段必须导出 |
| Node.js | `JSON.stringify` 原生 | Stream + pipeline | 无 schema，校验需 Zod |
| Python | `json` / pydantic | 迭代器 | 动态类型，字段名靠 dict key |
| Java | Jackson / Gson | Streaming API | 需注解或 POJO |
| Rust | serde + derive | `serde_json::Deserializer` | 编译期生成，零反射 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
