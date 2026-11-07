# 07 · Python · 包管理与依赖

## 这个例子讲什么

演示导入机制的运行时基础：`sys.path` / `sys.modules` 缓存、`importlib.util.find_spec`、`importlib.import_module`、`types.ModuleType` 动态造包、`import ... as` 别名、包目录的 `__path__` 标记，并说明 pip / venv / pyproject.toml 的职责（本课不创建任何依赖文件）。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/07-package-management/main.py
```

## 关键点

- **`sys.path` 决定去哪找模块，`sys.modules` 是已导入模块的缓存**；重复 `import` 只取缓存，不重复执行。
- `find_spec()` 只查找不执行，用来「探测某个依赖是否可用」最合适。
- **每个 `.py` 文件都是模块，带 `__init__.py` 的目录是包**，包对象有 `__path__`（普通模块没有）。
- `import ... as` / `from ... import ... as` 解决命名冲突；动态 `ModuleType` + 注册 `sys.modules` 展示导入的最小模型。
- 真实项目：`python3 -m venv .venv` 隔离，依赖写进 `pyproject.toml`（或 `requirements.txt`），由 `pip` 装进 `site-packages`；`__pycache__` 是字节码缓存，可随时删除。

## 和其他语言的差异

- **Go**：`go.mod` / `go.sum`、module 取代 `GOPATH`、`go get`、`internal/` 可见性边界。
- **Node.js/TS**：`package.json` / `node_modules`、`npm` / `pnpm`、`package-lock.json`、ESM 与 CJS 两套解析。
- **Java**：Maven `pom.xml` / Gradle、`~/.m2` 本地仓库、`classpath` / `module-info.java`。
- **Rust**：`Cargo.toml` / `Cargo.lock`、`crates.io`、`cargo build` 产出 `target/`。
- **Python**：`pip` + `pyproject.toml` / `requirements.txt`、`venv`、`site-packages`、`sys.path`、`importlib`、`__pycache__`。
