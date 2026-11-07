import functools
from typing import Callable


def make_adder(n: int) -> Callable[[int], int]:
    # 闭包：内层函数捕获外层变量 n
    def add(x: int) -> int:
        return x + n
    return add


def make_counter() -> Callable[[], int]:
    count = 0

    def inc() -> int:
        nonlocal count  # 修改捕获的可变状态必须声明 nonlocal
        count += 1
        return count
    return inc


def apply_twice(f: Callable[[int], int], x: int) -> int:
    # 高阶函数：函数作为值传入
    return f(f(x))


def log_calls(f: Callable[..., int]) -> Callable[..., int]:
    # 返回函数 + functools.wraps 保留元信息
    @functools.wraps(f)
    def wrapper(*args: object) -> int:
        print(f"    [decorator] 调用 {f.__name__}{args}")
        return f(*args)  # type: ignore[arg-type]
    return wrapper


def main() -> None:
    print("1) 闭包捕获与 __closure__ 单元")
    add5 = make_adder(5)
    print(f"  add5(10)={add5(10)}；cell 数={len(add5.__closure__ or ())}")
    print(f"  捕获的 n={add5.__closure__[0].cell_contents}（生命周期随闭包延长）")

    print("2) nonlocal 修改捕获的可变状态")
    c = make_counter()
    print(f"  连续调用: {c()}, {c()}, {c()}（状态被闭包持有）")

    print("3) 迟绑定坑：循环里建闭包全指向最后的变量")
    funcs = [lambda: i for i in range(3)]
    print(f"  [lambda: i ...] 结果={[f() for f in funcs]}（都是 2）")
    fixed = [lambda i=i: i for i in range(3)]  # 默认参数当场绑定
    print(f"  默认参数修复={[f() for f in fixed]}")

    print("4) 函数作为值 / 返回函数 / 常见高阶函数")
    print(f"  apply_twice(lambda x: x*2, 3)={apply_twice(lambda x: x * 2, 3)}")
    print(f"  map={list(map(lambda x: x*x, range(5)))}；filter={list(filter(lambda x: x%2==0, range(6)))}")
    print(f"  functools.reduce 求和={functools.reduce(lambda a, b: a + b, range(5))}")
    print(f"  sorted(key=len)={sorted(['aaa','b','cc'], key=len)}")

    print("5) 装饰器、partial 与绑定方法自带 self")
    @log_calls
    def triple(x: int) -> int:
        return x * 3
    print(f"  triple(4)={triple(4)}；装饰后名字={triple.__name__}")
    pow2 = functools.partial(pow, 2)  # 固定第一个参数
    print(f"  partial(pow,2)(10)={pow2(10)}")

    class Robot:
        def __init__(self, name: str) -> None:
            self.name = name
        def hi(self) -> str:
            return f"hi {self.name}"
    cb = Robot("py").hi
    print(f"  绑定方法自带 self: {cb()}；cb.__self__ 存在={cb.__self__ is not None}")


if __name__ == "__main__":
    main()
