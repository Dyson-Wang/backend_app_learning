def classify(n: int) -> str:
    # if/elif/else 配合链式比较：0 < n < 10 等价于 0 < n and n < 10
    if 0 < n < 10:
        return "小"
    elif n < 0:
        return "负"
    else:
        return "大"


def main() -> None:
    print("1) 分支：if/elif/else 与链式比较、条件表达式（三元）")
    for n in (-3, 5, 42):
        print(f"  {n:>3} -> {classify(n)}；三元: {'偶数' if n % 2 == 0 else '奇数'}")

    print("2) for 循环与 range / enumerate / zip / 解包")
    total = 0
    for i in range(1, 6):
        total += i
    print(f"  range(1,6) 求和={total}")
    rows = [("a", 1), ("b", 2)]
    print(f"  enumerate: {[(i, k) for i, (k, _) in enumerate(rows)]}")
    print(f"  zip: {[f'{k}={v}' for k, v in zip(['x', 'y'], [10, 20])]}")

    print("3) for...else：循环未被 break 才执行 else（Python 独有）")
    for prime in (2, 3, 4, 5, 7):
        for d in range(2, prime):
            if prime % d == 0:
                print(f"  {prime} 不是质数，被 {d} 整除")
                break
        else:
            print(f"  {prime} 是质数（内层 for 的 else 触发）")

    print("4) while / break / continue / 海象运算符 :=")
    buf: list[int] = []
    while (v := len(buf)) < 5:
        buf.append(v * v)
        if v == 2:
            continue  # 演示 continue，逻辑照常推进
    print(f"  while + := 结果={buf}")

    print("5) any/all 与生成器式短路")
    nums = [1, 3, 5, 8]
    print(f"  any(偶数)={any(n % 2 == 0 for n in nums)}；all(奇数)={all(n % 2 for n in nums)}")


if __name__ == "__main__":
    main()
