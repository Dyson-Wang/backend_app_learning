# 22 · Node.js · 反射、注解、宏与元编程

## 这个例子讲什么

用 `Reflect`（`get`/`set`/`has`/`apply`/`construct`/`ownKeys`）与 `Proxy`（`get`/`set`/`has` 陷阱）演示运行时类型信息、按名字访问字段/方法、动态构造与调用；用属性描述符近似「注解/标签」；再用两个边界实验说明 Proxy 不能违背对象不变式、也拦不到 `#` 私有字段。最后讲清 JS 没有装饰器/注解/宏，TS 装饰器元数据是编译器产物。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/22-reflection-annotations-and-macros/index.mjs
```

零依赖，`Reflect`/`Proxy` 都是语言内置。

## 关键点

- 运行时类型信息：`typeof` / `instanceof` / `constructor.name` / `Object.prototype.toString.call` 各有局限（如 `typeof null === 'object'`）。
- 按名字访问：`obj[name]`、`Reflect.get/set/has`、`Reflect.apply(fn, thisArg, args)`、`Reflect.construct(Ctor, args)` 让「运行时才知道的属性/方法」可被调用。
- `Reflect` 的 13 个静态方法与 Proxy 陷阱一一对应，**陷阱内用 Reflect 转发默认行为**是标准写法，避免自己重造语义。
- JS 无注解/标签：最接近的是 `Object.defineProperty` 的属性描述符（`writable`/`enumerable`/`configurable`），`Object.getOwnPropertyDescriptor` 可读；`Reflect.ownKeys` 能看到 `Object.keys` 看不到的不可枚举与 symbol 键。
- **Proxy 边界**：对不可写且不可配置的自身属性，`get` 陷阱不能谎报别的值（违反不变式会抛 `TypeError`）；`#` 私有字段不是普通属性操作，Proxy 拦不到、`Reflect.ownKeys` 也看不到。
- **动态调用代价**：Proxy 破坏 V8 的内联缓存与隐藏类优化，比直接调用慢；只应在框架、ORM、校验层使用，别放热点路径。
- JS 没有原生装饰器（TC39 Stage 3 推进中）、没有注解、没有宏、没有编译期代码生成；TypeScript 的装饰器 + `emitDecoratorMetadata` 会由 `tsc` 生成 `design:type`/`design:paramtypes` 元数据，配合 `reflect-metadata` 运行时读取，但这是**编译器产物**而非语言运行时能力，纯 `.mjs` 跑不了。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 运行时反射 | `reflect` 包 | `inspect`/`getattr`/`__annotations__` | `java.lang.reflect` | 极弱（`Any`/`TypeId`） | `Reflect` / `Proxy` |
| 按名访问 | `FieldByName`/`MethodByName` | `getattr`/`setattr` | `Class.getField`/`Method.invoke` | 不支持 | `obj[name]` / `Reflect.get` |
| 注解/标签 | struct tag | 装饰器 / 类型注解 | Annotation + 反射读取 | 属性宏 `#[...]` | 属性描述符 / TS 装饰器 |
| 宏/代码生成 | `go:generate` | 元类 / 装饰器（运行时） | APT / Lombok（编译期） | `macro_rules!` / 过程宏 | 无（tsc/Babel 变换） |
| 动态调用开销 | reflect 较慢、类型不安全 | 灵活但动态 | `Method.invoke` 慢 | 无动态调用 | Proxy 破坏内联缓存 |
| 编译期元编程 | 弱 | 无 | APT | 强（宏） | 无 |

Rust 是这里「编译期元编程最强、运行时反射最弱」的一端：靠 `macro_rules!`/过程宏生成代码，几乎没有运行时反射；Go 的 struct tag 是编译进二进制的字符串，靠 `reflect` 读取；Java 有成熟的 Annotation + `java.lang.reflect`；Python 万物皆对象、`getattr` 与装饰器最灵活。JS/Node 处在中间：有强大的 `Proxy` 运行时拦截，却完全没有注解与宏，TS 的装饰器元数据只能算「编译器打补丁」，和 Java/Rust 的原生机制不在同一层。
