"""幂等键：同一个 Idempotency-Key 只创建一次订单。"""
import threading
import time


class IdemStore:
    def __init__(self) -> None:
        self._mu = threading.Lock()
        self._done: dict[str, dict[str, object]] = {}
        self._inflight: dict[str, threading.Event] = {}
        self.created = 0
        self.replayed = 0

    def create(self, key: str, amount: int) -> tuple[dict[str, object], bool]:
        with self._mu:
            if key in self._done:
                self.replayed += 1
                return self._done[key], True
            ev = self._inflight.get(key)
            owner = ev is None
            if owner:
                ev = threading.Event()
                self._inflight[key] = ev

        if not owner:
            ev.wait(timeout=2)          # 等第一个请求跑完
            with self._mu:
                self.replayed += 1
                return self._done[key], True

        time.sleep(0.02)                # 模拟写库
        with self._mu:
            self.created += 1
            order = {"id": self.created, "key": key, "amount": amount}
            self._done[key] = order
            self._inflight.pop(key, None)
        ev.set()
        return order, False


def main() -> None:
    store = IdemStore()
    key = "idem-7f3a"
    results: list[tuple[dict[str, object], bool]] = []
    lock = threading.Lock()

    def worker() -> None:
        r = store.create(key, 9900)
        with lock:
            results.append(r)

    print("== 1. 20 个并发请求，同一个 Idempotency-Key ==")
    ts = [threading.Thread(target=worker) for _ in range(20)]
    for t in ts:
        t.start()
    for t in ts:
        t.join()

    created = sum(1 for _, replayed in results if not replayed)
    replayed = sum(1 for _, replayed in results if replayed)
    ids = {r[0]["id"] for r in results}
    print(f"真实建单 = {created}，回放 = {replayed}，订单 id 去重后 = {len(ids)}")

    print("== 2. 换一个 key ==")
    order, replayed = store.create("idem-other", 100)
    print(f"新 key 建单 id = {order['id']}，replayed = {replayed}")

    print("== 3. 结论 ==")
    print("幂等键必须落库并加唯一索引；内存字典只是演示，进程重启就没了")
    print("并发同 key 要用「占位 + 等待」而不是「先查再写」，否则照样重复")


if __name__ == "__main__":
    main()
