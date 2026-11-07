# 01 · Node.js · hello 与工具链

## 这个例子讲什么

确认 Node 版本、ESM 模块形态，并用 4 行代码直观看到事件循环：同步代码先跑完，`setTimeout(0)` 的回调最后执行。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/01-hello-and-toolchain/index.mjs
# 期望输出顺序：1: sync / 2: sync / 3: timer callback
```

## 关键点

- `.mjs` 强制 ESM；项目用 `package.json` 里 `"type": "module"` 时 `.js` 也是 ESM。本仓库 `basics/` 全部零依赖，不需要 `npm install`。
- Node 只有一个 JS 执行线程，**I/O 交给 libuv 线程池**，完成后把回调放回事件队列。输出顺序就是这件事的证据。
- 代价：任何同步阻塞（大 JSON 解析、CPU 密集循环、`*Sync` 文件 API）都会卡住整个进程的所有请求。
- `process.on('unhandledRejection')` 只是示例；生产应该配合进程管理器与 graceful shutdown（见 `10`）。
- 后续引入 TypeScript 时，用 `tsc --noEmit` 做类型检查、用 `tsx` 或 `node --experimental-strip-types` 运行。

## 和其他语言的差异

| 语言 | 入口 | 交付 |
|---|---|---|
| Node.js | 直接执行 `.mjs` | 源码 + Node 运行时 |
| Go | `package main` | 单二进制 |
| Python | 直接执行 `.py` | 源码 + CPython |
| Java | `main` 方法 | 字节码 + JVM |
| Rust | `fn main` | 单二进制 |

Node 和 Python 是这里唯二「源码即交付」的，部署时必须保证运行时版本一致。

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
