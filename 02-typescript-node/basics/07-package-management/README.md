# 07 · Node.js · 包管理与依赖

## 这个例子讲什么

用 `import.meta.resolve`、`node:module` 的 `builtinModules`、`process.versions` 讲清 Node 的模块解析与依赖组织：`node:` 前缀的意义、ESM 解析规则、内置模块清单、版本表，以及 `package.json` / semver / lockfile / npm 命令各自解决什么问题。全程不联网、不创建 `package.json`。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/07-package-management/index.mjs
```

## 关键点

- `node:` 前缀强制走内置模块，避免 `node_modules` 里同名包劫持；`import.meta.resolve('fs')` 会被规范化成 `'node:fs'`。
- `import.meta.resolve(specifier)` 在 Node 22 是**同步**的，找不到包抛 `ERR_MODULE_NOT_FOUND`，可用来做依赖探测。
- `builtinModules` 返回不带 `node:` 前缀的名字；理解它才能区分「内置」与「第三方」。
- `process.versions` 是编译进二进制的版本表（`node`/`v8`/`uv` 等），不是运行时去查网络。
- `package.json`：`"type":"module"` 决定 `.js` 按 ESM 解析、`"exports"` 精确控制包暴露的子路径（封装边界）。
- semver：`^` 放开 minor+patch、`~` 只放开 patch；`package-lock.json` 锁整棵树，`npm ci` 严格按锁装以保证可复现。
- 本仓库坚持零第三方依赖，所有课程只用 `node:` 内置模块，因此不需要 `npm install`。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 清单文件 | `go.mod` | `pyproject.toml`/`requirements.txt` | `pom.xml`/`build.gradle` | `Cargo.toml` | `package.json` |
| 锁文件 | `go.sum` | `poetry.lock`/`uv.lock` | 无原生 | `Cargo.lock` | `package-lock.json` |
| 包管理器 | `go mod` | `pip` / `uv` / `poetry` | Maven / Gradle | `cargo` | `npm` / `pnpm` / `yarn` |
| 中央仓库 | proxy.golang.org | PyPI | Maven Central | crates.io | npm registry |
| 内置前缀 | 标准库 import path | 标准库直接 import | JDK 标准库 | `std` | `node:` |
| 版本格式 | 语义版本 + `+incompatible` | PEP 440 | 无强规范 | semver | semver |

Rust 的 `Cargo.toml` + `Cargo.lock` 与 Node 的 `package.json` + lock 思路最接近；Go 用 `go.mod` 且默认最小版本选择（MVS），Java 则长期没有官方锁文件，依赖解析最依赖构建工具。Node 独特的 `node:` 前缀是其余四门都没有的命名空间隔离手段。
