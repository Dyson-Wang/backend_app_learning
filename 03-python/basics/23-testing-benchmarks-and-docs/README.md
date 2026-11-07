# 23 · Python · 测试、基准与文档

## 这个例子讲什么

`unittest` 的表驱动 `subTest` + 带信息的断言 + `assertRaises`、`doctest` 把文档示例当断言跑、`timeit` 微型基准、docstring / `inspect.getdoc` / `pydoc` 文档生成；入口脚本直接 `python3 main.py` 即可运行。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/23-testing-benchmarks-and-docs/main.py
```

## 关键点

- **表驱动测试**：`(输入, 期望)` 列表 + `with self.subTest(...)`，某行失败不会中断其余行，且失败信息自带参数。
- 本例用 `unittest.TextTestRunner` 显式跑 `TestSuite`（不调 `unittest.main()`，避免 `sys.exit`），保证 `main.py` 直接运行仍返回退出码 0。
- `doctest` 把 docstring 里的 `>>>` 当可断言示例，`DocTestFinder` + `DocTestRunner` 可在普通脚本里跑并拿到 `TestResults`。
- `timeit.timeit(stmt, number=N)` / `repeat(...)` 做基准；本课**只演示 API 与布尔结论**，不把机器相关的耗时数值当结论。
- 文档：`__doc__` + `inspect.getdoc` + `pydoc.render_doc`；覆盖率思路用标准库 `trace` 模块（`python3 -m trace --count`），本仓库不引第三方。

## 和其他语言的差异

- **Go**：`go test` 表驱动 + `testing.B` 基准 + `go test -cover` + `go doc` / godoc。
- **Node.js/TS**：`node:test` / Jest / Vitest、`benchmark`、`JSDoc` / TypeDoc、`c8` 覆盖率。
- **Java**：JUnit 5、`JMH` 基准、`Javadoc`、`JaCoCo` 覆盖率。
- **Rust**：`#[cfg(test)]` + `cargo test`、`cargo bench`、`///` 文档注释 + `rustdoc`。
- **Python**：`unittest`（`subTest`）、`doctest`、`timeit`、`docstring` + `pydoc`、`trace` 覆盖率思路。
