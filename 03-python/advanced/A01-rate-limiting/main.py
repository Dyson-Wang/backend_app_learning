class TokenBucket:
    """burst 个突发额度，按 rate/s 匀速补充。"""

    def __init__(self, rate, burst, clock=lambda: 0.0):
        self.rate = rate
        self.burst = float(burst)
        self.tokens = float(burst)
        self.last = clock()
        self.clock = clock

    def allow(self):
        now = self.clock()
        elapsed = now - self.last
        self.tokens = min(self.burst, self.tokens + elapsed * self.rate)
        self.last = now
        if self.tokens >= 1.0:
            self.tokens -= 1.0
            return True
        return False


def main():
    state = {"t": 0.0}
    bucket = TokenBucket(rate=2, burst=5, clock=lambda: state["t"])

    def hit(label):
        ok = bucket.allow()
        print(f"  {label} -> allow={ok}, tokens={bucket.tokens:.2f}")
        return ok

    print("== t=0ms：burst=5，连续 6 次 ==")
    for i in range(1, 7):
        hit(f"#{i}")

    print("== t=500ms：补 1 个令牌 ==")
    state["t"] = 0.5
    hit("#7")
    hit("#8")

    print("== t=1500ms：补 2 个令牌 ==")
    state["t"] = 1.5
    hit("#9")
    hit("#10")
    hit("#11")

    print("结论：令牌桶允许突发（burst），但长期速率被 rate 卡死")


if __name__ == "__main__":
    main()
