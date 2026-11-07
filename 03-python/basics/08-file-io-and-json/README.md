# 08 · Python · 文件 I/O 与 JSON

## 这个例子讲什么

用 `pathlib.Path` + `with` 做文本/二进制读写，`open('a')` 追加，UTF-8 编码下字符数与字节数的差异，`json.dumps/loads/dump/load` 与 `ensure_ascii`，以及 `JSONDecodeError` 处理；全部写在 `TemporaryDirectory` 里，退出即清理。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/08-file-io-and-json/main.py
```

## 关键点

- **`with open(...)` 自动关闭文件**；务必显式 `encoding="utf-8"`，否则依赖平台默认编码。
- **字符数 ≠ 字节数**：`len("中文")==2`，但 `.encode("utf-8")` 是 6 字节；`Content-Length` 要按字节算（10 课 HTTP 服务踩过这个坑）。
- `json.dumps(..., ensure_ascii=False, sort_keys=True)` 输出中文且键有序，便于对比；`loads` 还原为 `dict`。
- 大对象用 `json.dump/load` 直接接文件对象，省掉中间字符串。
- 非法输入抛 `json.JSONDecodeError`（`ValueError` 子类），用 `try/except` 兜住，别让服务直接 500。

## 和其他语言的差异

- **Go**：`os` / `bufio` / `encoding/json`，靠 struct tag 映射字段，`defer f.Close()` 关闭。
- **Node.js/TS**：`node:fs` / `fs/promises`、`JSON.parse/stringify`、`Buffer` 处理二进制。
- **Java**：`java.nio.file.Files`、Jackson / Gson、`try-with-resources`、`StandardCharsets.UTF_8`。
- **Rust**：`std::fs`、`serde_json`（第三方）、`?` 传播 `io::Result`。
- **Python**：`pathlib.Path`、`open(encoding=...)` + `with`、`json` 模块、`ensure_ascii`、`JSONDecodeError`。
