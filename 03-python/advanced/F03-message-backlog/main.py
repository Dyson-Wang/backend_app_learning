"""消息积压。"""

BACKLOG = 1000
PER_SEC = 10      # 单消费者每秒 10 条
PRODUCERS = 20    # 生产端每秒 20 条


def main() -> None:
    print("== 1. 单消费者，逐条处理 ==")
    print(f"{BACKLOG} 条 / {PER_SEC} 条每秒 = {BACKLOG // PER_SEC} 秒")

    print("== 2. 扩到 3 个消费者 ==")
    three = -(-BACKLOG // (PER_SEC * 3))
    print(f"{BACKLOG} 条 / {PER_SEC * 3} 条每秒 = {three} 秒")

    print("== 3. 批量 10 条一次（单消费者） ==")
    batch = -(-BACKLOG // (PER_SEC * 10))
    print(f"{BACKLOG} 条 / {PER_SEC * 10} 条每秒 = {batch} 秒")

    print("== 4. 生产仍在继续（每秒 20 条） ==")
    net = PER_SEC * 10 - PRODUCERS
    clear = -(-BACKLOG // net)
    print(f"净消化 {net} 条每秒，清空 {BACKLOG} 条需要 {clear} 秒")

    print("== 5. 结论 ==")
    print("积压先算「净消化速率」，不是「消费速率」——生产端还在往里灌")
    print("扩消费者最快但受分区数限制；批量提升单机吞吐，代价是延迟和失败重放")


if __name__ == "__main__":
    main()
