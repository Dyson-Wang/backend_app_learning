#!/usr/bin/env python3
"""扫描实际文件，重算进度并重写全部索引。

新增一课或一个场景后跑一次即可，不用手工维护进度表：

    python3 scripts/refresh_indexes.py

判定标准：对应目录里存在该语言的可运行入口文件，且存在 README.md。
"""
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parent.parent

BASICS = [
    ("01-hello-and-toolchain", "hello 与工具链 / 项目结构"),
    ("02-variables-and-types", "变量、类型、常量、零值"),
    ("03-control-flow", "控制流：分支、循环、迭代"),
    ("04-functions-and-errors", "函数与错误处理"),
    ("05-collections-and-iterators", "集合与迭代器"),
    ("06-structs-classes-and-modules", "结构体 / 类、接口、模块"),
    ("07-package-management", "包管理与依赖"),
    ("08-file-io-and-json", "文件 I/O 与 JSON"),
    ("09-concurrency-primitives", "并发原语入门"),
    ("10-first-http-service-with-db", "第一个 HTTP 服务 + 数据层"),
    ("11-operators-and-expressions", "运算符与表达式"),
    ("12-strings-and-encoding", "字符串与字符编码"),
    ("13-enums-and-pattern-matching", "枚举、模式匹配与代数数据类型"),
    ("14-generics-and-type-parameters", "泛型与类型参数"),
    ("15-iterators-and-generators", "迭代器与生成器"),
    ("16-closures-and-higher-order-functions", "闭包与高阶函数"),
    ("17-modules-and-visibility", "模块进阶：可见性与命名空间"),
    ("18-time-regex-and-stdlib", "时间、正则与常用标准库"),
    ("19-concurrent-communication", "并发通信：channel / 队列 / 同步原语"),
    ("20-async-and-cancellation", "异步编程与取消 / 上下文"),
    ("21-memory-model", "内存模型：GC / 所有权 / 生命周期"),
    ("22-reflection-annotations-and-macros", "反射、注解、宏与元编程"),
    ("23-testing-benchmarks-and-docs", "测试、基准与文档"),
    ("24-capstone-http-service-revisited", "综合实战：语言特性用回 HTTP 服务"),
]

SCENARIOS = [
    ("A · 流量与稳定性", [
        ("A01", "A01-rate-limiting", "限流：令牌桶 / 滑动窗口，单机与分布式"),
        ("A02", "A02-timeout-retry-backoff", "超时、重试、指数退避与抖动"),
        ("A03", "A03-circuit-breaker", "熔断与降级"),
        ("A04", "A04-graceful-shutdown", "优雅关闭、健康检查、上游摘除"),
        ("A05", "A05-backpressure", "背压与并发上限"),
    ]),
    ("B · 一致性与并发", [
        ("B01", "B01-transactions-and-isolation", "事务与隔离级别、乐观锁"),
        ("B02", "B02-idempotency", "幂等键与重复请求"),
        ("B03", "B03-distributed-lock", "分布式锁"),
        ("B04", "B04-race-and-deadlock", "竞态、死锁与调试"),
        ("B05", "B05-context-propagation", "上下文取消与跨服务传播"),
    ]),
    ("C · 数据与缓存", [
        ("C01", "C01-connection-pool", "连接池与资源耗尽"),
        ("C02", "C02-n-plus-one", "N+1 查询与批量写"),
        ("C03", "C03-cursor-pagination", "游标分页替代 OFFSET"),
        ("C04", "C04-cache-stampede", "缓存穿透 / 击穿 / 雪崩"),
        ("C05", "C05-cache-db-consistency", "缓存与 DB 双写一致性"),
    ]),
    ("D · 消息与异步", [
        ("D01", "D01-at-least-once-delivery", "至少一次投递与重复消费"),
        ("D02", "D02-dead-letter-queue", "死信队列与重试策略"),
        ("D03", "D03-ordering", "顺序性保证"),
        ("D04", "D04-streaming", "SSE / WebSocket / gRPC 流式"),
    ]),
    ("E · 安全与工程化", [
        ("E01", "E01-jwt-refresh-revoke", "JWT 刷新与吊销"),
        ("E02", "E02-password-rbac", "密码哈希与 RBAC"),
        ("E03", "E03-validation-error-code", "输入校验与统一错误码"),
        ("E04", "E04-observability", "结构化日志 / trace id / 指标"),
        ("E05", "E05-config-secrets-migration", "配置、密钥与 DB migration"),
        ("E06", "E06-testing-containerization", "单测、集成测试与容器化"),
    ]),
    ("F · 故障复盘", [
        ("F01", "F01-slow-query", "慢查询拖垮连接池"),
        ("F02", "F02-retry-storm", "重试风暴"),
        ("F03", "F03-message-backlog", "消息积压"),
    ]),
]

