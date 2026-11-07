"""事务与隔离级别、乐观锁。"""
import threading


class Account:
    def __init__(self, balance: int) -> None:
        self.balance = balance
        self.version = 0
        self._lock = threading.Lock()

    def read(self) -> tuple[int, int]:
        with self._lock:
            return self.balance, self.version

    def cas(self, expected_version: int, new_balance: int) -> bool:
        # 只有 version 没被别人改过才写入，等价于 UPDATE ... WHERE version = ?
        with self._lock:
            if self.version != expected_version:
                return False
            self.balance = new_balance
            self.version += 1
            return True


def naive_transfer(acc: Account, read_balance: int, amount: int) -> None:
    # 错误写法：拿着事务开始时读到的旧值直接写回
    acc.balance = read_balance - amount


def optimistic_debit(acc: Account, amount: int, max_retry: int = 10) -> tuple[bool, int]:
    for attempt in range(max_retry + 1):
        balance, version = acc.read()
        if acc.cas(version, balance - amount):
            return True, attempt
    return False, max_retry + 1


def main() -> None:
    print("== 1. 无锁：两个事务都基于旧值写回 ==")
    acc = Account(100)
    b1, _ = acc.read()
    b2, _ = acc.read()          # 两个事务都读到 100
    naive_transfer(acc, b1, 10)
    naive_transfer(acc, b2, 10)  # 用同样过期的 100 覆盖
    print(f"期望 80，实际 {acc.balance}（丢失更新）")

    print("== 2. 乐观锁：version 校验 + 重试 ==")
    acc = Account(100)
    ok1, a1 = optimistic_debit(acc, 10)
    print(f"T1 ok={ok1} attempts={a1}，balance={acc.balance} version={acc.version}")
    stale = acc.version - 1
    print(f"T2 用旧 version={stale} 写入：{acc.cas(stale, 70)}")
    ok2, a2 = optimistic_debit(acc, 10)
    print(f"T2 重试 ok={ok2} attempts={a2}，balance={acc.balance} version={acc.version}")
    print(f"期望 80，实际 {acc.balance}")

    print("== 3. 结论 ==")
    print("悲观锁锁整行，吞吐低；乐观锁读不加锁、写校验 version，冲突才重试")
    print("读已提交 + 乐观锁是多数业务的折中，范围查询要小心幻读")


if __name__ == "__main__":
    main()
