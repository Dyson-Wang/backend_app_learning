from dataclasses import dataclass
from enum import Enum, IntEnum, auto
from typing import Union, assert_never


class Color(Enum):
    # 枚举：成员是单例，有 .name / .value
    RED = 1
    GREEN = 2
    BLUE = auto()


class HttpStatus(IntEnum):
    # IntEnum 成员可直接当 int 用
    OK = 200
    NOT_FOUND = 404


@dataclass
class Circle:
    r: float


@dataclass
class Rect:
    w: float
    h: float


Shape = Union[Circle, Rect]  # 用 typing.Union 表达 ADT（3.10+ 也可写 Circle | Rect）


def area(shape: Shape) -> float:
    # 类模式匹配：按类型 + 关键字属性解构
    match shape:
        case Circle(r=r):
            return 3.14 * r * r
        case Rect(w=w, h=h):
            return w * h
        case _:
            raise AssertionError("未知形状")


def describe(value: object) -> str:
    # 字面量 / 捕获 / 通配 / 守卫 / OR / 序列 / 映射 / 类模式
    match value:
        case 0:
            return "零"
        case int(n) if n > 0:
            return f"正整数 {n}"
        case int(n):
            return f"非正整数 {n}"
        case [first, *rest] if rest:
            return f"序列 首={first} 其余={len(rest)}"
        case {"kind": "point", "x": x, "y": y}:
            return f"点 ({x},{y})"
        case Color.RED | Color.GREEN:
            return "暖色或绿色枚举"
        case str(s):
            return f"字符串 {s!r}"
        case _:
            return "其它"


def main() -> None:
    print("1) enum.Enum / IntEnum：取值、遍历、比较")
    print(f"  Color.RED={Color.RED} name={Color.RED.name} value={Color.RED.value}")
    print(f"  成员列表={[c.name for c in Color]}；Color(2)={Color(2)}")
    print(f"  IntEnum 当 int: HttpStatus.OK == 200 -> {HttpStatus.OK == 200}；+1={HttpStatus.OK + 1}")

    print("2) dataclass + Union 当 ADT，match 类模式解构")
    shapes: list[Shape] = [Circle(2.0), Rect(3.0, 4.0)]
    for sh in shapes:
        print(f"  {sh.__class__.__name__} 面积={area(sh)}")

    print("3) match case 其它模式：守卫 / OR / 序列 / 映射 / 类")
    cases: list[object] = [0, 5, -3, [1, 2, 3], {"kind": "point", "x": 1, "y": 2}, Color.RED, "hi", None]
    for c in cases:
        print(f"  {c!r:>34} -> {describe(c)}")

    print("4) 穷尽性：Python 无编译期检查，靠 case _ 兜底")
    print(f"  未覆盖对象走通配={describe(object())!r}")
    print(f"  assert_never 可用（供 mypy 静态检查）={callable(assert_never)}")


if __name__ == "__main__":
    main()
