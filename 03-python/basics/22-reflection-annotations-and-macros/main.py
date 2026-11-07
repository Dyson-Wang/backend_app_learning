import functools
import inspect
from dataclasses import dataclass, fields
from typing import get_type_hints


@dataclass
class User:
    uid: int
    name: str = "anon"


REGISTRY: dict[str, type] = {}


class Meta(type):
    # 元类：拦截类创建，自动登记
    def __new__(mcls, name: str, bases: tuple, ns: dict) -> type:
        cls = super().__new__(mcls, name, bases, ns)
        REGISTRY[name] = cls
        return cls


class Service(metaclass=Meta):
    def __init__(self, tag: str) -> None:
        self.tag = tag

    def run(self, x: int) -> int:
        return x * 2


def logged(func):
    # 装饰器：运行时替换函数对象（Python 没有宏，这是最接近的元编程）
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        print(f"    [decorator] 调用 {func.__name__}")
        return func(*args, **kwargs)
    return wrapper


@logged
def add(a: int, b: int) -> int:
    return a + b


def main() -> None:
    print("1) 运行时类型信息：type / __annotations__ / get_type_hints")
    print(f"  User.__name__={User.__name__}；User.__annotations__={User.__annotations__}")
    hints = {k: getattr(v, "__name__", str(v)) for k, v in get_type_hints(add).items()}
    print(f"  get_type_hints(add)={hints}")
    print(f"  isinstance(User(1), object)={isinstance(User(1), object)}")

    print("2) 按名字访问字段/方法：getattr / setattr / hasattr / callable")
    u = User(7)
    print(f"  getattr(u,'name')={getattr(u, 'name')!r}；hasattr(u,'uid')={hasattr(u, 'uid')}")
    setattr(u, "name", "alice")
    print(f"  setattr 后 name={u.name!r}；__init__ 可调用={callable(getattr(u, '__init__'))}")

    print("3) dataclasses.fields 与 inspect.signature")
    print(f"  fields(User)={[(f.name, getattr(f.type, '__name__', str(f.type))) for f in fields(User)]}")
    sig = inspect.signature(add)
    print(f"  signature(add)={sig}；参数名={list(sig.parameters)}")

    print("4) 装饰器：运行时替换 + functools.wraps 保留元信息")
    print(f"  add(2,3)={add(2,3)}；装饰后 __name__={add.__name__}")

    print("5) 元类：拦截类创建并注册")
    print(f"  type(Service).__name__={type(Service).__name__}；已注册={list(REGISTRY)}")
    print(f"  Service('x').run(5)={Service('x').run(5)}")

    print("6) 动态建类、动态调用与边界")
    Dynamic = type("Dynamic", (), {"value": 99})
    print(f"  type() 动态建类: Dynamic().value={Dynamic().value}")
    method = getattr(Service('y'), "run")
    print(f"  getattr 取方法再调用: {method(10)}")
    print("  反射代价：属性查找走字典，比直接调用慢；只能拿到运行期真正存在的名字")


if __name__ == "__main__":
    main()
