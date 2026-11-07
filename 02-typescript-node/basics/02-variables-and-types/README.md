# 02 · Node.js · 变量、类型、常量与零值

## 这个例子讲什么

把 JS 的变量声明（`let` / `const` / `var`）、七种原始类型、`typeof`、大整数边界、`NaN` 判等以及「JS 没有统一零值」一次性摊开。重点看 `const` 只锁绑定不锁对象内容、`typeof null === 'object'` 这类历史坑。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/02-variables-and-types/index.mjs
```

零依赖，无需 `npm install`。

## 关键点

- `let` / `const` 是块作用域且有 TDZ；`var` 是函数/模块作用域并会提升，ESM 里基本不用。
- `const` 禁止重新赋值，但对象/数组内部依然可变；要冻结用 `Object.freeze`（见 06）。
- 七种原始类型：`number` `string` `boolean` `bigint` `symbol` `null` `undefined`；`typeof null` 返回 `'object'` 是历史 bug。
- JS 没有 Go 式零值：未赋值是 `undefined`，不是 `0`/`""`/`false`；显式转换用 `Number()` / `String()` / `Boolean()`。
- 整数安全上限 `Number.MAX_SAFE_INTEGER`（2^53-1），超出精度丢失，用 `bigint`（`1n`）解决。
- 判等优先 `===`（不做隐式转换）；`NaN === NaN` 为假，用 `Number.isNaN`，别用会先转换的全局 `isNaN`。
- 落地 TypeScript 后这些动态类型才被 `let x: number`、`type` / `interface` 在编译期约束，运行时仍是同一套 JS。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 变量声明 | `var` / `:=` | 直接赋值 | `int x` / `var` | `let` / `let mut` | `let` / `const` / `var` |
| 常量 | `const` | 约定全大写 | `final` | `const` | `const`（仅锁绑定） |
| 零值 | 统一零值 `0`/`""`/`nil` | 无，未定义报 `NameError` | 字段有默认零值 | 无，必须初始化 | 无，只有 `undefined` |
| 大整数 | `int64` / `big.Int` | 任意精度 `int` | `BigInteger` | `i128` / 需库 `BigInt` | `number` + `bigint` |
| 类型检查 | 编译期静态 | 运行期动态 | 编译期静态 | 编译期静态 | 运行期动态（TS 才补静态） |

JS 是这里唯一「`const` 不等于不可变对象、且没有统一零值」的语言，与 Go / Java 的默认值语义差别最大。
