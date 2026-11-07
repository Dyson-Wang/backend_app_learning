# 17 · 模块进阶：可见性与命名空间

## 这个例子讲什么

Go 的可见性规则、包别名、`reflect` 看到的包/字段信息。用单文件可观察的方式演示「首字母大小写决定导出」，并在关键点里把只能在多包工程里体现的「目录即包 / 循环依赖 / internal 边界」讲清楚。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/17-modules-and-visibility
```

## 关键点

- **导出规则极简**：标识符首字母大写即导出（包外可见），小写即私有。**没有 `public`/`private`/`export` 关键字**，规则作用于包级别而非类型。
- **可观察证据**：`encoding/json` 只序列化导出字段，所以 `cfg` 里 `secret`/`version` 不出现在 JSON 中；`reflect.StructField.PkgPath` 非空即未导出。同一包内私有成员仍可互相访问（本例 `buildTag()`）。
- **命名冲突与别名**：两个包都叫 `rand` 时写 `mrand "math/rand"`、`crand "crypto/rand"`；别名还有 `_`（仅副作用的空导入）和 `.`（点导入，少用）。
- **目录即包**：一个目录 = 一个包，包名是导入路径最后一段（如 `net/url` 的包名是 `url`），完整路径由 `go.mod` 的 module 加目录决定。本例 `reflect.Type.PkgPath()` 打印 `net/url` 作证。
- **循环依赖**：Go **在编译期禁止 import cycle**，包之间必须是 DAG；出现环直接编译失败，只能抽公共包或改依赖方向。
- **internal 边界**：放在 `<parent>/internal/...` 下的包，只能被 `<parent>` 目录树内的代码导入，编译器强制；这是标准库和大型项目隐藏实现细节的惯例。
- 单目录单 `package main` 无法演示多包导入环，故这些边界以说明为主，机制本身是编译期硬约束。

## 和其他语言的差异

| 语言 | 导出/私有 | 命名空间与别名 | 循环依赖 / 内部边界 |
|---|---|---|---|
| Go | 首字母大小写 | 目录即包；`import alias "path"`、空导入、点导入 | 编译期禁止 import cycle；`internal/` 目录边界 |
| Node.js | ESM `export`/`import`；`#private` 字段 | `package.json` 的 `exports`；`import x as y` | 循环导入允许但易得到未初始化值；无 internal 机制 |
| Python | 约定 `_name` 私有（非强制）；`__all__` | 文件即模块、包即目录；`import x as y` | 循环导入常需延迟 import；无强制私有 |
| Java | `public/protected/默认/private` 四级 | `package` 对应目录；`import`；JPMS `module-info.java` | 循环依赖编译可过但设计不推荐；JPMS 有 `exports` 控制 |
| Rust | `pub` / `pub(crate)` / 默认私有 | crate + `mod` 树；`use ... as` | 模块可互相引用；`pub(crate)` 等可见性粒度更细 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
