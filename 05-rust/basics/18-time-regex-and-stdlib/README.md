# 18 · Rust · 时间、正则与常用标准库

## 这个例子讲什么

标准库里最常用的三块：时间（`Instant` 单调时钟 / `Duration` 时长 / `SystemTime` 墙上时间）、排序（`sort` / `sort_by` / `sort_by_key`）、以及**标准库没有的东西**——正则和随机数，用最小手写实现替代，并说明生产上该用什么 crate。

## 怎么跑

```
cd backend-roadmap/05-rust/basics/18-time-regex-and-stdlib
cargo run --quiet
```

## 关键点

- 测耗时用 `Instant`，不要用 `SystemTime`：后者会被 NTP 校时和用户改表影响，可能倒退。
- `Duration` 支持 `from_millis` / `as_secs` / `as_millis`，跨语言对照时注意「秒 vs 毫秒」的单位混用。
- **Rust 标准库没有正则**，生产用 `regex` crate；本例用 `starts_with` + `chars().skip().take()` 演示最小替代。
- **标准库也没有随机数**，生产用 `rand` crate；本例用固定种子 LCG 保证输出可复现，这也是「示例代码要稳定」的正确做法。
- `sort_by_key` 适合「按某个可比较字段排序」，比手写 `sort_by` 清晰。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 单调计时 | `time.Since` | `performance.now()` | `time.perf_counter` | `System.nanoTime` | `Instant::now` |
| 时长 | `time.Duration` | 毫秒数字 | `timedelta` | `Duration` | `Duration` |
| 正则 | `regexp` 标准库 | `RegExp` 内置 | `re` 标准库 | `java.util.regex` 标准库 | **无，需 `regex` crate** |
| 随机数 | `math/rand` | `Math.random` | `random` | `Random` | **无，需 `rand` crate** |

Go / Python / Java / JS 的正则和随机数都在标准库或语言内置，Rust 走「核心库极小 + crates.io 生态」路线，所以标准库只保留最底层能力。
