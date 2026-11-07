"""上下文取消与跨服务传播。"""
import asyncio
import time


class Deadline:
    """模拟 Go 的 context.WithTimeout：一个可传递的截止时间。"""

    def __init__(self, budget_ms: int) -> None:
        self.deadline = time.monotonic() + budget_ms / 1000

    def remaining(self) -> float:
        return self.deadline - time.monotonic()

    def expired(self) -> bool:
        return self.remaining() <= 0


async def downstream(deadline: Deadline | None, work_ms: int) -> str:
    # 错误写法：deadline 为 None 就完全不看上游预算
    if deadline is None:
        await asyncio.sleep(work_ms / 1000)
        return "ok"
    if deadline.remaining() < work_ms / 1000:
        # 有预算意识：跑到截止时间就放弃，而不是硬跑完
        await asyncio.sleep(max(0.0, deadline.remaining()))
        return "deadline exceeded"
    await asyncio.sleep(work_ms / 1000)
    return "ok"


async def handler_with_deadline(budget_ms: int, work_ms: int) -> tuple[str, int]:
    d = Deadline(budget_ms)
    t0 = time.monotonic()
    res = await downstream(d, work_ms)
    return res, int((time.monotonic() - t0) * 1000)


async def handler_dropping_deadline(work_ms: int) -> tuple[str, int]:
    t0 = time.monotonic()
    res = await downstream(None, work_ms)   # 中间层忘记传
    return res, int((time.monotonic() - t0) * 1000)


async def main() -> None:
    print("== 1. 上游预算 50ms，下游要 200ms ==")
    res, cost = await handler_with_deadline(50, 200)
    print(f"res = {res}，耗时 = {cost}ms（远小于 200ms）")

    print("== 2. 中间层丢掉 deadline ==")
    res, cost = await handler_dropping_deadline(200)
    print(f"res = {res}，耗时 = {cost}ms（上游再急也没用）")

    print("== 3. 正确写法：asyncio.timeout 传下去 ==")
    t0 = time.monotonic()
    try:
        async with asyncio.timeout(0.05):
            await asyncio.sleep(0.2)
    except TimeoutError:
        print(f"TimeoutError，耗时 = {int((time.monotonic() - t0) * 1000)}ms")

    print("== 4. 结论 ==")
    print("deadline 必须一路显式传；任何一层忘记传，取消就断在这里")
    print("Python 用 asyncio.timeout / TaskGroup，跨服务用 HTTP 头传剩余预算")


if __name__ == "__main__":
    asyncio.run(main())
