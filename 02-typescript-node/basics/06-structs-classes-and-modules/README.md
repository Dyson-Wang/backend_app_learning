# 06 · Node.js · 类、接口（鸭子类型）与模块

## 这个例子讲什么

演示 `class` 语法：构造器、`#私有字段`、`static`、`getter/setter`、`extends`/`super`、`instanceof`，并用「鸭子类型」说明 JS 没有运行时接口——TypeScript 的 `interface` 编译后会被擦除。最后点明 ESM 的静态导入与顶层 `this`。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/06-structs-classes-and-modules/index.mjs
```

只用 `node:util` 的 `inspect`，零依赖。

## 关键点

- `#field` 是语言级真私有，`Reflect.ownKeys` 都取不到；老代码里的 `_field` 只是命名约定。
- `static` 属于类，`get/set` 让方法像属性一样读写，`extends` + `super` 复用父类逻辑。
- JS 无 `interface` / `struct`：对象只要「长得像」就能用（鸭子类型）。TS 用 `interface`/`type` 做**编译期**结构检查，运行期完全擦除（14 课细讲）。
- 「结构体」风格可用 `Object.freeze({...})` 做不可变记录；深冻结要自己递归或用 `structuredClone`。
- ESM 的 `import`/`export` 是静态声明，顶层 `this === undefined`，`import.meta.url` 给出模块自身路径（17 课进阶）。
- 类本质仍是构造函数 + 原型（`prototype`），`instanceof` 查的是原型链，不是「类型标签」。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 数据载体 | `struct` | `class` / `dataclass` | `class` / `record` | `struct` / `enum` | `class` / 对象字面量 |
| 接口 | 隐式满足 `interface` | `Protocol`/ABC（鸭子） | 显式 `implements` | `trait` | 无，鸭子类型（TS 编译期） |
| 私有 | 首字母大小写 | `_name` 约定 | `private` | 模块 `pub` 可见性 | `#field` 真私有 |
| 继承 | 嵌入组合，无继承 | 多继承 | 单继承 `extends` | 无继承，用 trait | 原型链 `extends` |
| 模块 | `package` + `import` | `import` 模块 | `package` + JPMS | `mod` / `crate` | ESM `import`/`export` |

Go 用组合 + 隐式接口、Rust 用 trait + 无继承，Java 是最「正统」的名义类型 + `implements`；JS/TS 则是结构类型（长得像即可），这是最根本的分歧。类的底层是原型而非类，也和其余四门都不同。
