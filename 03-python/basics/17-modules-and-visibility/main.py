import importlib
import json
import sys
import types


class Vault:
    def __init__(self) -> None:
        self.__secret = 42  # 名称改写：实际变成 _Vault__secret

    def reveal(self) -> int:
        return self.__secret


def main() -> None:
    print("1) 私有约定与 __all__：动态模块 + 星号导入")
    pkg = types.ModuleType("brm_vis", "演示可见性的内存包")
    pkg.__all__ = ["PUBLIC", "func"]
    pkg.PUBLIC = "可见"
    pkg._hidden = "私有约定"
    pkg.func = lambda: "ok"
    sys.modules["brm_vis"] = pkg
    ns: dict = {}
    exec("from brm_vis import *", ns)
    print(f"  星号导入拿到 func={ns['func']()}；PUBLIC 存在={'PUBLIC' in ns}")
    print(f"  _hidden 被 __all__ 挡住={'_hidden' not in ns}；但显式取仍={getattr(pkg, '_hidden')!r}")

    print("2) 名称改写 __x -> _Class__x（不是真私有）")
    v = Vault()
    print(f"  reveal()={v.reveal()}；hasattr(_Vault__secret)={hasattr(v, '_Vault__secret')}")
    print(f"  hasattr(__secret)={hasattr(v, '__secret')}；实际字段名={[k for k in vars(v) if 'secret' in k]}")

    print("3) import ... as 别名解决命名冲突、包目录即命名空间")
    import math as mth
    import json as js
    print(f"  别名 mth.pi>3={mth.pi > 3}；js.dumps={js.dumps({'a': 1})}")
    print(f"  json 是包(有 __path__)={hasattr(json, '__path__')}；包名={json.__package__!r}")

    print("4) 循环依赖：把 import 放进函数（惰性）打破导入期环")
    a = types.ModuleType("brm_a")
    b = types.ModuleType("brm_b")
    sys.modules["brm_a"] = a
    sys.modules["brm_b"] = b

    def a_value():
        # 调用时才 import，避免两个模块在顶层互相引用
        mod = importlib.import_module("brm_b")
        return mod.VALUE  # type: ignore[attr-defined]

    a.a_value = a_value  # type: ignore[attr-defined]
    b.VALUE = 100  # type: ignore[attr-defined]
    print(f"  brm_a 惰性取 brm_b.VALUE={a.a_value()}（无导入期死锁）")
    print(f"  brm_a 已登记在 sys.modules={'brm_a' in sys.modules}")

    print("5) __name__ 守卫：被 import 时不执行 main")
    print(f"  本文件 __name__={__name__}；作为脚本运行 __spec__ 为 None={__spec__ is None}")


if __name__ == "__main__":
    main()
