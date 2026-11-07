import abc
import json
from dataclasses import dataclass, field
from typing import Protocol


@dataclass
class User:
    # dataclass 自动生成 __init__ / __repr__ / __eq__
    uid: int
    name: str
    tags: list[str] = field(default_factory=list)

    @property
    def label(self) -> str:
        # 用计算属性代替 getter
        return f"[{self.uid}] {self.name}"

    def __len__(self) -> int:
        # 鸭子类型：实现了 __len__ 就能被 len() 调用
        return len(self.tags)


class Repository(abc.ABC):
    # 抽象基类 = 名义上的「接口」，靠继承实现
    @abc.abstractmethod
    def get(self, uid: int) -> str: ...

    def exists(self, uid: int) -> bool:
        return bool(self.get(uid))


class InMemoryRepo(Repository):
    def __init__(self) -> None:
        self._data = {1: "alice", 2: "bob"}

    def get(self, uid: int) -> str:
        return self._data.get(uid, "")


class Greeter(Protocol):
    # 结构化接口：不继承也行，只要有 greet 方法
    def greet(self) -> str: ...


class Robot:
    def greet(self) -> str:
        return "beep"


def welcome(g: Greeter) -> str:
    return f"欢迎：{g.greet()}"


class Tracked:
    # 类属性被所有实例共享；实例属性各自独立
    created = 0

    def __init__(self, name: str) -> None:
        Tracked.created += 1
        self.name = name

    def __repr__(self) -> str:
        return f"Tracked({self.name!r})"


def main() -> None:
    print("1) dataclass + property + 鸭子类型方法")
    u = User(1, "alice", ["admin", "dev"])
    print(f"  dataclass: {u}；label={u.label}；len(u)={len(u)}")
    u2 = User(1, "alice", ["admin", "dev"])
    print(f"  __eq__ 按值比较: u == u2 -> {u == u2}")

    print("2) abc.ABC 抽象基类作为名义接口")
    try:
        Repository()  # 带未实现抽象方法，不能实例化
    except TypeError:
        print("  实例化抽象类触发 TypeError")
    repo = InMemoryRepo()
    print(f"  repo.get(1)={repo.get(1)!r}；exists(9)={repo.exists(9)}")

    print("3) Protocol 结构化接口：无需继承")
    print(f"  Robot 没继承 Greeter，仍可: {welcome(Robot())}")

    print("4) 类属性共享、实例属性独立")
    c1, c2 = Tracked("a"), Tracked("b")
    print(f"  类属性 Tracked.created={Tracked.created}；c1 实例字典={c1.__dict__}")
    print(f"  c1 is c2: {c1 is c2}；名字不同: {c1.name != c2.name}")

    print("5) 模块：每个 .py 文件即模块，__name__ / __module__")
    print(f"  本模块名={__name__}；User.__module__={User.__module__}")
    print(f"  json 是模块对象={type(json).__name__}；User 注解={User.__annotations__}")


if __name__ == "__main__":
    main()
