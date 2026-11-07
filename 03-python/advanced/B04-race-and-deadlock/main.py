import threading
import time

# 1. 无锁读改写：Barrier 强制两个线程都读完再写
racy = 0
barrier = threading.Barrier(2)


def racy_inc():
    global racy
    value = racy
    barrier.wait()
    racy = value + 1


def main():
    global racy
    print("== 1. 无锁读改写：两个线程各加一次 ==")
    t1 = threading.Thread(target=racy_inc)
    t2 = threading.Thread(target=racy_inc)
    t1.start()
    t2.start()
    t1.join()
    t2.join()
    print(f"期望 2，实际 {racy}（丢更新）")

    print("== 2. 加锁：4 线程 × 1000 次 ==")
    lock = threading.Lock()
    counter = {"v": 0}

    def worker():
        for _ in range(1000):
            with lock:
                counter["v"] += 1

    threads = [threading.Thread(target=worker) for _ in range(4)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()
    print(f"期望 4000，实际 {counter['v']}")

    print("== 3. 反序加锁导致死锁 ==")
    a = threading.Lock()
    b = threading.Lock()
    done = threading.Event()

    def order_ab():
        with a:
            time.sleep(0.01)
            with b:
                done.set()

    def order_ba():
        with b:
            time.sleep(0.01)
            with a:
                done.set()

    t1 = threading.Thread(target=order_ab, daemon=True)
    t2 = threading.Thread(target=order_ba, daemon=True)
    t1.start()
    t2.start()

    if done.wait(0.1):
        print("未死锁")
    else:
        print("检测到死锁：100ms 内无人完成（各持一把锁等对方）")

    print("== 4. 固定加锁顺序 ==")
    c = threading.Lock()
    d = threading.Lock()
    finished = []

    def lock_both(name):
        with c:
            time.sleep(0.005)
            with d:
                finished.append(name)

    ts = [threading.Thread(target=lock_both, args=(f"thread-{i}",)) for i in (1, 2)]
    for t in ts:
        t.start()
    for t in ts:
        t.join(timeout=1)
    print("  完成:", sorted(finished))

    print("结论：死锁靠统一加锁顺序或带超时的 acquire 避免；daemon 线程不阻止进程退出")


if __name__ == "__main__":
    main()
