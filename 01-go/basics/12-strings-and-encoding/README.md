# 12 · 字符串与字符编码

## 这个例子讲什么

字符串的不可变性、UTF-8 字节与 Unicode 码点的区别、两种遍历方式、拼接与格式化、转义与原始字符串、大小写/查找/切分，以及编码解码（rune↔UTF-8、base64、hex）。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/12-strings-and-encoding
```

## 关键点

- **不可变**：`s[i] = x` 编译失败；想要可变就先 `[]byte(s)` / `[]rune(s)` 转成切片，改完再转回来。
- **`len(s)` 是字节数，不是字符数**：`"Go语言"` 是 8 字节、4 个码点。用 `utf8.RuneCountInString` 数字符。
- **遍历差异**：`for i := 0; i < len(s); i++` 拿到字节；`for idx, r := range s` 拿到 rune，`idx` 是该 rune 的**字节偏移**（中文会跳 3）。
- **string 是 UTF-8 字节序列**，要按字符处理请用 `for range` 或 `[]rune(s)`；`unicode/utf8` 提供 `EncodeRune` / `DecodeRune` / `RuneLen`。
- **拼接**：少量用 `+`，带格式用 `fmt.Sprintf`，循环里用 `strings.Builder`（避免反复分配）。
- **转义 vs 原始**：双引号字符串解释 `\n \t \u4e2d`；反引号原始字符串不解释任何转义，适合正则、多行文本。
- **base64/hex 是字节级编码**，与「字符编码」是两件事，别混淆。

## 和其他语言的差异

| 语言 | 内部表示 | 长度语义 | 可变性 / 编码 |
|---|---|---|---|
| Go | UTF-8 字节序列 | `len` = 字节，需 `utf8.RuneCountInString` 数码点 | 不可变；`[]rune` 转码点；`encoding/base64`、`hex` |
| Node.js | UTF-16 code unit | `.length` 是 code unit 数，emoji 会算 2 | 不可变；`for...of` 按码点；`Buffer` 管字节；`TextEncoder` |
| Python | Unicode 码点（PEP 393 紧凑存储） | `len` = 码点数 | 不可变 `str`；`encode`/`decode` 转 `bytes`；`codecs` |
| Java | UTF-16 `char` | `.length()` 是 code unit；`codePointCount` 才准 | 不可变 `String`；`StringBuilder` 可变；`Charset` |
| Rust | UTF-8（`String`/`&str`） | `.len()` = 字节；`.chars().count()` 数码点 | `String` 需 `mut` 才可变；`.bytes()` / `.chars()`；索引不允许 `s[i]` |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
