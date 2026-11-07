# 12 · Rust · 字符串与字符编码

## 这个例子讲什么

`String`（拥有所有权的堆字符串）与 `&str`（借用切片）的区别、UTF-8 变长编码下「字节数 ≠ 字符数」、为什么 Rust 禁止 `s[0]` 下标、`get(range)` 的安全切片语义、`push`/`format!` 拼接、原始字符串 `r"..."` / `r#"..."#`、查找切分与 `parse`。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/12-strings-and-encoding
cargo run --quiet
```

## 关键点

- `String` 本质是 `Vec<u8>` 的 UTF-8 包装，`len()` 返回**字节数**；`"世界"` 的 `len()` 是 6，`chars().count()` 是 2。
- 索引 `s[0]` 编译不过：UTF-8 是变长的，按字节下标可能落在字符中间。只能用 `get(..)` 拿 `Option<&str>`，或 `chars()` 迭代。
- `get(0..2)` 对 `"héllo"` 返回 `None`，因为 `é` 占 2 字节、边界在 1 和 3，2 是非法切点。这是 Rust 用类型系统强制你处理编码边界。
- `&str` 是借用视图，`String` 才能增长；函数参数优先用 `&str`（可同时接受两者）。
- `r"..."` 不解析转义，`r#"..."#` 还能包含双引号，写正则和 Windows 路径很实用。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 字符串类型 | `string` 不可变 UTF-8 字节 | `string` 不可变 UTF-16 | `str` 不可变 Unicode 码点 | `String` 不可变 UTF-16 | `String` 可变 / `&str` 借用 |
| 长度语义 | `len()` 字节 | `.length` UTF-16 码元 | `len()` 码点 | `.length()` UTF-16 码元 | `len()` 字节 / `chars().count()` 字符 |
| 下标访问 | `s[i]` 字节 | `s[i]` 码元 | `s[i]` 码点 | `s.charAt(i)` 码元 | 禁止，只能 `get(range)` |
| 可变拼接 | `+=` 生成新串 | `+=` 生成新串 | 不可变，`join` | `StringBuilder` | `push_str` 原地增长 |

Go 的 `len()` 也是字节、下标也是字节，但允许越界切片 panic；Python 下标是码点最直观；Java/JS 用 UTF-16，遇到 emoji（代理对）长度会算成 2；Rust 直接把「按字节下标」这条路封死。
