import itertools


class Countdown:
    # 自定义迭代器：实现 __iter__ 与 __next__
    def __init__(self, start: int) -> None:
        self.cur = start

    def __iter__(self) -> "Countdown":
        return self

    def __next__(self) -> int:
        if self.cur <= 0:
            raise StopIteration
        self.cur -= 1
        return self.cur + 1


def squares(n: int):
    # 生成器函数：yield 逐项产出，惰性
    for i in range(1, n + 1):
        yield i * i


def flatten(nested):
    # yield from 把子迭代器「让渡」出去
    for sub in nested:
        yield from sub


def main() -> None:
    print("1) 迭代协议：iter() / next() / 自定义 __iter__+__next__")
    cd = Countdown(3)
    print(f"  iter(cd) is cd={iter(cd) is cd}；list(Countdown(3))={list(Countdown(3))}")

    print("2) 生成器惰性：调用生成器函数不会执行函数体")
    traced: list[int] = []

    def watch():
        for i in range(3):
            traced.append(i)
            yield i

    g = watch()
    print(f"  刚创建时 traced={traced}（未执行）；next 一次得 {next(g)}，traced={traced}")

    print("3) 一次性消费陷阱：生成器耗尽后为空")
    gen = squares(4)
    print(f"  首次 list={list(gen)}；再次 list={list(gen)}（已耗尽）")

    print("4) itertools：chain / takewhile / map-filter 惰性管道")
    chained = list(itertools.chain([1, 2], [3, 4]))
    pipeline = (x * x for x in range(10) if x % 2 == 0)
    print(f"  chain={chained}；惰性平方管道={list(pipeline)}")
    print(f"  takewhile(<5)={list(itertools.takewhile(lambda x: x < 5, range(10)))}")

    print("5) 无限序列取前 N 项：itertools.count + islice")
    evens = itertools.islice((x for x in itertools.count() if x % 2 == 0), 6)
    print(f"  前 6 个偶数={list(evens)}")

    print("6) yield from 与 iter(callable, sentinel)")
    print(f"  flatten={list(flatten([[1, 2], [3], [4, 5]]))}")
    src = iter([9, 8, 7, 0, 6])
    print(f"  iter(遇到 0 停)={list(iter(lambda: next(src), 0))}")


if __name__ == "__main__":
    main()
