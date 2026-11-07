# 21 · Java · 内存模型：GC / 生命周期

## 这个例子讲什么

栈与堆、值拷贝 vs 引用拷贝、`==`/`equals` 与 `Integer` 缓存陷阱、GC Roots 可达性、`System.gc()` 不确定性、`WeakReference`、`try-with-resources` 确定性清理、常见内存泄漏。

## 怎么跑

```
cd backend-roadmap/04-java/basics/21-memory-model
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **基本类型是值、对象是引用**：`int y = x` 复制值，`Holder h2 = h1` 复制引用（改 `h2.n` 会影响 `h1`）。
- 局部变量/引用在**栈帧**，对象实体在**堆**；JIT 的逃逸分析可能把未逃逸对象标量替换/栈上分配，但语义上仍由 GC 管理。
- **GC Roots**（栈局部变量、静态字段、活动线程、JNI 引用等）出发的可达链决定对象生死；不可达才被回收。
- **`System.gc()` 只是建议**，可能忽略、可能延迟，绝不依赖其时机——示例调用它但不读取结果，保证输出稳定。
- `WeakReference` 不阻止回收；**`try-with-resources`（`AutoCloseable`）** 才是确定性的 `close()` 时机，对应 RAII。
- `Integer` 缓存 `-128..127` 导致 `127 == 127` 为 `true`、`128 == 128` 为 `false`：比较包装类内容一律用 `equals`。
- 静态集合长期持有对象引用是最常见的内存泄漏。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 回收机制 | JVM GC（分代/ZGC 等） | GC（三色标记） | V8 GC | 引用计数 + gc 模块 | **无 GC，所有权** |
| 值/引用 | 基本类型值 / 对象引用 | 值 / 指针 / slice | 原始值 / 对象引用 | 一切皆对象引用 | 值 / 引用 `&T` |
| 弱引用 | `WeakReference` | 无（1.24 有 weak pointer） | `WeakRef`/`WeakMap` | `weakref` | `Weak<T>` |
| 确定性释放 | `try-with-resources`/`AutoCloseable` | `defer` | 无（`using` 提案） | `with`/`__del__` | `Drop` |
| 主要陷阱 | 静态集合持有 / 缓存身份 | 逃逸到堆 | 闭包持有 | 引用循环 / GIL | 悬垂引用（编译期挡住） |
