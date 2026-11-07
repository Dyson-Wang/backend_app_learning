"""慢查询拖垮连接池。"""


def simulate(durations: list[int], pool: int, wait_timeout: int) -> tuple[int, int, int]:
    free_at = [0] * pool
    ok = timeout = total = 0
    for d in durations:
        best = min(range(pool), key=lambda i: free_at[i])
        wait = free_at[best]
        if wait > wait_timeout:
            timeout += 1
            continue
        free_at[best] = wait + d
        total = max(total, free_at[best])
        ok += 1
    return ok, timeout, total


def main() -> None:
    normal = [30] * 10

    print("== 1. 正常：全部 30ms ==")
    ok, to, total = simulate(normal, 2, 100)
    print(f"成功 = {ok}，超时 = {to}，总耗时 = {total}ms")

    print("== 2. 混入 1 个 300ms 慢查询（等待上限 100ms） ==")
    slow = list(normal)
    slow[0] = 300
    ok, to, total = simulate(slow, 2, 100)
    print(f"成功 = {ok}，超时 = {to}，总耗时 = {total}ms")

    print("== 3. 结论 ==")
    print("池子只有 2 个连接，1 个慢查询就占掉一半；等待超时把雪崩挡在入口")
    print("慢查询必须配 statement_timeout + 索引，光加连接数只是把问题推后")


if __name__ == "__main__":
    main()
