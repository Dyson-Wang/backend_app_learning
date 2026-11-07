# 05 · 集合与迭代器

## 这个例子讲什么

数组是值、切片是描述符、map 是引用类型。重点看 slice 的底层数组共享——生产里最常见的数据串改 bug。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/05-collections-and-iterators
```

## 关键点

- `slice` 头部只有三字段（指针 / len / cap），**赋值和传参都复制描述符，底层数组仍共享**。
- `append` 在 cap 足够时会原地写入，因此 `b := append(a, x)` 可能改变 `base`；需要隔离时用 `copy` 或三索引切片 `a[:2:2]`。
- `map` 取值用 comma-ok 区分「不存在」和「存在但为零值」；map 不是并发安全的，多 goroutine 读写必须加锁。
- `range` 会复制值，改元素要用下标 `s[i] = ...`。
- Go 没有内置 `map/filter/reduce`；1.18+ 可以自己写泛型函数，标准库的 `slices` / `maps` 包（1.21+）也提供了一批。

## 和其他语言的差异

| 语言 | 动态数组 | 映射 | 迭代 |
|---|---|---|---|
| Go | `[]T` slice（共享底层数组） | `map[K]V` | `range` |
| Node.js | `Array` | `Map` / 普通对象 | `for-of` / `.map()` |
| Python | `list` | `dict` | `for` + 推导式 |
| Java | `ArrayList` | `HashMap` | for-each / Stream |
| Rust | `Vec<T>`（所有权移动） | `HashMap` | `Iterator` + 适配器 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
