"""连接池与资源耗尽。"""
import threading
import time


class Pool:
    def __init__(self, size: int) -> None:
        self._sem = threading.Semaphore(size)
        self.in_use = 0
        self.peak = 0
        self._mu = threading.Lock()

    def acquire(self, wait_ms: int | None = None) -> bool:
        if wait_ms is None:
            ok = self._sem.acquire()          # 无超时：一直等
        else:
            ok = self._sem.acquire(timeout=wait_ms / 1000)
        if ok:
            with self._mu:
                self.in_use += 1
                self.peak = max(self.peak, self.in_use)
        return ok

    def release(self) -> None:
        with self._mu:
            self.in_use -= 1
        self._sem.release()


def run(pool: Pool, n: int, hold_ms: int, wait_ms: int | None) -> tuple[int, int, int]:
    ok = 0
    timeout = 0
    lock = threading.Lock()

    def worker() -> None:
        nonlocal ok, timeout
        if pool.acquire(wait_ms):
            try:
                time.sleep(hold_ms / 1000)
                with lock:
                    ok += 1
            finally:
                pool.release()
        else:
            with lock:
                timeout += 1

    t0 = time.monotonic()
    ts = [threading.Thread(target=worker) for _ in range(n)]
    for t in ts:
        t.start()
    for t in ts:
        t.join()
    return ok, timeout, int((time.monotonic() - t0) * 1000)


def main() -> None:
    print("== 1. 池大小 3，10 个请求，无等待超时 ==")
    p = Pool(3)
    ok, to, cost = run(p, 10, 120, None)
    print(f"ok={ok} timeout={to} 耗时={cost}ms peak={p.peak}")

    print("== 2. 同样的池，等待超过 40ms 就快速失败 ==")
    p = Pool(3)
    ok, to, cost = run(p, 10, 120, 40)
    print(f"ok={ok} timeout={to} 耗时={cost}ms peak={p.peak}")

    print("== 3. 结论 ==")
    print("连接池/并发上限必须配等待超时：没有超时，故障时请求无限排队")
    print("Python 里 ThreadPoolExecutor / asyncpg Pool 都要显式配 timeout")


if __name__ == "__main__":
    main()
