import threading
from concurrent.futures import ThreadPoolExecutor


def worker(name: str, counter: dict, lock: threading.Lock) -> None:
    # 每个线程把共享计数加 1000 次；必须加锁，否则 += 会丢更新
    for _ in range(1000):
        with lock:
            counter["value"] += 1
    with lock:
        counter["threads"].append(name)


def square(n: int) -> int:
    return n * n


def main() -> None:
    print("1) threading.Thread：起 4 个线程并 join 等待")
    counter = {"value": 0, "threads": []}
    lock = threading.Lock()
    threads = [threading.Thread(target=worker, args=(f"t{i}", counter, lock)) for i in range(4)]
    for t in threads:
        t.start()
    for t in threads:
        t.join(timeout=1.0)
    print(f"  共享计数（加锁）={counter['value']}，期望 4000，正确={counter['value'] == 4000}")
    print(f"  完成的线程名={sorted(counter['threads'])}")
    print(f"  所有线程已结束={all(not t.is_alive() for t in threads)}")

    print("2) 主线程身份、daemon 属性")
    print(f"  主线程名={threading.current_thread().name}；活跃线程数={threading.active_count()}")
    bg = threading.Thread(target=lambda: None, daemon=True)
    bg.start()
    bg.join(timeout=1.0)
    print(f"  daemon 线程是守护线程={bg.daemon}，已结束={not bg.is_alive()}")

    print("3) ThreadPoolExecutor：map / submit 的确定性结果")
    with ThreadPoolExecutor(max_workers=3) as pool:
        mapped = list(pool.map(square, range(8)))
        futures = [pool.submit(square, n) for n in range(5)]
        collected = [f.result(timeout=1.0) for f in futures]
    print(f"  pool.map 结果={mapped}")
    print(f"  submit+result 结果={collected}")

    print("4) GIL：线程适合 I/O 密集，CPU 密集需 multiprocessing")
    print("  GIL 下字节码串行；本例全是极小任务，秒级内自行结束")


if __name__ == "__main__":
    main()
