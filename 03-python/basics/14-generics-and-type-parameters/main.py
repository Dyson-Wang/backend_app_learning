from typing import Generic, Protocol, TypeVar, get_type_hints, runtime_checkable

# 旧写法：TypeVar + Generic
T = TypeVar("T")


class Stack(Generic[T]):
    def __init__(self) -> None:
        self._items: list[T] = []

    def push(self, item: T) -> None:
        self._items.append(item)

    def pop(self) -> T:
        return self._items.pop()


# 约束：用 Protocol 表达「必须实现 __len__」的结构化 bound
@runtime_checkable
class HasLen(Protocol):
    def __len__(self) -> int: ...


# 3.12+ PEP 695 新语法：def f[T](...)、class Box[T]、以及内联 bound T: HasLen
def first[T](items: list[T]) -> T:
    return items[0]


def length_of[T: HasLen](obj: T) -> int:
    return len(obj)


class Box[T]:
    def __init__(self, value: T) -> None:
        self.value = value

    def get(self) -> T:
        return self.value


def main() -> None:
    print("1) 旧写法 TypeVar + Generic：泛型集合 Stack[T]")
    s: Stack[int] = Stack()
    s.push(10)
    s.push(20)
    print(f"  Stack[int] pop={s.pop()}；Stack[int].__args__={Stack[int].__args__}")

    print("2) 3.12 新语法 def first[T](...) 的泛型函数")
    print(f"  first([1,2,3])={first([1,2,3])}；first(['a','b'])={first(['a','b'])}")
    print(f"  get_type_hints(first) 键={sorted(get_type_hints(first).keys())}")

    print("3) 3.12 新语法 class Box[T] 的泛型类型")
    b = Box[str]("hello")
    print(f"  Box[str].get()={b.get()!r}；Box[int](7).get()={Box[int](7).get()}")

    print("4) 类型擦除：运行时只剩原始类，泛型参数不进实例")
    print(f"  type(Box[int](1)) is Box -> {type(Box[int](1)) is Box}")
    print(f"  list[int].__args__={list[int].__args__}；dict[str,int].__args__={dict[str,int].__args__}")

    print("5) 约束 / 结构化 Protocol：T 必须有 __len__")
    print(f"  length_of([1,2,3])={length_of([1,2,3])}；length_of('abc')={length_of('abc')}")
    print(f"  isinstance([1], HasLen)={isinstance([1], HasLen)}（runtime_checkable 才可用 isinstance）")

    print("6) 单态化 vs 类型擦除：Python 走擦除，注解不产生多份机器码")
    print(f"  Box[int] 不是类而是泛型别名={not isinstance(Box[int], type)}")
    print(f"  TypeVar 有约束写法: TypeVar('N', int, float) -> {TypeVar('N', int, float).__constraints__}")


if __name__ == "__main__":
    main()
