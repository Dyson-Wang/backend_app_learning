import time


class DownstreamError(Exception):
    def __init__(self, message, retryable):
        super().__init__(message)
        self.retryable = retryable


def do_with_retry(max_attempts, base_ms, fn):
    last = None
    for attempt in range(1, max_attempts + 1):
        try:
            return {"ok": True, "attempt": attempt, "value": fn(attempt)}
        except DownstreamError as e:
            last = e
            if not e.retryable:
                return {"ok": False, "attempt": attempt, "error": str(e), "permanent": True}
            if attempt == max_attempts:
                break
            delay = base_ms * (2 ** (attempt - 1))
            print(f"  attempt {attempt} 失败（{e}），{delay}ms 后重试")
            time.sleep(delay / 1000)
    return {"ok": False, "attempt": max_attempts, "error": str(last), "permanent": False}


def main():
    calls = {"n": 0}

    def flaky(_attempt):
        calls["n"] += 1
        if calls["n"] < 4:
            raise DownstreamError("503 upstream", retryable=True)
        return f"第 {calls['n']} 次成功"

    def permanent(_attempt):
        raise DownstreamError("400 bad request", retryable=False)

    def always_fail(_attempt):
        raise DownstreamError("timeout", retryable=True)

    print("== 1. 可重试错误：指数退避后成功 ==")
    print("结果:", do_with_retry(4, 10, flaky))

    print("== 2. 永久错误：立即停止 ==")
    print("结果:", do_with_retry(4, 10, permanent))

    print("== 3. 一直失败：耗尽 3 次 ==")
    print("结果:", do_with_retry(3, 5, always_fail))

    print("结论：只重试可重试错误，退避必须指数增长；生产还要加抖动，避免重试风暴")


if __name__ == "__main__":
    main()
