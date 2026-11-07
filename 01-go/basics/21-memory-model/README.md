# 21 · 内存模型：GC / 所有权 / 生命周期

## 这个例子讲什么

值类型与引用语义的区别、逃逸到堆的实证、slice/map 共享底层数据、`append` 扩容与「子切片持有大数组」等内存陷阱，以及 Go 追踪式 GC 与对象生命周期的定位。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/21-memory-model
# 看逃逸分析结果（需要 go 工具链）
go build -gcflags=-m ./basics/21-memory-model
```

## 关键点

- **值 vs 引用**：struct 赋值是**按值拷贝**（改副本不动原值）；想共享就用指针。`slice`/`map`/`chan`/`func` 是描述符，内部指向共享的底层数组/哈希表，赋值只拷描述符，本例 `t[0]=9` 让 `s` 跟着变。
- **逃逸分析决定栈/堆**：`escapesToHeap` 返回局部变量 `&x`，栈帧销毁后指针仍必须有效，编译器把它分配到堆——这就是「逃逸」。是否逃逸由编译器静态分析，不写在源码里；`go build -gcflags=-m` 可看报告。
- **GC 而非所有权**：Go 用并发**三色标记-清除** GC，对象是否存活由「可达性」决定，没有引用计数、没有 `free`、没有所有权/借用/生命周期标注。`runtime.GC()` 只是**建议**触发一次，不保证立即回收。
- **陷阱一 `append` 扩容**：容量不足会分配新数组，旧切片与新切片就此分离；本例 `base` 与 `grown` 互不影响。
- **陷阱二 子切片泄漏**：`head := big[:3]` 的 `cap` 仍是 1024，只要 `head` 活着，整个大数组就无法被 GC 回收。修正用 `append([]byte(nil), head...)` 拷到独立小切片。
- **值传递成本**：大结构体按值传会整份拷贝，通常用指针传参。

## 和其他语言的差异

| 语言 | 回收机制 | 值 / 引用 | 逃逸 / 移动 / 借用 |
|---|---|---|---|
| Go | 并发三色标记 GC，无手动 free | struct 值拷贝，slice/map 引用语义 | 逃逸分析决定栈/堆；无所有权/借用 |
| Node.js | V8 分代 GC（Scavenge/Mark-Sweep） | 对象皆引用，原始类型按值 | `WeakRef`/`FinalizationRegistry`；无栈堆概念暴露 |
| Python | **引用计数**为主 + `gc` 处理循环引用 | 一切皆对象引用 | `sys.getrefcount`、`__del__`、弱引用；GIL 影响线程内存模型 |
| Java | JVM GC（G1/ZGC 等），GC Roots 可达性 | 基本类型按值，对象按引用 | `WeakReference`/`SoftReference`；JIT 有逃逸分析与栈上分配 |
| Rust | **无 GC**，编译期所有权 + 借用检查 | 移动语义，`Copy` vs `Clone` | 所有权/生命周期/`&`/`&mut`、`Box`/`Rc`/`Arc`、Drop 时机确定 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
