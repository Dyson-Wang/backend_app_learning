import asyncio
import contextvars

request_id: contextvars.ContextVar[str] = contextvars.ContextVar("request_id", default="root")


async def fetch(name: str, delay: float) -> str:
    # 上下文传值：每个 Task 复制一份 Context，互不干扰
    request_id.set(name)
    await asyncio.sleep(delay)
    return f"{name} 完成，ctx={request_id.get()}"


async def slow_cleanup() -> None:
    # 被取消时在 finally 里清理，让取消传播可观察
    try:
        await asyncio.sleep(1)
    finally:
        print("  [cleanup] slow_cleanup 被取消，finally 已执行")


async def demo_timeout() -> None:
    print("1) asyncio.timeout：超时取消内部协程并转成 TimeoutError")
    try:
        async with asyncio.timeout(0.02):
            await asyncio.sleep(1)
    except TimeoutError:
        print("  捕获 TimeoutError（而非 CancelledError）")


async def demo_cancel() -> None:
    print("2) 显式取消 + finally 清理")
    task = asyncio.create_task(slow_cleanup())
    await asyncio.sleep(0.01)  # 确保任务已进入 sleep(1)
    task.cancel()
    try:
        await task
    except asyncio.CancelledError:
        print(f"  task.cancel() 后 await 抛 CancelledError；task.cancelled()={task.cancelled()}")


async def demo_taskgroup() -> None:
    print("3) TaskGroup：并发等待多个任务，任一失败则取消其余")
    async with asyncio.TaskGroup() as tg:
        t1 = tg.create_task(fetch("A", 0.01))
        t2 = tg.create_task(fetch("B", 0.01))
        t3 = tg.create_task(fetch("C", 0.01))
    print(f"  全部完成: {t1.result()!r} | {t2.result()!r} | {t3.result()!r}")
    print(f"  外层 request_id 仍是 root: {request_id.get() == 'root'}")


async def demo_exception_group() -> None:
    print("4) TaskGroup 把异常聚合成 ExceptionGroup，用 except* 捕获")
    async def boom() -> None:
        raise ValueError("故意失败")
    try:
        async with asyncio.TaskGroup() as tg:
            tg.create_task(boom())
            tg.create_task(asyncio.sleep(0.01))
    except* ValueError as eg:
        print(f"  except* 抓到 ExceptionGroup，子异常数={len(eg.exceptions)}")


async def demo_contextvars() -> None:
    print("5) contextvars.ContextVar：任务间上下文隔离 + copy_context")
    request_id.set("main")

    async def child(val: str) -> None:
        request_id.set(val)
        await asyncio.sleep(0.001)
        print(f"  子任务 {val} 看到 request_id={request_id.get()}")

    async with asyncio.TaskGroup() as tg:
        tg.create_task(child("X"))
        tg.create_task(child("Y"))
    print(f"  父任务 request_id 未被污染={request_id.get() == 'main'}")
    ctx = contextvars.copy_context()
    print(f"  copy_context 快照含变量数={len(ctx)}")


async def main() -> None:
    await demo_timeout()
    await demo_cancel()
    await demo_taskgroup()
    await demo_exception_group()
    await demo_contextvars()
    print("6) 全部异步演示结束（总耗时 < 0.2s）")


if __name__ == "__main__":
    asyncio.run(main())
