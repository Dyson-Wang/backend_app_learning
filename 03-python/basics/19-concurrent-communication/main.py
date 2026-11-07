import queue
import threading
from concurrent.futures import ThreadPoolExecutor, as_completed


def consumer(q: queue.Queue, out: list[int], lock: threading.Lock) -> None:
    # 关闭信号：拿到 None（毒丸）就退出；每消费者一个哨兵，避免 task_done 失衡
    while True:
        try:
            item = q.get(timeout=1.0)
        except queue.Empty:
            break
        try:
            if item is None:
                break
            with lock:
                out.append(item * item)
        finally:
            q.task_done()


def double(n: int) -> int:
    return n * 2


def main() -> None:
    print("1) queue.Queue：有界队列 + 生产消费 + task_done/join")
    q: queue.Queue = queue.Queue(maxsize=4)  # 满时 put 阻塞 = 天然背压
    out: list[int] = []
    lock = threading.Lock()
    t = threading.Thread(target=consumer, args=(q, out, lock))
    t.start()
    for n in range(1, 7):
        q.put(n)
    q.put(None)
    q.join()
    t.join(timeout=1.0)
    print(f"  消费结果(平方)={sorted(out)}；队列已空={q.empty()}")
    print(f"  消费者线程结束={not t.is_alive()}")

    print("2) 关闭信号传播：N 个哨兵广播给 N 个消费者（扇出扇入）")
    fan_q: queue.Queue = queue.Queue()
    results: list[int] = []
    rlock = threading.Lock()
    workers = [threading.Thread(target=consumer, args=(fan_q, results, rlock)) for _ in range(3)]
    for w in workers:
        w.start()
    for n in range(1, 8):
        fan_q.put(n)
    for _ in workers:
        fan_q.put(None)  # 一个消费者一个毒丸
    fan_q.join()
    for w in workers:
        w.join(timeout=1.0)
    print(f"  3 个消费者扇入结果={sorted(results)}；全部结束={all(not w.is_alive() for w in workers)}")

    print("3) 同步原语：Event / Semaphore / Barrier 的可见行为")
    ev = threading.Event()
    ev.set()
    sem = threading.Semaphore(2)
    sem.acquire(timeout=1.0)
    sem.acquire(timeout=1.0)
    got = sem.acquire(timeout=0.01)  # 已耗尽，超时返回 False
    print(f"  Event.is_set()={ev.is_set()}；Semaphore(2) 第三次获取={got}")
    sem.release()
    sem.release()
    bar = threading.Barrier(2, timeout=1.0)
    notes: list[str] = []

    def wait_bar() -> None:
        try:
            bar.wait()
            notes.append("过了栅栏")
        except threading.BrokenBarrierError:
            pass

    bthreads = [threading.Thread(target=wait_bar) for _ in range(2)]
    for b in bthreads:
        b.start()
    for b in bthreads:
        b.join(timeout=1.0)
    print(f"  Barrier(2) 两侧都通过={len(notes) == 2}")

    print("4) concurrent.futures：扇出 + as_completed 等待多个任务")
    with ThreadPoolExecutor(max_workers=3) as pool:
        futs = [pool.submit(double, n) for n in range(6)]
        gathered = sorted(f.result(timeout=1.0) for f in as_completed(futs))
    print(f"  as_completed 汇总={gathered}")

    print("5) 共享状态保护：out/results 都靠 Lock，Queue 自身线程安全")


if __name__ == "__main__":
    main()
