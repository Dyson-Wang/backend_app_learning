import importlib
import importlib.util
import json
import sys
import types


def main() -> None:
    print("1) 解释器与隔离环境：executable / prefix / base_prefix")
    print(f"  sys.executable 非空={bool(sys.executable)}")
    print(f"  在虚拟环境中={sys.prefix != sys.base_prefix}（本仓库 basics 不需要 venv）")
    print(f"  sys.path 是列表={isinstance(sys.path, list)}；条目数>0={len(sys.path) > 0}")

    print("2) 查找模块：find_spec 判断能否导入，不真正执行")
    ok = importlib.util.find_spec("json")
    bad = importlib.util.find_spec("no_such_pkg_xyz_987")
    print(f"  find_spec('json') 找到={ok is not None}；不存在的包={bad is None}")
    print(f"  已加载模块缓存: 'json' in sys.modules = {'json' in sys.modules}")

    print("3) 动态创建模块：ModuleType + sys.modules 注册 + import_module 取回")
    m = types.ModuleType("brm_demo_pkg", "一个内存里造出来的包")
    m.VERSION = "1.0"
    m.__all__ = ["VERSION"]
    sys.modules["brm_demo_pkg"] = m
    got = importlib.import_module("brm_demo_pkg")
    print(f"  import_module 取回同一对象={got is m}；VERSION={got.VERSION}")
    print(f"  import_module('json') 命中缓存={importlib.import_module('json') is sys.modules['json']}")

    print("4) import ... as 别名与 from ... import 的命名空间")
    import math as mth
    from json import dumps as to_json
    print(f"  别名 mth.sqrt(16)={mth.sqrt(16)}；局部别名 to_json={to_json({'k': 1})}")

    print("5) 包 = 带 __init__.py 的目录：__path__ / __package__")
    print(f"  json 是包(有 __path__)={hasattr(json, '__path__')}；json.__name__={json.__name__}")
    print(f"  本文件作为脚本运行时 __package__={__package__!r}")

    print("6) 项目级依赖文件（本课只讲解，不创建）")
    print("  pip + pyproject.toml / requirements.txt；venv 隔离；__pycache__ 字节码缓存")


if __name__ == "__main__":
    main()
