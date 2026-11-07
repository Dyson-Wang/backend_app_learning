import os
import platform
import sys
from typing import Iterable


def greet(name: str) -> str:
    # 类型注解只是提示，运行时不强制；需要强校验请用 pydantic / mypy
    return f"hello, {name}"


def total(values: Iterable[int]) -> int:
    return sum(values)


def main() -> None:
    print(greet("backend roadmap"))
    print(f"python:   {platform.python_version()} ({platform.python_implementation()})")
    print(f"platform: {sys.platform} {platform.machine()}")
    print(f"pid:      {os.getpid()}")
    print(f"cpu:      {os.cpu_count()}")
    print(f"sum:      {total([1, 2, 3, 4])}")

    # 动态类型：同一个变量可以换成完全不同的类型
    value: object = 42
    print(f"value:    {value!r} ({type(value).__name__})")
    value = "now a string"
    print(f"value:    {value!r} ({type(value).__name__})")

    # GIL：CPython 同一时刻只有一个线程执行字节码
    # CPU 密集用 multiprocessing，I/O 密集用 asyncio / 线程池
    print("note:     GIL 决定 CPU 密集用多进程、I/O 密集用 asyncio")


if __name__ == "__main__":
    main()
