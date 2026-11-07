# 12 · Java · 字符串与字符编码

## 这个例子讲什么

`String` 不可变性、拼接与格式化、大小写 / 查找 / 切分、码点 vs UTF-16 码元 vs UTF-8 字节、按码点遍历、编码解码、转义与 text block、`strip`/`repeat`。

## 怎么跑

```
cd backend-roadmap/04-java/basics/12-strings-and-encoding
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **`String` 不可变、线程安全**：`concat`/`replace` 都返回新对象，循环拼接用 `StringBuilder`。
- Java 字符串是 **UTF-16 码元序列**：`length()` 是 `char` 数，`😀` 占 2；码点数要用 `codePointCount` / `codePoints()`；UTF-8 里它占 **4 字节**。这是遍历非 BMP 字符最容易踩的坑。
- **编码 / 解码必须显式指定 `StandardCharsets.UTF_8`**，别用平台默认编码；`new String(bytes, charset)` 往返。
- **text block**（`"""`，Java 15+）是 Java 的「原始字符串」近似物，自动处理缩进；普通字符串仍用 `\n`、`\t`、`\"` 转义。
- `String.format` / `formatted`、`String.join`、`split`、`strip`（Unicode 感知，优于 `trim`）。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 字符串可变性 | 不可变 | 不可变（[]byte 可变） | 不可变 | 不可变 | `String` 可变需 `String::push_str` |
| 内部编码 | UTF-16 | UTF-8 字节 | UTF-16 | Unicode 码点 | UTF-8（强制） |
| 码点遍历 | `codePoints()` | `for range`（rune） | `for...of`（码点） | 直接按码点 | `.chars()` |
| 原始字符串 | text block `"""` | 反引号 raw string | 模板字符串 | `r"..."` | `r#"..."#` |
| 格式化 | `String.format`/`%` | `fmt.Sprintf` | 模板字符串 | f-string / `.format` | `format!` 宏 |
