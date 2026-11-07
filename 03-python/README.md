# Python · 03-python

**运行时模型**：GIL + asyncio / 多进程 · GC（引用计数 + 分代）· 动态类型 · CPython 运行时

## 环境准备

本机 Python 3.13 可用，示例已现场验证。零第三方依赖：

```bash
cd 03-python
python3 basics/01-hello-and-toolchain/main.py
python3 basics/10-first-http-service-with-db/main.py   # 起服务
```

`advanced/` 讲到 FastAPI / SQLAlchemy / asyncio 时会引入依赖，届时在每个目录单独给 `requirements.txt`，不污染全局环境：

```bash
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
```

## 目录

- `basics/` —— 基础入门 10 课（统一编号）
- `advanced/` —— 生产场景（统一 A01 ~ F03 编号）

## 重点补什么

GIL 决定了「CPU 密集用多进程、I/O 密集用 asyncio」这条分界线；动态类型决定了运行时错误更晚暴露，所以 `advanced/E03`（输入校验）和 `E06`（测试）在这门语言里格外重要。

## 与你现有环境的关系

homelab 的运维自动化、备份脚本、AI 服务接口，Python 是最顺手的一门。示例会尽量复用 `docker-services` 里已有的 PostgreSQL / Redis，而不是另起一套。

## 跨语言对照

每课 / 每个场景做完，把「这门语言和其他语言的差异」补到 `../99-cross-language/scenario-matrix.md`。同一编号横向打开，是这套库唯一正确的读法。
