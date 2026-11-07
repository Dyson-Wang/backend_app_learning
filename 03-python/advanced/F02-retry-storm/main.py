"""重试风暴。"""

CLIENTS = 100
DOWNSTREAM_CAP = 100


def main() -> None:
    print("== 1. 无退避：失败立即重试 3 次 ==")
    no_backoff = CLIENTS * (1 + 3)
    print(f"下游收到 {no_backoff} QPS（容量 {DOWNSTREAM_CAP}，超载 {no_backoff // DOWNSTREAM_CAP} 倍）")

    print("== 2. 指数退避 + 全抖动 ==")
    peak = CLIENTS + CLIENTS // 3
    print(f"下游峰值 ≈ {peak} QPS（接近容量，不再放大）")

    print("== 3. 再加熔断 ==")
    print("下游请求 = 0 QPS（open 状态本地快速失败，直接返回降级结果）")

    print("== 4. 结论 ==")
    print("重试必须有上限、退避和抖动；三个都没有，就是给故障加杠杆")
    print("熔断才是止损的关键：下游已经挂了，再发请求只会让它更晚恢复")


if __name__ == "__main__":
    main()
