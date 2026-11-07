import sys


def bump(label: str, out: list[str]) -> bool:
    # 有副作用，用来观察短路求值是否跳过右侧
    out.append(label)
    return True


def main() -> None:
    print("1) 算术：+ - * / // % ** 与 floor 取整")
    print(f"  7/2={7/2}；7//2={7//2}；-7//2={-7//2}（向下取整）；7%(-3)={7 % -3}")
    print(f"  2**10={2**10}；divmod(7,2)={divmod(7,2)}")

    print("2) 比较、链式比较、逻辑运算")
    print(f"  1<2<3={1 < 2 < 3}；'a'<'b'={'a' < 'b'}；[1,2]==[1,2]={[1,2]==[1,2]}")
    print(f"  and/or 返回操作数: 0 or 'x'={0 or 'x'!r}；1 and 2={1 and 2}")

    print("3) 短路求值：and 左侧为假 / or 左侧为真则右侧不求值")
    trace: list[str] = []
    r1 = bump("left", trace) and bump("right", trace)
    print(f"  True and True 调用序列={trace}，结果={r1}")
    trace.clear()
    r2 = (0 and bump("skipped", trace)) or bump("used", trace)
    print(f"  短路后只执行 or 右侧: trace={trace}，r2={r2}")

    print("4) 位运算：& | ^ ~ << >>")
    print(f"  1010&1100={0b1010 & 0b1100:04b}；1010|0101={0b1010|0b0101:04b}；1010^1100={0b1010 ^ 0b1100:04b}")
    print(f"  ~5={~5}；1<<4={1<<4}；16>>2={16>>2}")

    print("5) 赋值、复合赋值与优先级")
    x = 2
    x += 3
    x **= 2
    print(f"  x=2; x+=3; x**=2 -> {x}；2+3*4={2+3*4}（先乘后加）")
    items = [1]
    alias = items
    alias += [2]  # list 的 += 是原地扩展，别名同步变化
    print(f"  list += 原地扩展: items={items}；alias is items={alias is items}")

    print("6) 整数溢出与类型转换")
    huge = 2 ** 200
    print(f"  int 任意精度，2**200 位数={huge.bit_length()}（无回绕）；sys.maxsize={sys.maxsize}")
    print(f"  float 溢出: 1e308*1e10={1e308*1e10}")
    print(f"  隐式提升 1+0.5={1+0.5}（{type(1+0.5).__name__}）；int('42')={int('42')}；bool('')={bool('')}")
    print(f"  nan != nan={float('nan') != float('nan')}")

    print("7) 三元表达式与海象运算符")
    score = 75
    print(f"  条件表达式: {'及格' if score >= 60 else '不及格'}")
    if (n := len([1, 2, 3])) > 2:
        print(f"  := 赋值并判断: n={n}")


if __name__ == "__main__":
    main()
