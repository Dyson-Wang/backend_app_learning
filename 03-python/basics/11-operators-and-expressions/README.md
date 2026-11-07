# 11 · Python · 运算符与表达式

## 这个例子讲什么

算术 / 比较 / 逻辑 / 位运算、短路求值（用带副作用的函数实证）、优先级、赋值与复合赋值、整数不溢出与 `float` 的 `inf`、隐式类型转换、条件表达式与海象运算符。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/11-operators-and-expressions/main.py
```

## 关键点

- `//` 是 **floor 向下取整**（`-7//2 == -4`），`%` 结果与除数同号，和 C/Java 的「向零取整」不同。
- `**` 幂、`divmod(a,b)` 一次拿商和余；`2**200` 也不溢出，因为 `int` 任意精度。
- **短路求值**：`and` 左假、`or` 左真就不求值右侧；`and/or` 返回**操作数本身**而非 `bool`，所以 `0 or 'x'` 得到 `'x'`。
- 链式比较 `1<2<3` 只求值一次中间量；`list += [x]` 是**原地扩展**（`__iadd__`），别名会跟着变。
- `float` 才会溢出成 `inf`；`int + float` 隐式提升为 `float`；`bool` 是 `int` 子类。

## 和其他语言的差异

- **Go**：无三元、`&^` 位清除、定长整型**溢出回绕**、`&&`/`||` 短路、`math/bits` 处理进位。
- **Node.js/TS**：`===`/`==`、`??`、`?.`、位运算按 32 位有符号、`BigInt` 才任意精度。
- **Java**：`>>>` 无符号右移、整数溢出回绕、`Math.addExact` 显式检测、`instanceof` 模式匹配。
- **Rust**：`checked_add` / `wrapping_add` / `saturating_add` 显式溢出语义、无 `++`、`as` 显式转换。
- **Python**：`int` 任意精度不溢出、`//` floor 整除、`**`、链式比较、海象 `:=`、`A if c else B` 三元。
