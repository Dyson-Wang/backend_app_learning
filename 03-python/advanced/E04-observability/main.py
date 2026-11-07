import contextvars
import json
import logging
import time

trace_id_var = contextvars.ContextVar("trace_id", default="-")
_seq = 0


def next_trace_id():
    global _seq
    _seq += 1
    return f"trace-{_seq:04d}"


class Metrics:
    def __init__(self):
        self.requests = 0
        self.errors = 0
        self.total_ms = 0
        self.max_ms = 0

    def observe(self, status, ms):
        self.requests += 1
        if status >= 500:
            self.errors += 1
        self.total_ms += ms
        self.max_ms = max(self.max_ms, ms)

    def snapshot(self):
        avg = self.total_ms / self.requests if self.requests else 0.0
        return self.requests, self.errors, avg, self.max_ms


class JsonFormatter(logging.Formatter):
    def format(self, record):
        payload = {
            "level": record.levelname.lower(),
            "trace_id": getattr(record, "trace_id", "-"),
            "method": getattr(record, "method", "-"),
            "path": getattr(record, "path", "-"),
            "status": getattr(record, "status", 0),
            "latency_ms": getattr(record, "latency_ms", 0),
        }
        msg = getattr(record, "msg_extra", None)
        if msg:
            payload["msg"] = msg
        return json.dumps(payload, ensure_ascii=False)


handler = logging.StreamHandler()
handler.setFormatter(JsonFormatter())
logger = logging.getLogger("brm")
logger.setLevel(logging.INFO)
logger.addHandler(handler)
logger.propagate = False


def handle(metrics, method, path, status, ms):
    trace_id = next_trace_id()
    token = trace_id_var.set(trace_id)
    try:
        metrics.observe(status, ms)
        extra = {
            "trace_id": trace_id_var.get(),
            "method": method,
            "path": path,
            "status": status,
            "latency_ms": ms,
        }
        if status >= 500:
            extra["msg_extra"] = "upstream failed"
            logger.error("", extra=extra)
        else:
            logger.info("", extra=extra)
    finally:
        trace_id_var.reset(token)


def main():
    metrics = Metrics()
    reqs = [
        ("GET", "/healthz", 200, 1),
        ("GET", "/users", 200, 12),
        ("POST", "/users", 201, 18),
        ("GET", "/users/999", 404, 3),
        ("POST", "/users", 400, 2),
        ("GET", "/orders", 500, 240),
        ("GET", "/users", 200, 9),
        ("GET", "/users/1", 200, 7),
        ("GET", "/orders", 500, 260),
        ("DELETE", "/users/1", 204, 5),
    ]
    for method, path, status, ms in reqs:
        handle(metrics, method, path, status, ms)

    start = time.perf_counter()
    time.sleep(0.005)
    elapsed = (time.perf_counter() - start) * 1000

    requests, errors, avg, max_ms = metrics.snapshot()
    print("---- RED 指标 ----")
    print(
        f"requests={requests} errors={errors} "
        f"error_rate={errors / requests * 100:.2f}% avg_ms={avg:.1f} max_ms={max_ms}"
    )
    print(f"采集自身耗时={elapsed:.0f}ms（仅演示计时）")
    print("排查顺序：先看 error_rate 与 max_ms，再用 trace_id 串起同一次请求的全部日志")


if __name__ == "__main__":
    main()
