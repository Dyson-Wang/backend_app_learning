# TypeScript / Node.js · 02-typescript-node

**运行时模型**：单线程事件循环 + libuv 非阻塞 I/O · GC · 渐进静态类型 · Node 运行时

## 环境准备

本机 Node 22 可用，示例已现场验证。零第三方依赖，用 `.mjs` + Node 内置模块：

```bash
cd 02-typescript-node
node basics/01-hello-and-toolchain/index.mjs
node basics/10-first-http-service-with-db/index.mjs   # 起服务，Ctrl-C 结束
```

`basics/` 先用纯 ESM 保证零依赖跑通；后续讲到类型时会引入 `typescript` / `tsx`，并在对应 README 说明。

## 目录

- `basics/` —— 基础入门 10 课（统一编号）
- `advanced/` —— 生产场景（统一 A01 ~ F03 编号）

## 重点补什么

单线程也能高并发——代价是**任何同步阻塞都会拖垮整个进程**。这一门的 `advanced/` 会特别关注事件循环延迟、背压、流（stream）以及 `AbortSignal` 的取消传播。

## 与你现有项目的对照

`linnet` 已经是 TypeScript（Vite + React + Zod）。这里只补服务端心智：请求上下文用 `AsyncLocalStorage`、校验用 Zod、流用 `ReadableStream`，类型和工具链可以完全复用，是 5 门里切换成本最低的一门。

## 跨语言对照

每课 / 每个场景做完，把「这门语言和其他语言的差异」补到 `../99-cross-language/scenario-matrix.md`。同一编号横向打开，是这套库唯一正确的读法。
