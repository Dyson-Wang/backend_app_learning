import queue
import threading

capacity = 5
total = 20

q = queue.Queue(maxsize=capacity)
gate = threading.Event()
result = {}


def consumer():
    gate.wait()
    processed = 0
    while True:
        try:
            q.get_nowait()
        except queue.Empty:
            break
        processed += 1
        q.task_done()
    result["processed"] = processed


worker = threading.Thread(target=consumer)
worker.start()

accepted = 0
rejected = 0
for i in range(total):
    try:
        q.put_nowait(i)
        accepted += 1
    except queue.Full:
        rejected += 1  # 满则快速拒绝

gate.set()
worker.join()

print(f"容量={capacity}，提交={total}")
print(f"accepted={accepted} rejected={rejected} processed={result['processed']}")
print("结论：拒绝发生在入口，内存和延迟都是常数；无界 Queue 只是把 OOM 推迟")

unbounded = list(range(total))
print(f"无界队列长度={len(unbounded)}（全部接收，没有上限）")
