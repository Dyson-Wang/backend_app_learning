import gc
import sys
import tracemalloc
import weakref


class Node:
    # 普通类默认带 __weakref__，可被弱引用
    def __init__(self, name: str) -> None:
        self.name = name

    def __repr__(self) -> str:
        return f"Node({self.name})"


events: list[str] = []


class Finalized:
    def __init__(self, name: str) -> None:
        self.name = name

    def __del__(self) -> None:
        # 析构器：引用计数归零即触发（CPython）
        events.append(self.name)


def main() -> None:
    print("1) 值与引用：名字是贴在对象上的标签")
    a = [1, 2]
    b = a  # b 与 a 指向同一堆对象
    b.append(3)
    print(f"  b.append 后 a={a}；a is b={a is b}（别名共享）")

    print("2) sys.getrefcount：引用计数增减可观察（注意它会临时 +1）")
    o = object()
    base = sys.getrefcount(o)
    alias = o
    print(f"  新增一个引用后 +1={sys.getrefcount(o) == base + 1}")
    del alias
    print(f"  删除别名后复原={sys.getrefcount(o) == base}")

    print("3) 引用计数即时释放 + __del__")
    obj = Finalized("tmp")
    del obj
    print(f"  del 之后 __del__ 已跑: events={events}")

    print("4) 引用循环 + gc.collect() + weakref 观察回收")
    x = Node("x")
    y = Node("y")
    x.peer = y
    y.peer = x  # 互相引用，引用计数永不归零
    ref = weakref.ref(x, lambda _: events.append("weakref 回调"))
    del x, y
    found = gc.collect()  # 分代 GC 专门收拾循环
    print(f"  gc.collect 找到不可达对象>0={found > 0}；弱引用已失效={ref() is None}")
    print(f"  弱引用回调被触发={'weakref 回调' in events}")

    print("5) gc 开关与阈值、对象是否被跟踪")
    print(f"  gc.isenabled()={gc.isenabled()}；gc.is_tracked([])={gc.is_tracked([])}；gc.is_tracked(1)={gc.is_tracked(1)}")
    print(f"  分代阈值是三元组={len(gc.get_threshold()) == 3}")

    print("6) tracemalloc 追踪堆分配（只比较量级，不打印绝对值）")
    tracemalloc.start()
    before = tracemalloc.get_traced_memory()[0]
    chunk = [bytearray(4096) for _ in range(64)]
    cur, peak = tracemalloc.get_traced_memory()
    tracemalloc.stop()
    print(f"  分配后已追踪字节增长={cur > before}；峰值>0={peak > 0}；chunk 长度={len(chunk)}")

    print("7) GIL 对内存模型的影响")
    print("  引用计数的增减在 GIL 下是原子字节码，线程安全；但复合操作仍需 Lock（见 09/19 课）")
    print("  Python 无所有权/移动语义：对象全在堆上，栈帧只存引用")


if __name__ == "__main__":
    main()
