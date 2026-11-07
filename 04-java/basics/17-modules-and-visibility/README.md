# 17 · Java · 模块进阶：可见性与命名空间

## 这个例子讲什么

四级访问修饰符、命名冲突（无 import 别名）、目录即包、循环依赖编译单元内可解析、内部可见性边界，以及 JPMS `exports` 的定位。

## 怎么跑

```
cd backend-roadmap/04-java/basics/17-modules-and-visibility
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **可见性四级**：`public`（任意）> `protected`（同包 + 子类）> 包私有（同包）> `private`（本类）。本系列单文件都在**默认包**，包私有 ≈ 同文件可见。
- **命名冲突无别名**：`import` 只能省略包名，遇到 `java.util.Date` / `java.sql.Date` 同名时，Java 没有 `import as`，只能写**全限定名**——这是与 Python / Rust 别名机制的明显差异。
- **目录即包**：`package a.b;` 必须对应 `a/b/` 目录；全限定类名 = 包 + 类名，是命名空间的根本。
- **循环依赖**：`A` 持有 `B`、`B` 持有 `A`，同一编译单元内 javac 能正常解析（多文件用 `javac *.java` 也能处理）；跨模块循环依赖才会被 JPMS 拒绝。
- **JPMS**（Java 9）：`module-info.java` 用 `requires` / `exports` / `opens` 做真正的模块导出与强封装，`exports` 决定哪些包对外可见；单文件示例处于**未命名模块**，用访问修饰符演示边界。
- 注意：Java 11 起「嵌套类」互为 nestmate，同一顶层类内部可互访 `private`；真正的 private 边界要跨文件 / 跨包才体现。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 导出规则 | `public` / JPMS `exports` | 首字母大写即导出 | ESM `export` | 约定（`_`/`__all__`） | `pub` / `pub(crate)` |
| 命名空间 | `package` + 全限定名 | import path（URL） | 模块 + 包名 | 包 + 模块 | `crate::mod` |
| 冲突/别名 | 无别名，用全限定名 | `import alias "path"` | `import {x as y}` | `import x as y` | `use ... as ...` |
| 目录即包 | 是（强制对应） | 目录即 import path | 目录即模块路径 | 目录 + `__init__.py` | `mod` + 文件/目录 |
| 循环依赖 | 同编译单元可解析 | 包级禁止循环导入 | 可（但易出 undefined） | 可（谨慎） | 模块树禁止循环 |
