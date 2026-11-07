from typing import Any, Final

# 常量在 Python 里只是约定：全大写 + Final 注解，运行时并不禁止改写
MAX_RETRY: Final[int] = 3


def show(value: Any) -> None:
    # 用 type() 与 bool() 观察动态类型；类型注解只给静态检查器看
    print(f"  value={value!r:<12} type={type(value).__name__:<8} 真值={bool(value)}")


def main() -> None:
    print("1) 动态类型：同一个名字可反复绑定到不同类型")
    x: Any = 42
    show(x)
    x = "四十二"
    show(x)
    x = [1, 2, 3]
    show(x)

    print("2) 零值：没有语言级零值，空值统一是 None")
    nothing = None
    print(f"  nothing is None: {nothing is None}")
    print(f"  假值: bool(None)={bool(None)} bool(0)={bool(0)} bool('')={bool('')} bool([])={bool([])}")

    print("3) 整数任意精度，不会像定长整数那样溢出")
    big = 2 ** 80
    print(f"  2**80 = {big} (位数={big.bit_length()})")
    print(f"  int 无溢出: {big + 1 > big}")

    print("4) 常量约定与类型注解")
    print(f"  MAX_RETRY={MAX_RETRY}；函数注解={show.__annotations__}")
    n: int = 1
    n = "其实我变成了字符串"  # 注解不阻止赋值
    print(f"  注解不强制: n 现在是 {type(n).__name__}")

    print("5) bool 是 int 子类；链式比较；序列解包")
    print(f"  issubclass(bool, int)={issubclass(bool, int)}；True+True={True + True}")
    a, b = 1, 2
    a, b = b, a
    print(f"  交换赋值后 a={a} b={b}")

    print("6) 身份：is 判断是不是同一个对象")
    obj = object()
    alias = obj
    print(f"  alias is obj: {alias is obj}")
    print(f"  [] is []: {[] is []}（每次新建）；[] == []: {[] == []}（值相等）")


if __name__ == "__main__":
    main()
