# 21 · Python · 内存模型：GC / 所有权 / 生命周期

## 这个例子讲什么

值与引用（别名共享）、`sys.getrefcount` 观察引用计数、`__del__` 即时析构、**引用循环**靠 `gc.collect()` 回收、`weakref` 弱引用与回调、`gc.is_tracked` / 分代阈值、`tracemalloc` 堆分配追踪，以及 GIL 对内存模型的影响。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/21-memory-model/main.py
```

## 关键点

- CPython 以**引用计数为主**：计数归零立刻回收并触发 `__del__`；`sys.getrefcount(x)` 自身会临时 +1，所以比较「增量」而非绝对值。
- 引用计数解决不了**循环引用**（`x.peer=y; y.peer=x`），由**分代 GC**（`gc.collect()`）定期清理；`weakref` 不增加计数，可无副作用地观察对象何时被回收。
- `gc.is_tracked([])` 为 `True`、`gc.is_tracked(小整数)` 为 `False`（不可变且无容器的对象不进 GC 跟踪）。
- **GIL** 让 `ob_refcnt` 的增减成为原子字节码、线程安全；但 `counter += 1` 这种复合操作仍非原子，必须 `Lock`（09/19 课）。
- Python 没有所有权/移动/借用：对象全在堆上，栈帧与变量只存引用；不像 Rust 有编译期生命周期，也不像 Go 有逃逸分析。

## 和其他语言的差异

- **Go**：三色并发标记 GC、逃逸分析决定栈/堆、指针但无所有权、`unsafe` 可绕过。
- **Node.js/TS**：V8 分代 GC（Scavenge / Mark-Sweep）、闭包易致内存驻留、`WeakRef` / `WeakMap`。
- **Java**：JVM GC Roots 可达性、`System.gc()` 不保证执行、`WeakReference`、`try-with-resources`。
- **Rust**：所有权 / 移动 / 借用 / 生命周期编译期检查，`Drop` 确定性释放，**无 GC**。
- **Python**：引用计数 + 分代 GC、`sys.getrefcount`、`__del__`、`weakref`、`gc.collect()`、GIL 保护引用计数。
