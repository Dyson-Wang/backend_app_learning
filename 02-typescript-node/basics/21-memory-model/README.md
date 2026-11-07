# 21 · Node.js · 内存模型：GC / 值与引用 / 生命周期

## 这个例子讲什么

演示 JS 的值与引用语义（原始值拷贝、对象引用共享、浅拷贝 vs `structuredClone` 深拷贝）、栈与堆（V8 分代堆、`Buffer` 堆外内存）、追踪式 GC 与 GC Roots、`WeakRef`/`WeakMap`/`FinalizationRegistry` 弱引用语义，以及三类常见内存陷阱（全局缓存、闭包持有大对象、未清理的 timer/listener）。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/21-memory-model/index.mjs
```

只用 `node:v8` / `node:process`；GC 相关只演示 API，不依赖不可预测的回收时机，输出稳定。

## 关键点

- 原始值按**值**传递（拷贝），对象按**引用**传递（多个变量指向同一块堆内存）；`{...obj}` 是浅拷贝，`structuredClone` 才是深拷贝。
- JS 没有栈上对象/堆上对象的显式区分：`Buffer`/`ArrayBuffer` 分配在 **V8 堆外的 external memory**；`v8.getHeapStatistics()` / `process.memoryUsage()` 看堆指标。
- V8 是**追踪式分代 GC**：新生代 Scavenge、老生代 Mark-Sweep/Mark-Compact，从 GC Roots 标记可达对象并回收不可达对象，**不是引用计数**，开发者无法手动 `free`。
- `WeakRef.deref()`、`WeakMap`/`WeakSet` 不产生强引用，不阻止回收；`FinalizationRegistry` 回调是否/何时触发**不确定**，绝不能用来做业务逻辑（示例只演示 API 存在）。
- 对象生命周期：创建 → 被强引用（可达）→ 最后一个强引用消失（不可达）→ 某次 GC 回收。`strong = null` 后只是「等待 GC」。
- 常见陷阱：全局 `Map` 缓存只增不减、闭包捕获大对象、忘记 `clearTimeout`/`off` 的监听器——本质都是「意外的强引用让对象一直可达」。
- JS 没有 Rust 的所有权/移动/借用、没有 Go 的逃逸控制；V8 的逃逸分析可能把短命对象优化到栈上，但对开发者完全透明。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 回收机制 | 追踪 GC（三色标记） | 引用计数 + 循环 GC | 追踪 GC（G1/ZGC） | 编译期所有权 + `Drop` | 追踪分代 GC（V8） |
| 值/引用 | 值类型 + 指针/切片 | 全对象引用 | 基本类型值 + 对象引用 | 值 + `&`/`Box` | 原始值 + 对象引用 |
| 栈/堆控制 | 逃逸分析自动 | 全在堆 | 全在堆（标量替换优化） | 编译期明确 `Box`/栈 | V8 透明，Buffer 在堆外 |
| 弱引用 | 无（1.24 `weak.Pointer` 仍新） | `weakref` | `WeakReference` | `Weak` / `Rc::downgrade` | `WeakRef` / `WeakMap` |
| 移动/借用 | 无 | 无 | 无 | `move` / 借用检查 | 无 |
| 内存陷阱 | 切片持有底层数组 | 引用循环、`__del__` | 静态集合、监听器泄漏 | 循环 `Rc` | 全局缓存、闭包、未清理 listener |

Rust 用所有权/借用/移动在**编译期**管理内存，完全没有 GC；Go/Java/JS 都靠运行时追踪 GC，但 Go 有值语义与逃逸分析、Java 有 `WeakReference` 与多种 GC 器。Python 以引用计数为主、另配循环检测器，是这里最特殊的。JS 既没有弱类型系统约束、也没有手动释放手段，因此「意外强引用」是最常见的泄漏来源，`WeakMap`/`WeakRef` 是缓解工具而非根治手段。
