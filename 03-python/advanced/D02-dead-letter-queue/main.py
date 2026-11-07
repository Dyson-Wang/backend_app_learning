"""死信队列与重试策略。"""
from dataclasses import dataclass


@dataclass
class Msg:
    id: str
    failures: int          # 前 N 次投递失败（可重试）
    permanent: bool = False  # 永久错误，重试也没用


def deliver(msg: Msg, max_attempts: int) -> tuple[bool, int]:
    for attempt in range(1, max_attempts + 1):
        if msg.permanent:
            print(f"{msg.id} -> 永久错误，直接进死信（不重试）")
            return False, attempt
        if attempt > msg.failures:
            print(f"{msg.id} -> 成功（第 {attempt} 次）")
            return True, attempt
        if attempt < max_attempts:
            print(f"{msg.id} -> 可重试错误，退避 {10 * 2 ** (attempt - 1)}ms 后重试")
            continue
        print(f"{msg.id} -> 重试耗尽，进死信")
    return False, max_attempts


def main() -> None:
    msgs = [
        Msg("m1", 0),
        Msg("m2", 2),
        Msg("m3", 0, permanent=True),
        Msg("m4", 99),
        Msg("m5", 0),
    ]

    print("== 1. 逐条投递（最多 3 次） ==")
    ok = dead = total = 0
    for m in msgs:
        success, attempts = deliver(m, 3)
        total += attempts
        if success:
            ok += 1
        else:
            dead += 1

    print("== 2. 统计 ==")
    print(f"成功 = {ok}，死信 = {dead}，总投递 = {total}")

    print("== 3. 结论 ==")
    print("永久错误不要重试，重试只会拖慢整个队列")
    print("死信队列必须有人看：它进 DLQ 说明要么代码有 bug，要么数据有问题")


if __name__ == "__main__":
    main()
