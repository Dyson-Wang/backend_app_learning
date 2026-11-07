# 08 · Rust · 文件 I/O 与 JSON

## 这个例子讲什么

用 `std::fs` 把一段 JSON 写入临时文件、读回、逐行解析再清理。零第三方依赖，**手写**极简序列化与取字段，目的和 10 课一致：看清 `serde` / `serde_json` 替你做了什么。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/08-file-io-and-json
cargo run
```

## 关键点

- `File::create` 截断/新建，`read_to_string` 一次性读，`BufReader::lines()` 逐行读；三者都要 `std::io::Result`，用 `?` 传播。
- `write_all` / `lines` 等方法来自 `Write` / `BufRead` trait，必须 `use` 到作用域才可见——这是 Rust 常见的新手坑。
- 写文件用 `payload.as_bytes().len()` 计字节；中文 UTF-8 下 `chars().count()` 会算错。
- `split_once` + `find` 手写的 `get_string` 只认最简单的 `"key":"value"`，遇到转义、嵌套、数字就失效——**这就是「为什么需要 serde」**。
- `serde` 靠派生宏在**编译期**生成（零反射开销），对应 Go 的 `encoding/json`、Java 的 Jackson、Python 的 `json`。
- 本机未装 cargo，示例**未经编译验证**。

## 和其他语言的差异

| 能力 | Go | Node.js | Python | Java | Rust(std) |
|---|---|---|---|---|---|
| 读文件 | `os.ReadFile` | `fs.readFile` | `open().read()` | `Files.readString` | `fs::read_to_string` |
| 写文件 | `os.WriteFile` | `fs.writeFile` | `open().write()` | `Files.writeString` | `fs::write` / `File` |
| 逐行 | `bufio.Scanner` | `readline` | 迭代文件对象 | `BufferedReader` | `BufReader::lines` |
| JSON | `encoding/json`（反射） | 内置 `JSON` | 内置 `json` | Jackson / Gson | std 无，生产用 serde |
| 错误 | `(v, err)` | 回调 / Promise | 异常 | 异常 | `io::Result` + `?` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
