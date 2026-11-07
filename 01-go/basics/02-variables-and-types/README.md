# 02 · 变量、类型、常量、零值

## 这个例子讲什么

Go 的类型系统底线：每个变量都有零值、没有隐式类型转换、`string` 是字节序列、指针不能做算术。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/02-variables-and-types
```

## 关键点

- **零值**：数值 0、字符串 `""`、布尔 false、指针/切片/map/通道/函数/接口 nil。没有「未初始化变量」这种状态，这消灭了一大类 bug。
- **`:=`** 只能在函数内用，且左边至少有一个新变量；包级变量必须用 `var`。
- **没有隐式转换**：`int` 和 `int64` 不能直接相加，必须显式转换。自定义类型 `UserID` 与 `int64` 也是不同类型。
- `len(string)` 返回**字节数**，中文要 `len([]rune(s))`；`range` 字符串时拿到的是 rune 和字节下标。
- `iota` 从 0 开始按行自增，常用于枚举。
- 指针只用于「共享 / 修改」，不能做 `p++` 这类指针运算。

## 和其他语言的差异

| 语言 | 空值 | 类型转换 | 字符串 |
|---|---|---|---|
| Go | 零值，没有 null | 必须显式 | UTF-8 字节序列，`len` 是字节数 |
| Node.js | `null` / `undefined` | 隐式（宽松） | UTF-16 |
| Python | `None` | 动态类型，无转换 | Unicode 码点，`len` 是字符数 |
| Java | `null`（基本类型有默认值） | 窄化需显式 | UTF-16 |
| Rust | 没有 null，用 `Option` | 必须显式 `as` / `From` | UTF-8 字节序列 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
