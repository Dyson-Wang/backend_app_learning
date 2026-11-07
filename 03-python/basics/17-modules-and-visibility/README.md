# 17 · Python · 模块进阶：可见性与命名空间

## 这个例子讲什么

`_x` 私有约定、`__all__` 控制 `from pkg import *`、名称改写 `__x` → `_Class__x`、`import ... as` 别名、包目录的 `__path__` / `__package__`、用「函数内惰性 import」打破循环依赖、`if __name__ == "__main__"` 守卫。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/17-modules-and-visibility/main.py
```

## 关键点

- Python **没有强制 private**：`_x` 只是约定（不挡 `from x import *` 之外的显式访问）；`__x` 在类里被改写成 `_Class__x`，仍可按改写名取到。
- **`__all__` 决定 `from pkg import *` 导出哪些名字**，是模块公开 API 的唯一声明；没写 `__all__` 时星号导入会跳过所有 `_` 开头名字。
- 包 = 带 `__init__.py` 的目录，包对象有 `__path__`；`import ... as` 给模块起别名，解决同名冲突。
- **循环依赖**：两个模块顶层互相 `import` 会卡在半初始化状态；把 `import` 挪进函数（惰性导入）或抽公共模块即可打破。
- `if __name__ == "__main__"` 让文件既可被导入、又能直接运行；脚本方式运行时 `__spec__` 为 `None`。

## 和其他语言的差异

- **Go**：大写首字母导出、小写包内私有、`internal/` 目录是编译器强制的可见性边界、`package` 名即目录。
- **Node.js/TS**：ESM `export` / `import` / `default`、barrel `index`、类字段 `#private` 才是真私有。
- **Java**：`public/protected/private/package-private` 四档关键字、`module-info.java`（JPMS）模块系统。
- **Rust**：`pub` / `pub(crate)` / `pub(super)` 分级可见性、`mod` 树、crate 边界。
- **Python**：`_` 约定私有、`__all__` 管星号导入、`__x` 名称改写、无 `private` 关键字、`__init__.py` 包、惰性 import 破循环。
