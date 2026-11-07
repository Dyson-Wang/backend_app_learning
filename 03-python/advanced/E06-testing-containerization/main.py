"""极简测试框架：只看失败输出长什么样。"""

failures = 0


def check(name: str, ok: bool, detail: str) -> None:
    global failures
    if ok:
        print(f"PASS {name}")
        return
    failures += 1
    print(f"FAIL {name}: {detail}")


def add(a: int, b: int) -> int:
    return a + b


def divide(a: int, b: int) -> float:
    if b == 0:
        raise ValueError("division by zero")
    return a / b


def parse_age(s: str) -> int:
    if s == "":
        raise ValueError("empty")
    return 18            # 故意不校验非数字，用来触发失败用例


def main() -> None:
    print("== 1. 单元测试 ==")
    check("add(1,2)=3", add(1, 2) == 3, "期望 3")

    threw = False
    try:
        divide(1, 0)
    except ValueError:
        threw = True
    check("divide(1,0) 抛错", threw, "期望抛错")

    age_err = None
    try:
        parse_age("abc")
    except ValueError as e:
        age_err = e
    check('parse_age("abc") 抛错', age_err is not None, "期望抛错，实际正常返回 18")

    print("== 2. 汇总 ==")
    print(f"3 个用例：{3 - failures} 通过，{failures} 失败")

    print("== 3. 结论 ==")
    print("最后一条是故意失败的用例，用来看失败输出长什么样")
    print("单测不碰网络/DB；集成测试才起真实依赖，两者分开跑")
    print("容器化：多阶段构建、非 root、HEALTHCHECK 打 /healthz")


if __name__ == "__main__":
    main()
