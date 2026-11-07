# 23 · Node.js · 测试、基准与文档

## 这个例子讲什么

用 `node:assert/strict` 搭一个可运行的**表驱动**微型测试 runner（带通过/失败统计与断言失败信息 `actual/expected/operator`），补上 `deepStrictEqual`/`throws`/`ok` 等断言形式；用 `node:perf_hooks` 的 `performance.now()` 做固定迭代的基准演示；再讲清 JSDoc 文档注释、覆盖率思路（`node --test` / `NODE_V8_COVERAGE` / c8）与标准 `node:test` 写法。入口脚本 `node index.mjs` 仍可直接运行并退出码 0。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/23-testing-benchmarks-and-docs/index.mjs
# 也可对照标准写法（本文件未引入）：
# node --test basics/23-.../index.mjs
# node --experimental-test-coverage --test basics/23-.../index.mjs
```

只用 `node:assert/strict` / `node:perf_hooks`，零第三方依赖。

## 关键点

- 表驱动测试：把用例写成 `{name, input, expected}` 数组，循环里逐个 `assert.strictEqual`，一处断言覆盖所有用例，失败信息带用例名。
- `assert` 失败会抛 `AssertionError`，含 `actual` / `expected` / `operator` / `message`，这就是「失败信息」的来源；示例故意埋一个错误用例演示它。
- 断言族：`strictEqual`（`===`）、`deepStrictEqual`（深比较）、`throws`（验证抛错类型）、`ok`（真值）。
- **基准**：用单调时钟 `performance.now()` 量「差值」而非绝对时间；固定迭代次数保证可复现，示例只打印迭代数与类型，不依赖波动的耗时数值。
- **覆盖率思路**：`node --experimental-test-coverage --test` 直接打印覆盖率；`NODE_V8_COVERAGE=dir` 导出 V8 原始数据，再由 c8 / nyc 生成 HTML 报告。这些是工具层能力，脚本本身零依赖。
- **文档注释**：`/** ... */` JSDoc 的 `@param`/`@returns` 由 IDE 与 `jsdoc` 读取；TypeScript 项目用 `tsc --declaration` / `--emitDeclarationOnly` 生成 `.d.ts`，`typedoc` 生成站点（均为外部工具，不引入本仓库）。
- 标准工程用 `node:test` 的 `test()`/`describe()` + `node --test`；本课用自制 runner 是为了让入口脚本「直接跑」且输出完全可控。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 测试框架 | 内置 `testing` / `go test` | `unittest` / `pytest` | JUnit | 内置 `#[cfg(test)]` / `cargo test` | `node:test` / `node --test` |
| 断言 | 手写 `if` + `t.Errorf` | `assert` / `self.assertEqual` | `assertEquals` 等 | `assert!`/`assert_eq!` 宏 | `node:assert` |
| 表驱动 | 惯例（切片子测试） | `subTest` / 参数化 | 参数化测试 | 手写循环 | 数组用例循环 |
| 基准 | `Benchmark` / `go test -bench` | `timeit` | JMH | `cargo bench`（需 nightly） | `performance.now()` / 手写 |
| 覆盖率 | `go test -cover` | `coverage.py` | JaCoCo | `cargo llvm-cov` | `--experimental-test-coverage` / c8 |
| 文档 | `go doc` / godoc | docstring + `pydoc` | Javadoc | `///` + `rustdoc` | JSDoc / `tsc --declaration` |

Go 是这里唯一把「测试 + 基准 + 覆盖率 + 文档」全部塞进一条 `go test`/`go doc` 命令的语言；Rust 的 `#[cfg(test)]` + `cargo test`/`rustdoc` 同样是一等公民。Java 依赖 JUnit + JMH + JaCoCo 三套外部库，Python 也是 `unittest`/`timeit`/`coverage.py` 拼装。Node 的 `node:test` 较新、覆盖率仍带 `--experimental-` 前缀，基准没有官方框架（需 `performance.now()` 或三方库），生态上最接近「自己拼」的 Python。
