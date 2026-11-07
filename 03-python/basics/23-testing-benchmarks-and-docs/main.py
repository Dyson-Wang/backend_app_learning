import doctest
import inspect
import io
import pydoc
import timeit
import unittest


def add(a: int, b: int) -> int:
    """返回两个整数的和，供 doctest / pydoc 演示。

    >>> add(2, 3)
    5
    >>> add(-1, 1)
    0
    >>> add(0, 0)
    0
    """
    return a + b


class TestAdd(unittest.TestCase):
    # 表驱动测试：一组 (输入, 期望) 用 subTest 逐个跑，失败信息带参数
    def test_table(self) -> None:
        cases = [(1, 2, 3), (0, 0, 0), (-1, 1, 0), (10, -3, 7)]
        for a, b, expected in cases:
            with self.subTest(a=a, b=b):
                self.assertEqual(add(a, b), expected, f"{a}+{b} 应等于 {expected}")

    def test_raises(self) -> None:
        with self.assertRaises(ValueError):
            int("not-a-number")


def run_tests() -> None:
    print("1) unittest：表驱动 subTest + 失败信息 + assertRaises")
    suite = unittest.defaultTestLoader.loadTestsFromTestCase(TestAdd)
    result = unittest.TextTestRunner(stream=io.StringIO(), verbosity=0).run(suite)
    print(f"  通过={result.wasSuccessful()}；用例数={result.testsRun}；失败数={len(result.failures)}")


def run_doctest() -> None:
    print("2) doctest：文档字符串里的 >>> 就是可执行断言")
    finder = doctest.DocTestFinder()
    runner = doctest.DocTestRunner()
    for test in finder.find(add, "add", globs=globals()):
        runner.run(test)
    res = runner.summarize(verbose=False)
    print(f"  示例数={res.attempted}；失败数={res.failed}；全部通过={res.failed == 0}")


def run_benchmark() -> None:
    print("3) timeit 基准：固定 number，只演示 API 不依赖具体数值")
    elapsed = timeit.timeit("sum(range(100))", number=2000)
    print(f"  重复次数=2000；耗时为正={elapsed > 0}")
    print(f"  timeit.repeat 返回组数={len(timeit.repeat('1+1', number=10, repeat=2))}")


def run_docs() -> None:
    print("4) 文档字符串 / inspect.getdoc / pydoc 生成")
    print(f"  add.__doc__ 存在={add.__doc__ is not None}；含示例={'>>>' in (add.__doc__ or '')}")
    print(f"  inspect.getdoc 首句={inspect.getdoc(add).splitlines()[0]}")
    rendered = pydoc.render_doc(add)
    print(f"  pydoc.render_doc 含函数名={'add' in rendered}")


def main() -> None:
    run_tests()
    run_doctest()
    run_benchmark()
    run_docs()
    print("5) 入口脚本可直接 python3 main.py 运行（无需 pytest，覆盖率可用 trace 模块）")


if __name__ == "__main__":
    main()
