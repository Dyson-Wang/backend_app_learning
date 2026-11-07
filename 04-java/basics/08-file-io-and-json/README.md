# 08 · Java · 文件 I/O 与 JSON

## 这个例子讲什么

`java.nio.file.Files` 的 `writeString/readString/readAllBytes`、显式 UTF-8、临时文件与清理、text block、手写 JSON 序列化与解析、字符串转义。

## 怎么跑

```
cd backend-roadmap/04-java/basics/08-file-io-and-json
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- 新 I/O 用 **`java.nio.file` 的 `Path` / `Files`**，配合 `StandardCharsets.UTF_8` 显式指定编码，别依赖平台默认编码。
- `Files.writeString` / `readString` / `readAllBytes` 是 Java 11+ 的便捷方法；`String.lines()` 按行切分流。
- 写临时文件务必 **try/finally 清理**（`deleteIfExists`），示例不往代码库里留垃圾文件。
- Java 标准库**没有内置 JSON 解析器**：这里手写只为教学，生产用 Jackson（Spring Boot 默认）或 Gson + DTO。
- text block `"""..."""`（Java 15+）适合内嵌多行 JSON / SQL；`.strip()` 去掉首尾空白。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 文件读写 | `Files.readString/writeString` | `os.ReadFile/WriteFile` | `node:fs` `readFileSync` | `open()` / `pathlib` | `std::fs::read_to_string` |
| 编码 | 显式 `StandardCharsets.UTF_8` | `[]byte` + 手动 | `Buffer` + 编码参数 | `encoding=` | `String` 本身 UTF-8 |
| JSON | 无内置（用 Jackson） | `encoding/json` | 内置 `JSON` | 内置 `json` | 无内置（用 serde） |
| 资源清理 | try-with-resources/finally | `defer` | finally / `fs` 同步 API | `with` | `Drop` / RAII |
| 多行字符串 | text block | 反引号 raw string | 模板字符串 | 三引号 | raw string `r#"..."#` |
