# 18 · Node.js · 时间、正则与常用标准库

## 这个例子讲什么

用**固定时刻**演示时间点/时长/时区（`Intl.DateTimeFormat`）、格式化与解析；演示正则的 `exec`/命名捕获组/`matchAll`/`replace` 回调与 `lastIndex` 陷阱；用固定种子 LCG 给可复现随机序列、`node:crypto` 只演示 API；再覆盖排序（`toSorted`/`Intl.Collator`）与 `util.format`/`inspect`/`types`、`structuredClone`、`Object.groupBy`、`Intl.NumberFormat` 等常用工具。

## 怎么跑

```
cd backend-roadmap/02-typescript-node
node basics/18-time-regex-and-stdlib/index.mjs
```

只用 `node:crypto` / `node:util` / `node:perf_hooks`。

## 关键点

- 时间用**固定时间戳**演示格式化/差值，不打印当前绝对时间；`Intl.DateTimeFormat` 配 `timeZone` 是 JS 标准时区方案（`Date` 本身只存 UTC 毫秒 + 本地时区偏移）。
- `Date` 没有「时长」类型，时长就是两个时刻相减的毫秒数；单调时钟用 `performance.now()`。
- 正则：`g` 正则的 `lastIndex` 是**有状态**的，`exec` 会推进它，复用前必须归零，否则 `matchAll` 漏掉开头；`s` 标志让 `.` 匹配换行，命名组用 `(?<name>...)` + `groups.name`。
- `replace` 配回调函数可做分组级替换；`matchAll` 返回所有匹配的迭代器。
- 随机数：`crypto.randomInt`/`randomUUID`/`randomBytes` 是密码学安全随机；输出要稳定就只演示 API（打印长度/类型），或像示例里用**固定种子 LCG** 产生可复现序列。
- `sort` 原地改数组，`toSorted` 返回新数组；本地化排序用 `Intl.Collator`（12 课）。
- 工具函数：`util.format/inspect/types`、`structuredClone`（深拷贝，替代 `JSON.parse(JSON.stringify(x))`）、`Object.groupBy`（Node 21+）、`Intl.NumberFormat`。

## 和其他语言的差异

| 机制 | Go | Python | Java | Rust | Node.js |
|---|---|---|---|---|---|
| 时间点 | `time.Time` | `datetime` | `java.time.Instant` | `SystemTime`（需 chrono） | `Date`（UTC 毫秒） |
| 时长 | `time.Duration` | `timedelta` | `Duration` | `Duration` | 毫秒数（无专用类型） |
| 时区 | `time.LoadLocation` | `zoneinfo` | `ZoneId` | chrono `Tz` | `Intl.DateTimeFormat` |
| 格式化 | `Format("2006-01-02")` 参考时间 | `strftime` | `DateTimeFormatter` | chrono `format` | `Intl` / `toISOString` |
| 正则 | `regexp`（RE2） | `re` | `java.util.regex` | `regex` crate（三方） | 内置 `RegExp` |
| 命名捕获 | `(?P<name>...)` | `(?P<name>...)` | `(?<name>...)` | `(?P<name>...)` | `(?<name>...)` |
| 随机数 | `math/rand` / `crypto/rand` | `random` / `secrets` | `Random` / `SecureRandom` | `rand` crate（三方） | `node:crypto` |

Go 的 `time.Time` 自带时区与 `Duration`、用「参考时间」`2006-01-02` 做格式串是独一份；Java `java.time`、Python `datetime` 都有专门的时长类型，而 JS 的 `Date` 只是个 UTC 毫秒数，时区格式化全靠 `Intl`。正则方面 Go 用 RE2（无回溯、无环视），JS/Python/Java 是回溯引擎；命名组语法各家略有差异（Go/Python 用 `?P<>`，JS/Java 用 `?<>`）。随机与排序在 Rust 标准库最弱，通常依赖 `rand`/chrono 第三方 crate。
