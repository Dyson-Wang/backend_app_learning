class ValidationError(Exception):
    # 自定义异常：带一个业务错误码
    def __init__(self, code: str, message: str) -> None:
        super().__init__(message)
        self.code = code


def power(base: float, exp: int = 2, /, *, strict: bool = False, scale: float = 1.0) -> float:
    # / 之前是位置限定参数；* 之后是关键字限定参数
    if strict and exp < 0:
        raise ValidationError("NEG_EXP", "strict 模式不允许负指数")
    return round(base ** exp * scale, 4)


def sum_all(*args: int) -> int:
    """把任意多个整数求和。"""
    return sum(args)


def tags(**opts: object) -> str:
    # **opts 把多余关键字收成 dict
    return ",".join(sorted(opts.keys()))


def append_bad(x: int, box: list[int] = []) -> list[int]:  # 故意用可变默认值演示坑
    box.append(x)
    return box


def append_good(x: int, box: list[int] | None = None) -> list[int]:
    box = [] if box is None else box
    box.append(x)
    return box


def main() -> None:
    print("1) 参数形态：默认值 / 位置限定 / 关键字限定 / *args / **kwargs")
    print(f"  power(3)={power(3)}；power(2, 10)={power(2, 10)}")
    print(f"  power(2, 3, scale=0.5)={power(2, 3, scale=0.5)}")
    print(f"  sum_all(1,2,3,4)={sum_all(1,2,3,4)}；tags(b=2,a=1)={tags(b=2,a=1)}")
    print(f"  power 关键字默认值={power.__kwdefaults__}")

    print("2) 可变默认参数陷阱：默认值只在定义时求值一次")
    print(f"  第一次={append_bad(1)}；第二次={append_bad(2)}（残留被污染）")
    print(f"  修正后第一次={append_good(1)}；第二次={append_good(2)}")

    print("3) try / except / else / finally")
    for raw in ("16", "oops", "-2"):
        try:
            n = int(raw)
            if n < 0:
                raise ValidationError("NEG", "不能为负")
        except ValueError:
            print(f"  {raw!r}: 捕获 ValueError（int 转换失败）")
        except ValidationError as exc:
            print(f"  {raw!r}: 捕获业务异常 code={exc.code}")
        else:
            print(f"  {raw!r}: 无异常走 else，n={n}")
        finally:
            print(f"  {raw!r}: finally 一定执行")

    print("4) 异常链 raise ... from")
    try:
        try:
            int("bad")
        except ValueError as exc:
            raise ValidationError("PARSE", "包装一层") from exc
    except ValidationError as exc:
        print(f"  外层={exc.code}；__cause__ 类型={type(exc.__cause__).__name__}")

    print("5) assert 与 lambda")
    assert sum_all(1, 2, 3) == 6
    square = lambda x: x * x
    print(f"  assert 通过；lambda(7)={square(7)}；sum_all 有文档串={sum_all.__doc__ is not None}")


if __name__ == "__main__":
    main()
