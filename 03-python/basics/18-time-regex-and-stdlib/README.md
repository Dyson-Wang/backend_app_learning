# 18 · Python · 时间、正则与常用标准库

## 这个例子讲什么

`datetime` / `timedelta` / `timezone` 的时间点、时长、时区转换与格式化解析；`re` 的编译、命名捕获、查找、替换、切分；固定种子 `random.Random(42)`；`sorted` / `Counter` / `statistics` / `functools.reduce` / `math` 等常用工具。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/18-time-regex-and-stdlib/main.py
```

## 关键点

- **时间用固定值演示格式化 / 差值**，不打印当前绝对时间作为结论；`astimezone()` 做时区转换，`strftime`/`fromisoformat`/`strptime` 负责格式化与解析。
- `naive`（无 tzinfo）与 `aware`（带时区）`datetime` 不能直接相减；生产里建议统一存 UTC，展示时再转本地。
- `re.compile` 复用正则；`(?P<name>...)` 命名捕获组、`findall` / `sub` / `split` / `finditer` 覆盖常见操作。
- `random.Random(42)` 用**独立实例 + 固定种子**保证输出可复现（不污染全局 `random`）；安全场景用 `secrets`。
- `sorted(key=)` 稳定、`Counter.most_common` 取 TopN、`statistics.mean` / `math.isqrt` / `pow(a,b,m)` 都是标准库现成工具。

## 和其他语言的差异

- **Go**：`time.Time` / `time.Duration` / `time.Parse`、布局字符串 `'2006-01-02'` 很特别、`regexp`、`math/rand`。
- **Node.js/TS**：`Date`、`RegExp`、`Math.random`、`Intl.DateTimeFormat`。
- **Java**：`java.time`（`LocalDateTime` / `ZonedDateTime` / `Duration`）、`java.util.regex.Pattern`、`Random`。
- **Rust**：`std::time`（`Instant`/`Duration`），时区/正则/随机主要靠第三方 `chrono` / `regex` / `rand`。
- **Python**：`datetime`/`timedelta`/`timezone`/`zoneinfo`、`re`、`random.Random(seed)`、`sorted`/`Counter`/`statistics`/`itertools` 内置电池。
