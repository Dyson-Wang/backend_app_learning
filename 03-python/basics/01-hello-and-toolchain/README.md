# 01 · Python · hello 与工具链

## 这个例子讲什么

确认 CPython 版本，看动态类型、类型注解、`f-string` 与 `if __name__ == "__main__"` 惯用法，并点出 GIL 这条贯穿全路线的分界线。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/01-hello-and-toolchain/main.py
# 依赖隔离（advanced/ 用到第三方库时）
python3 -m venv .venv && source .venv/bin/activate
```

## 关键点

- **类型注解不强制**：`name: str` 只是给 `mypy` / IDE 看，运行时传 int 也照样跑。要运行时校验得用 pydantic（FastAPI 的核心）。
- `if __name__ == "__main__"` 让文件既能被 import 又能直接执行。
- **GIL**：CPython 同一时刻只有一个线程执行字节码。I/O 密集用 `asyncio` / 线程池能提速，CPU 密集必须上 `multiprocessing` 或用释放 GIL 的 C 扩展（numpy）。
- 依赖必须装进虚拟环境，别污染系统 Python；本仓库 `basics/` 零第三方依赖。
- 缩进即语法块，没有花括号；这一条让 Python 代码天然统一，但也让「复制粘贴错了缩进」变成语义错误而非格式错误。

## 和其他语言的差异

| 语言 | 类型检查 | 并发上限 | 交付 |
|---|---|---|---|
| Python | 动态，注解可选 | GIL，CPU 密集需多进程 | 源码 + CPython |
| Go | 静态，编译期 | 多核并行 | 单二进制 |
| Node.js | 动态（TS 静态） | 单线程事件循环 | 源码 + Node |
| Java | 静态强类型 | 多线程 / 虚拟线程 | 字节码 + JVM |
| Rust | 静态 + 所有权 | 多核并行，编译期防数据竞争 | 单二进制 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md`。
