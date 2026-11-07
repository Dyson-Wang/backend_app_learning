"""分布式锁：SET NX PX + 校验 token 解锁。"""
import time


class FakeRedis:
    """用字典模拟 Redis 的 SET NX PX / GET / DEL。"""

    def __init__(self) -> None:
        self._data: dict[str, tuple[str, float]] = {}

    def set_nx_px(self, key: str, token: str, ttl_ms: int) -> bool:
        now = time.monotonic()
        cur = self._data.get(key)
        if cur is not None and cur[1] > now:
            return False
        self._data[key] = (token, now + ttl_ms / 1000)
        return True

    def get(self, key: str) -> str:
        cur = self._data.get(key)
        if cur is None or cur[1] <= time.monotonic():
            return ""
        return cur[0]

    def delete(self, key: str) -> None:
        self._data.pop(key, None)

    def del_if_match(self, key: str, token: str) -> bool:
        # 真实环境用 Lua 保证「比较 + 删除」原子
        if self.get(key) == token:
            self._data.pop(key, None)
            return True
        return False


def main() -> None:
    r = FakeRedis()
    key = "lock:order:42"

    print("== 1. A 加锁 ==")
    print(f"A SetNX = {r.set_nx_px(key, 'token-A', 50)}, holder={r.get(key)}")

    print("== 2. B 尝试加锁（A 还持有时） ==")
    print(f"B SetNX = {r.set_nx_px(key, 'token-B', 50)}")

    print("== 3. A 的业务超时，锁过期 ==")
    time.sleep(0.06)
    print(f"锁是否还有效 = {bool(r.get(key))}")

    print("== 4. B 抢到锁 ==")
    print(f"B SetNX = {r.set_nx_px(key, 'token-B', 50)}, holder={r.get(key)}")

    print("== 5. A 醒来直接 DEL：删掉了 B 的锁（错误） ==")
    r.delete(key)
    print(f'holder="{r.get(key)}"（B 的锁被误删）')

    print("== 6. 修正：delIfMatch 校验 token ==")
    r.set_nx_px(key, "token-B", 50)
    print(f"A delIfMatch(token-A) = {r.del_if_match(key, 'token-A')}, holder={r.get(key)}")
    print(f"B delIfMatch(token-B) = {r.del_if_match(key, 'token-B')}, holder={r.get(key)}")

    print("== 7. 结论 ==")
    print("加锁 SET key token NX PX ttl；解锁必须比较 token 再删，且要原子")
    print("业务可能超过 ttl 时要续期，否则锁提前释放")


if __name__ == "__main__":
    main()
