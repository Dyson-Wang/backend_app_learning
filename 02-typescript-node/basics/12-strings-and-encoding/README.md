# 12 · Node.js · 字符串与字符编码

## 这个例子讲什么

把 JS 字符串的不可变性、UTF-16 码元 / Unicode 码点 / UTF-8 字节三者区别、按索引与 `for...of` 的遍历差异、拼接格式化、转义与 `String.raw`、大小写/查找/切分、以及 `TextEncoder`/`TextDecoder`/`Buffer` 的编码解码一次讲透。核心结论：`'😀'.length === 2`，但它是 1 个码点、4 个字节。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/12-strings-and-encoding/index.mjs
```

只用 `node:util` 的 `TextEncoder`/`TextDecoder`。

## 关键点

- 字符串是**不可变**的原始值：ESM 严格模式下 `s[0] = 'X'` 会抛 `TypeError`，所有「改」方法都返回新串。
- JS 字符串底层是 **UTF-16 码元**：`.length`、索引、`charCodeAt` 数的是码元；`for...of`、`[...s]`、`codePointAt` 按码点工作。
- 增补平面字符（如 emoji）占 2 个码元、1 个码点、UTF-8 下 4 字节；按索引遍历可能切出孤立的代理项（乱码）。
- `TextEncoder` 把串编成 UTF-8 `Uint8Array`，`TextDecoder` 反向；`Buffer.from(str,'utf8')` 与 `Buffer.byteLength` 是 Node 的字节视角。
- `String.raw` 作标签模板保留反斜杠，`String.fromCodePoint(0x1f600)` / `'\u{1f600}'` 写码点，`JSON.stringify` 负责转义。
- `replace` 默认只换第一个，`replaceAll` 或 `g` 正则才全换；`split` 会保留空串。
- `Buffer` 的 `hex`/`base64` 与 `encodeURIComponent` 是不同用途的编码，别混用。
- 排序默认按码元（大写在小写前）；本地化顺序用 `Intl.Collator`（18 课）。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 字符串内部 | UTF-8 字节 | Unicode 码点 | UTF-16 码元 | UTF-8 字节 | UTF-16 码元 |
| 可变性 | 不可变 | 不可变 | 不可变 | `String` 不可变 / `String::push_str` 可变 | 不可变 |
| `len` 含义 | 字节数 | 码点数 | 码元数 | 字节数 | **码元数** |
| 遍历 | `range` 给 rune | 给码点 | `codePoints()` | `.chars()` 给码点 | `for...of` 给码点，索引给码元 |
| 原始字符串 | 反引号 raw string | `r"..."` | 文本块 `"""` | `r#"..."#` | `String.raw\`...\`` |
| 编码类型 | `[]byte` / `rune` | `str` / `bytes` | `String` / `byte[]` | `String` / `Vec<u8>` | `string` / `Buffer` / `Uint8Array` |

Go 和 Rust 的 `len()` 是字节数、Python 是码点数，而 **JS 与 Java 是 UTF-16 码元数**——emoji 上 `.length === 2` 是 JS 最著名的坑。`String.raw` 与 Go 反引号、Python `r""` 对应；Java 直到 15 才有文本块，且没有 raw 前缀。