# 每门语言的可运行入口文件名
LANGS = [
    ("01-go", "Go", "main.go"),
    ("02-typescript-node", "TypeScript / Node.js", "index.mjs"),
    ("03-python", "Python", "main.py"),
    ("04-java", "Java", "Main.java"),
    ("05-rust", "Rust", "src/main.rs"),
]

ALL_SCENARIOS = [(code, slug, title) for _, items in SCENARIOS for code, slug, title in items]


def basics_done(lang_dir, entry):
    out = []
    for slug, _ in BASICS:
        d = ROOT / lang_dir / "basics" / slug
        out.append((d / entry).is_file() and (d / "README.md").is_file())
    return out


def advanced_done(lang_dir, entry):
    out = {}
    for code, slug, _ in ALL_SCENARIOS:
        d = ROOT / lang_dir / "advanced" / slug
        out[code] = (d / entry).is_file() and (d / "README.md").is_file()
    return out


def basics_index(label, flags):
    n = sum(flags)
    rows = []
    for (slug, topic), done in zip(BASICS, flags):
        rows.append("| `{}/` | {} | {} |".format(slug, topic, "已实现" if done else "待补"))
    return """# {} · basics 基础入门

统一 **24 课**，编号在 5 门语言里一致，方便打开相同编号目录直接对照。01-10 打基础，11-24 覆盖语法细节与语言特性。

进度：**{}/24**（由 `scripts/refresh_indexes.py` 自动扫描生成）

| 目录 | 主题 | 状态 |
|---|---|---|
{}

## 每课固定结构

- 一个可运行的最小工程（`11`~`24` 全部做成跑完即结束的脚本，不阻塞）
- `README.md` 四段：这个例子讲什么 / 怎么跑 / 关键点 / 和其他语言的差异

## 怎么读

先跑通，再改：把示例里的类型、边界、错误分支故意改错一次，看编译器或运行时怎么反馈。**编译器/运行时的报错信息本身就是这一课的一半内容。**

每课最后都留「和其他语言的差异」小节，读完请把结论补进 `../../99-cross-language/scenario-matrix.md`。
""".format(label, n, "\n".join(rows))


def advanced_index(label, done):
    n = sum(done.values())
    blocks = []
    for group, items in SCENARIOS:
        blocks.append("### {}\n".format(group))
        blocks.append("| 编号 | 目录 | 场景 | 状态 |")
        blocks.append("|---|---|---|---|")
        for code, slug, title in items:
            blocks.append("| {} | `{}` | {} | {} |".format(
                code, slug, title, "已实现" if done[code] else "待补"))
        blocks.append("")
    return """# {} · advanced 生产场景

这里不是语法练习，是生产上真会咬人的问题。每个场景统一结构：

1. **问题复现**：先让故障真的发生（有并发 / 自测脚本）
2. **错误写法**：展示直觉上会写错的版本
3. **正确写法**：可运行的实现
4. **怎么验证**：给出命令与预期输出

场景编号与目录名全局统一，完整对照表见 `../../99-cross-language/scenario-matrix.md`。
公共底座概念在 `../../00-foundations/`，先读那一层再来做场景，否则只会抄代码。

进度：**{}/{}**（由 `scripts/refresh_indexes.py` 自动扫描生成）

{}

## 完成一个场景的判据

能不看代码回答三个问题：故障是怎么发生的？这门语言用哪个机制解决？换成另一门语言时机制换成什么？
""".format(label, n, len(ALL_SCENARIOS), "\n".join(blocks))


def load_conclusions():
    """读取现有矩阵里手填的「差异结论」，刷新时保留，避免每次重写被清空。"""
    p = ROOT / "99-cross-language" / "scenario-matrix.md"
    if not p.is_file():
        return {}
    out = {}
    for line in p.read_text(encoding="utf-8").splitlines():
        parts = [c.strip() for c in line.strip().strip("|").split("|")]
        # 场景编号 A01~F03，基础知识编号 01~24，两者都要保留
        if len(parts) >= 8 and re.fullmatch(r"[A-F]\d{2}|\d{2}", parts[0]):
            out[parts[0]] = parts[-1]
    return out


