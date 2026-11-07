# 17 · Node.js · 模块进阶：可见性与命名空间

## 这个例子讲什么

用 `data:` URL 动态 import 在单文件里造出**真实的 ES 模块实例**，讲清导出/私有规则（未 `export` 的绑定不可见）、Module Namespace 对象（只读、键有序、`Symbol.toStringTag === 'Module'`）、命名冲突与别名（`import {x as y}` / `export {a as b}`）、活绑定 `live binding`，以及目录即包、`exports` 可见性边界与循环依赖的 TDZ 陷阱。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/17-modules-and-visibility/index.mjs
```

## 关键点

- 模块作用域天然私有：顶层 `const secret` 没 `export`，命名空间上就是 `undefined`；只有导出的绑定对外可见。
- `import * as ns` 得到 **Module Namespace 对象**：属性只读（严格模式赋值抛 `TypeError`）、键按字典序排列、`[Symbol.toStringTag] === 'Module'`。
- 别名：`import { value as v } from '...'` 是静态写法；动态 `import()` 后用 `const { value: v } = ns` 等价；模块内可 `export { a as b }`。
- `export default` 对应 `ns.default`；命名导出与默认导出可共存。
- **活绑定**：`export let counter` 被模块内部改写后，外部 `ns.counter` 同步变化（读的是同一绑定，不是拷贝）——这正是 ESM 能处理循环依赖的基础。
- **目录即包**：含 `package.json` 的目录就是包；`exports` 字段列出对外子路径，未列出的内部文件 import 报 `ERR_PACKAGE_PATH_NOT_EXPORTED`，构成可见性边界（07 课）。
- **循环依赖**：ESM 把「实例化」与「求值」分离，配合活绑定可以加载 A↔B；但求值期读取对方尚未初始化的 `let/const` 会命中 TDZ 抛 `ReferenceError`。Node 不会死锁，但要避免在顶层读取循环依赖的值。
- CommonJS 是 `module.exports` 的动态对象、值拷贝、可写，与 ESM 命名空间语义完全不同。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 单元 | `package`（目录） | 模块 / 包目录 | `package` / JPMS 模块 | `crate` / `mod` | ESM 模块 / 包目录 |
| 可见性 | 首字母大小写 | `_` 约定 / `__all__` | `public/private/protected` + JPMS `exports` | `pub` / `pub(crate)` | `export` + `exports` 字段 |
| 命名空间 | 包名 `pkg.X` | `import pkg; pkg.x` | 全限定类名 | `crate::mod::x` | `import * as ns` 命名空间对象 |
| 别名 | `import alias "path"` | `import x as y` | 无（全限定名） | `use x as y` | `import {x as y}` / `export {a as b}` |
| 循环依赖 | 编译报错（import cycle） | 运行时可能半成品 | 允许但易出初始化问题 | 允许（同 crate） | 允许（活绑定 + TDZ 风险） |
| 内部边界 | `internal/` 目录约定 | `_` 前缀约定 | JPMS `exports` 精确控制 | `pub(crate)` | `package.json` `exports` |

Go 直接禁止 import cycle，Rust/Java 靠 `pub(crate)` / JPMS `exports` 控制包内可见性，Python 靠 `_` 命名约定；只有 Node 用 `exports` 字段在**包解析层**强制拦截未暴露的子路径。ESM 的活绑定也是这里独有的机制：Go/Java/Python 的导入更像「值的拷贝/引用」，`export let` 同步更新命名空间的语义在其他语言里没有直接对应。
