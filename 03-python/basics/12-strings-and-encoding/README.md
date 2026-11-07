# 12 · Python · 字符串与字符编码

## 这个例子讲什么

字符串不可变性、Unicode 码点 vs UTF-8 字节、按码点/按字节遍历差异、`join`/f-string/`format`/`%` 拼接格式化、转义与原始字符串、大小写/查找/切分、`encode/decode` 与错误策略。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/12-strings-and-encoding/main.py
```

## 关键点

- `str` **不可变**：`s[0]='z'` 抛 `TypeError`，任何「修改」都产生新对象；循环里 `s += x` 是 O(n²)，应改用 `join`。
- `len(str)` 数的是**码点**，`len(str.encode('utf-8'))` 数的是**字节**；中文 3 字节、emoji 4 字节。
- `ord/chr` 在码点与整数间转换；`for c in s` 按码点、`for b in s.encode()` 按字节。
- `r""` 原始字符串不处理转义；`casefold()` 比 `lower()` 更适合大小写无关比较（如德语 ß）。
- `decode(..., errors='replace')` 决定坏字节如何兜底；`codecs.BOM_UTF8` 处理 BOM。

## 和其他语言的差异

- **Go**：`string` 是只读**字节**序列、`rune`=int32、`for range` 按 rune 解码、`unicode/utf8` 包。
- **Node.js/TS**：UTF-16 码元、`.length` 按码元（emoji 算 2）、`for...of` 按码点、`TextEncoder`/`Buffer`。
- **Java**：UTF-16 `char`、`String` 不可变、`codePointAt`、`StandardCharsets`。
- **Rust**：`String` 是**可变** UTF-8、`&str`、`chars()` 按码点、`len()` 按字节，按字节索引在非边界会 panic。
- **Python**：`str` 不可变 Unicode 码点序列、`encode/decode`、`ord/chr`、`r""`、f-string。