def matrix(status, conclusions=None):
    conclusions = conclusions or {}
    lines = ["# 生产场景 × 5 语言对照矩阵", "",
             "状态列由 `scripts/refresh_indexes.py` 自动扫描生成；「差异结论」列由你每做完一个场景后亲手填，刷新时会保留。", ""]
    for group, items in SCENARIOS:
        lines.append("### {}\n".format(group))
        lines.append("| 编号 | 场景 | Go | Node.js | Python | Java | Rust | 差异结论 |")
        lines.append("|---|---|---|---|---|---|---|---|")
        for code, _, title in items:
            cells = ["已实现" if status[d][code] else "待补" for d, _, _ in LANGS]
            lines.append("| {} | {} | {} | {} |".format(
                code, title, " | ".join(cells), conclusions.get(code, "待填")))
        lines.append("")
    lines.append("## 基础知识对照区\n")
    lines.append("| 编号 | 主题 | Go | Node.js | Python | Java | Rust | 结论 |")
    lines.append("|---|---|---|---|---|---|---|---|")
    bstatus = {d: basics_done(d, e) for d, _, e in LANGS}
    for i, (_, topic) in enumerate(BASICS):
        cells = ["已实现" if bstatus[d][i] else "待补" for d, _, _ in LANGS]
        code = "{:02d}".format(i + 1)
        lines.append("| {} | {} | {} | {} |".format(
            code, topic, " | ".join(cells), conclusions.get(code, "待填")))
    lines.append("""
## 五个机制速查（其余场景对照时反复用到）

| 能力 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 取消 / 超时 | `context.Context` | `AbortSignal` | `asyncio.timeout` | `orTimeout` | `tokio::time::timeout` |
| 请求上下文 | `context.WithValue` | `AsyncLocalStorage` | `contextvars` | `ThreadLocal` + MDC | `tower` Extension |
| 并发单元 | goroutine | 事件循环 + Promise | Task / 进程 | 平台线程 / 虚拟线程 | `tokio::spawn` |
| 共享状态 | `sync.Mutex` / channel | 单线程无锁 | `asyncio.Lock` / GIL | `synchronized` / Lock | `Arc<Mutex<T>>` |
| 错误表达 | `error` 返回值 | `throw` / rejected Promise | `Exception` | 受检 / 非受检异常 | `Result<T, E>` |
| 内存管理 | GC | GC | 引用计数 + GC | JVM GC | 所有权，编译期 |
| 交付 | 单二进制 | Node 运行时 | CPython 运行时 | JVM 字节码 | 单二进制 |
""")
    return "\n".join(lines)


def progress_table(status, bstatus):
    rows = ["| 模块 | 状态 |", "|---|---|",
            "| 目录骨架 / 总索引 | 完成 |",
            "| `00-foundations/` 六块文档与配置 | 完成 |"]
    for d, label, _ in LANGS:
        rows.append("| {} `basics/` | {}/24 |".format(label, sum(bstatus[d])))
        rows.append("| {} `advanced/` | {}/{} |".format(
            label, sum(status[d].values()), len(ALL_SCENARIOS)))
    rows.append("| `06-elective/` | 说明已完成，代码未开始 |")
    rows.append("| `99-cross-language/` | 矩阵自动生成，28 个场景 + 24 课结论已填满 |")
    return "\n".join(rows)


def main():
    status = {d: advanced_done(d, e) for d, _, e in LANGS}
    bstatus = {d: basics_done(d, e) for d, _, e in LANGS}

    for d, label, _ in LANGS:
        (ROOT / d / "basics" / "README.md").write_text(
            basics_index(label, bstatus[d]).strip() + "\n", encoding="utf-8")
        (ROOT / d / "advanced" / "README.md").write_text(
            advanced_index(label, status[d]).strip() + "\n", encoding="utf-8")

    (ROOT / "99-cross-language" / "scenario-matrix.md").write_text(
        matrix(status, load_conclusions()).strip() + "\n", encoding="utf-8")

    readme = ROOT / "README.md"
    text = readme.read_text(encoding="utf-8")
    table = progress_table(status, bstatus)
    start, end = "<!-- PROGRESS:START -->", "<!-- PROGRESS:END -->"
    block = "{}\n{}\n{}".format(start, table, end)
    if start in text and end in text:
        text = re.sub(re.escape(start) + r".*?" + re.escape(end), lambda _: block, text, flags=re.S)
    else:
        pattern = r"## 进度看板\n.*?(?=\n## )"
        repl = "## 进度看板\n\n状态由 `python3 scripts/refresh_indexes.py` 自动扫描生成，新增示例后重跑即可。\n\n{}\n".format(block)
        text, n = re.subn(pattern, lambda _: repl, text, count=1, flags=re.S)
        if n == 0:
            raise SystemExit("找不到进度看板段落，请检查 README.md")
    readme.write_text(text, encoding="utf-8")

    for d, label, _ in LANGS:
        print("{}: basics {}/{}, advanced {}/{}".format(
            label, sum(bstatus[d]), len(BASICS),
            sum(status[d].values()), len(ALL_SCENARIOS)))


if __name__ == "__main__":
    main()
