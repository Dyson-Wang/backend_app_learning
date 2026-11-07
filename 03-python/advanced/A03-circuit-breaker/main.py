import threading
import time
from enum import Enum


class State(Enum):
    CLOSED = "closed"
    OPEN = "open"
    HALF_OPEN = "half-open"


class Breaker:
    """三态熔断器：连续失败到阈值就 open，冷却后放一个探测。"""

    def __init__(self, threshold=5, cooldown=0.2, clock=time.monotonic):
        self._lock = threading.Lock()
        self.state = State.CLOSED
        self.failures = 0
        self.threshold = threshold
        self.cooldown = cooldown
        self.opened_at = 0.0
        self.probe_in_flight = False
        self._clock = clock

    def allow(self):
        with self._lock:
            if self.state is State.OPEN:
                if self._clock() - self.opened_at >= self.cooldown:
                    self.state = State.HALF_OPEN
                    self.probe_in_flight = True
                    return True, ""
                return False, "circuit open: fast fail"
            if self.state is State.HALF_OPEN:
                if self.probe_in_flight:
                    return False, "circuit half-open: probe already in flight"
                self.probe_in_flight = True
                return True, ""
            return True, ""

    def success(self):
        with self._lock:
            self.failures = 0
            self.probe_in_flight = False
            self.state = State.CLOSED

    def failure(self):
        with self._lock:
            self.failures += 1
            self.probe_in_flight = False
            if self.state is State.HALF_OPEN or self.failures >= self.threshold:
                self.state = State.OPEN
                self.opened_at = self._clock()


def call(breaker, downstream_ok):
    allowed, reason = breaker.allow()
    if not allowed:
        return f"rejected: {reason}"
    if downstream_ok:
        breaker.success()
        return "ok"
    breaker.failure()
    return "downstream error"


class FakeClock:
    def __init__(self):
        self.t = 0.0

    def __call__(self):
        return self.t

    def advance(self, seconds):
        self.t += seconds


def main():
    clock = FakeClock()
    b = Breaker(threshold=5, cooldown=0.2, clock=clock)

    print("== 1. 正常请求 ==")
    print(f"call -> {call(b, True)}, state={b.state.value}")

    print("== 2. 连续 5 次失败触发熔断 ==")
    for i in range(1, 6):
        print(f"failure #{i} -> {call(b, False)}, state={b.state.value}")

    print("== 3. 熔断期间快速失败 ==")
    print(f"call -> {call(b, True)}, state={b.state.value}")

    print("== 4. 冷却 200ms 后进入半开 ==")
    clock.advance(0.2)
    print(f"allow -> {b.allow()}, state={b.state.value}")

    print("== 5. 探测成功恢复 closed ==")
    b.success()
    print(f"state={b.state.value}, 后续请求 -> {call(b, True)}")

    print("== 6. 半开探测失败则重新 open ==")
    for _ in range(5):
        b.failure()
    clock.advance(0.2)
    b.allow()
    b.failure()
    print(f"state={b.state.value}（探测失败后立即回到 open）")


if __name__ == "__main__":
    main()
