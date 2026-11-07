from collections import Counter, defaultdict, deque


def main() -> None:
    print("1) 四种内置集合：list / tuple / set / dict")
    lst = [3, 1, 2, 1]
    tpl = (3, 1, 2)  # 不可变，可作 dict 键
    st = set(lst)
    d = {"a": 1, "b": 2}
    print(f"  list={lst} tuple={tpl} set(去重)={sorted(st)} dict键={sorted(d)}")

    print("2) 推导式：list / dict / set 三种")
    squares = {n: n * n for n in range(1, 6)}
    evens = [n for n in range(10) if n % 2 == 0]
    uniq_mod = {n % 3 for n in range(10)}
    print(f"  dict推导={squares}；list推导={evens}；set推导={sorted(uniq_mod)}")

    print("3) 切片、排序、反转")
    seq = list(range(10))
    print(f"  seq[2:8:2]={seq[2:8:2]}；seq[::-1][:3]={seq[::-1][:3]}")
    words = ["banana", "fig", "apple"]
    print(f"  sorted(key=len)={sorted(words, key=len)}（稳定排序）")

    print("4) collections：deque / Counter / defaultdict")
    dq = deque([1, 2, 3])
    dq.appendleft(0)
    dq.append(4)
    print(f"  deque={list(dq)}；Counter={dict(Counter('abracadabra'))}")
    dd: defaultdict[str, list[int]] = defaultdict(list)
    for k, v in [("x", 1), ("y", 2), ("x", 3)]:
        dd[k].append(v)
    print(f"  defaultdict={dict(dd)}")

    print("5) 迭代器协议：iter() / next() / 一次性消费")
    it = iter(["a", "b", "c"])
    print(f"  next 连续调用={next(it)},{next(it)},{next(it)}")
    try:
        next(it)
    except StopIteration:
        print("  再 next 抛 StopIteration（for 循环靠它结束）")
    gen = (n * n for n in range(3))
    print(f"  生成器首次 list={list(gen)}；再次 list={list(gen)}（已耗尽）")

    print("6) 解包与 dict 合并")
    first, *middle, last = range(6)
    print(f"  星号解包: first={first} middle={middle} last={last}")
    merged: dict[str, int] = {**{"a": 1}, "b": 2}
    print(f"  dict 合并={merged}")


if __name__ == "__main__":
    main()
