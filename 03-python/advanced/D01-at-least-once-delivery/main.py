"""至少一次投递：重复消息必须由消费端去重。"""

DELIVERIES = [
    ("e1", 10),
    ("e2", 10),
    ("e2", 10),   # ACK 丢了，Broker 重投
    ("e3", 10),
]


def run(seen: set[str] | None) -> tuple[int, int]:
    balance = 100
    applied = 0
    for msg_id, amount in DELIVERIES:
        if seen is not None:
            if msg_id in seen:
                continue            # 已处理，直接 ACK
            seen.add(msg_id)
        balance += amount
        applied += 1
    return balance, applied


def main() -> None:
    print("== 1. 消费端没有幂等表 ==")
    balance, _ = run(None)
    print(f"投递 {len(DELIVERIES)} 次，余额 100 -> {balance}（多加了 10）")

    print("== 2. 消费端按消息 ID 去重 ==")
    balance, applied = run(set())
    print(f"投递 {len(DELIVERIES)} 次，实际处理 {applied} 次，余额 100 -> {balance}")

    print("== 3. 结论 ==")
    print("至少一次投递下，重复是常态；正确性必须由消费端幂等保证")
    print("去重表要落库并和业务写在同一事务里，否则还是会重复")


if __name__ == "__main__":
    main()
