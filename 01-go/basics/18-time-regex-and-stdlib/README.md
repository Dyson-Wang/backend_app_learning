# 18 · 时间、正则与常用标准库

## 这个例子讲什么

`time` 的时间点/时长/时区/格式化解析、`regexp` 的匹配捕获替换、固定种子的随机数、`sort` 排序，以及 `strings` 常用工具函数。时间全部用**固定值**，只演示格式化与差值。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/18-time-regex-and-stdlib
```

## 关键点

- **魔法布局时间**：Go 用参考时间 `2006-01-02 15:04:05`（记忆口诀 1-2-3-4-5）作为 `Format`/`Parse` 的 layout，而不是 `yyyy-MM-dd` 这种占位符。
- **时间点带位置**：`time.Date(...)` 构造固定时刻；`time.FixedZone` 无需 tzdata 即可做时区演示（比 `LoadLocation` 更可控）；`t.In(loc)` 只是换视角，绝对时刻不变。
- **`Duration` 是 `int64` 纳秒**：`t2.Sub(t1)` 得到时长，`Hours()/Minutes()/Round()/String()` 做换算。本例只打印差值，不打印当前绝对时间，输出稳定可复现。
- **正则**：`regexp.MustCompile`（非法 pattern 直接 panic）、`FindStringSubmatch` 取捕获组、`(?P<name>...)` 命名组配合 `SubexpIndex`、`ReplaceAllString`、`FindAllString`。Go 正则用 RE2 语法，**无回溯、无反向引用**，但保证线性时间。
- **随机数**：`math/rand` 是伪随机，`rand.New(rand.NewSource(7))` 固定种子可复现；`crypto/rand` 才是密码学安全随机。
- **排序**：`sort.Ints/Strings` 就地排；`sort.Slice` 传 less 闭包；`sort.SearchInts` 在有序切片上二分。

## 和其他语言的差异

| 语言 | 时间 | 正则 | 随机 / 排序 |
|---|---|---|---|
| Go | `time.Time`/`Duration`/`time.FixedZone`，魔法 layout | `regexp`（RE2，无回溯） | `math/rand`（固定种子）、`sort` 包 |
| Node.js | `Date` / 新标准 `Temporal` | `RegExp`（有回溯/反向引用） | `Math.random()`（不可设种子）、`Array.sort` |
| Python | `datetime` / `zoneinfo` / `timedelta` | `re`（有回溯/反向引用） | `random.seed`、`sorted` / `list.sort` |
| Java | `java.time`（`LocalDateTime`/`ZonedDateTime`/`Duration`） | `java.util.regex.Pattern`（有回溯） | `java.util.Random`（可设种子）、`Collections.sort` |
| Rust | `std::time`（`Instant`/`Duration`）；时区/格式化需第三方 `chrono` | 标准库**无正则**，用第三方 `regex` crate（基于 RE2 思想） | `rand` 是第三方 crate；`sort`/`sort_by` 在标准库 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
