# 13 · Node.js · 枚举、ADT 与「模式匹配」

## 这个例子讲什么

JS 没有原生 `enum`，也没有模式匹配。用 `Object.freeze` 做枚举常量、用「判别联合」（带 `kind` 字段的对象）+ `switch` 模拟 ADT 与穷尽匹配、`switch(true)` 模拟守卫分支，并用注释给出 TypeScript 的 `as const` / `type` / `never` 等价写法，说明 `.mjs` 里跑不了类型。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/13-enums-and-pattern-matching/index.mjs
```

## 关键点

- `Object.freeze({...})` 只是「常量对象」，不是真枚举：没有反向映射、没有类型，要自己 `Object.fromEntries` 建反查表。
- ADT 在 JS 里就是**带判别字段（tag）的普通对象联合**：`{kind:'circle', r} | {kind:'rect', w, h} | ...`。
- `switch (s.kind)` 按 tag 分支，配 `const { r } = s` 解构取数据；`default` 抛错做**运行时**穷尽兜底。
- TS 里靠 `const _exhaustive: never = s`（`assertNever`）让**编译器**在漏分支时报错；JS 没有这层保障，未知 tag 只能运行时才发现。
- `switch (true)` + 条件表达式可模拟「守卫式」匹配，但远不如真正的模式匹配表达力强。
- `as const` + `typeof X[keyof typeof X]` 是 TS 里从常量对象导出字面量联合类型的惯用法；这些类型信息编译后全部被擦除（14 课）。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 枚举 | `iota` + 自定义类型 | `enum.Enum` | `enum`（可带字段） | `enum`（可带数据） | 无，`Object.freeze` 模拟 |
| ADT | 接口 + `type switch` | `dataclass` + `Union` | `sealed interface` + `record` | `enum` 原生 | 判别联合对象 |
| 模式匹配 | `type switch` / `switch` | 3.10 `match case` | 17+ `instanceof` 模式（预览） | 强大 `match` | 无，`switch` 值匹配 |
| 穷尽性 | 不强制 | 不强制 | sealed 部分保证 | 编译器强制 | 运行时 default / TS `never` |
| 解构 | 无 | `match` 捕获 | record 解构（预览） | `let` 解构 | `const {x} = obj` |

Rust 的带数据 `enum` + `match` 是这里表达力最强、且编译器强制穷尽的；Java 17 的 `sealed` + `record`、Python 3.10 的 `match case` 各有原生支持。JS/TS 只能靠判别联合 + `switch` 手工模拟，穷尽性在纯 JS 里完全没有编译器兜底，是最大的短板。
