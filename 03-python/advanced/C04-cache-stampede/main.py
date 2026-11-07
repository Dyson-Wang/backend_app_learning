import threading
import time
from concurrent.futures import ThreadPoolExecutor


class FakeDB:
    """一次查询 20ms，并统计回源次数。"""

    def __init__(self):
        self.hits = 0
        self._lock = threading.Lock()

    def query(self, key):
        with self._lock:
            self.hits += 1
        time.sleep(0.02)
        return f"value-of-{key}"


def unprotected_get(db, cache, lock, key):
    with lock:
        value = cache.get(key)
    if value is not None:
        return value
    value = db.query(key)  # 并发请求全部走到这里
    with lock:
        cache[key] = value
    return value


class GuardedCache:
    """按 key 加锁 + 双重检查，只放一个请求回源。"""

    def __init__(self):
        self.data = {}
        self.key_locks = {}
        self.lock = threading.Lock()

    def _key_lock(self, key):
        with self.lock:
            return self.key_locks.setdefault(key, threading.Lock())

    def get(self, db, key):
        with self.lock:
            if key in self.data:
                return self.data[key]

        key_lock = self._key_lock(key)
        with key_lock:
            with self.lock:
                if key in self.data:
                    return self.data[key]
            value = db.query(key)
            with self.lock:
                self.data[key] = value
            return value


def run(n, guarded):
    db = FakeDB()
    start = time.perf_counter()
    if guarded:
        cache = GuardedCache()
        with ThreadPoolExecutor(max_workers=n) as pool:
            list(pool.map(lambda _: cache.get(db, "user:1"), range(n)))
    else:
        cache = {}
        lock = threading.Lock()
        with ThreadPoolExecutor(max_workers=n) as pool:
            list(pool.map(lambda _: unprotected_get(db, cache, lock, "user:1"), range(n)))
    return db.hits, (time.perf_counter() - start) * 1000


def main():
    n = 50
    hits, ms = run(n, False)
    print(f"无保护：{n} 并发 -> 回源 {hits} 次，耗时 {ms:.0f}ms")

    hits, ms = run(n, True)
    print(f"按 key 加锁 + 双重检查：{n} 并发 -> 回源 {hits} 次，耗时 {ms:.0f}ms")

    print("结论：缓存击穿是「同一热 key 失效瞬间的并发回源」，GIL 并不能阻止它")


if __name__ == "__main__":
    main()
