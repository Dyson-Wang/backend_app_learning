# 18 · Java · 时间、正则与常用标准库

## 这个例子讲什么

`java.time` 时间点 / 时长 / 时区、格式化与解析、正则匹配捕获替换、固定种子随机数、排序与查找、Base64 / Math 等常用工具。

## 怎么跑

```
cd backend-roadmap/04-java/basics/18-time-regex-and-stdlib
javac Main.java && java Main

# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **`java.time`（Java 8+）不可变且线程安全**：`LocalDateTime` / `Duration`（秒纳秒级）/ `Period`（年月日）/ `ZonedDateTime` + `ZoneId`；`DateTimeFormatter` 负责格式化与解析。
- 全部用**固定时间点**演示，不打印当前墙上时钟；时区转换用固定时刻（`atZone` + `withZoneSameInstant`），无 DST 干扰。
- **正则** `Pattern` / `Matcher`：`find()` + `group(n)` 分组捕获，`replaceAll("$1-...")` 用组引用替换，`Pattern.matches` 整体匹配。
- **随机数**用 `new Random(42L)` 固定种子；示例只比较「两个同种子序列是否相等」这一稳定结论，不依赖具体数值。
- `Collections.sort/max/binarySearch/frequency`、`Arrays`、`Base64`、`Math` 都在 JDK 标准库内，零第三方依赖。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 时间 | `java.time`（内置） | `time.Time`/`Duration` | `Date`（Intl/Temporal 提案） | `datetime`/`zoneinfo` | `std::time`（chrono 是第三方） |
| 正则 | `java.util.regex`（内置） | `regexp`（内置） | `RegExp`（内置） | `re`（内置） | `regex` 是第三方 crate |
| 随机数 | `java.util.Random` 固定种子 | `math/rand` | `Math.random`（不可控） | `random.seed` | `rand` 是第三方 crate |
| 排序 | `Collections.sort`/`List.sort` | `sort.Slice` | `Array.prototype.sort` | `sorted`/`list.sort` | `slice.sort` |
| 标准库覆盖 | 全部内置 | 大部分内置 | 部分需提案 | 大部分内置 | 很多靠 crate |
