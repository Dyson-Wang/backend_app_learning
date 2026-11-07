"""顺序性：状态机 + 分区。"""
from dataclasses import dataclass


@dataclass
class Event:
    id: str
    kind: str


def apply(state: str, e: Event) -> tuple[str, str]:
    if e.kind == "create" and state == "":
        return "CREATED", ""
    if e.kind == "pay" and state == "CREATED":
        return "PAID", ""
    if e.kind == "ship" and state == "PAID":
        return "SHIPPED", ""
    return state, f"{e.kind} 在 {state or 'INIT'} 状态下不合法"


def consume(events: list[Event]) -> tuple[str, str]:
    state = ""
    for e in events:
        state, err = apply(state, e)
        if err:
            return state, err
    return state, ""


def partition(key: str, n: int) -> int:
    h = 0
    for ch in key:
        h = (h * 31 + ord(ch)) % n
    return h


def main() -> None:
    create, pay, ship = Event("c1", "create"), Event("c2", "pay"), Event("c3", "ship")

    print("== 1. 同一 key 被并发消费（顺序打乱） ==")
    state, err = consume([ship, create, pay])
    print(f"最终状态 = {state or 'INIT'}，错误 = {err}")

    print("== 2. 按 key 路由到同一分区、单消费者串行处理 ==")
    state, err = consume([create, pay, ship])
    print(f"最终状态 = {state}，错误 = {err or '无'}")

    print("== 3. 分区函数：同一 order_id 必然落到同一分区 ==")
    for k in ["order-42", "order-43", "order-42"]:
        print(f"{k} -> 分区 {partition(k, 4)}")

    print("== 4. 结论 ==")
    print("顺序性不是全局的，是「同一分区内」的；跨 key 的顺序本来就没有意义")
    print("要顺序就牺牲并发：单分区单消费者，或者消费端按 key 做串行队列")


if __name__ == "__main__":
    main()
