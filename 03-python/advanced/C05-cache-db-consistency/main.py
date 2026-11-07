"""缓存与 DB 双写一致性。"""


class Cache:
    def __init__(self) -> None:
        self.data: dict[str, str] = {}

    def get(self, k: str) -> str | None:
        return self.data.get(k)

    def set(self, k: str, v: str) -> None:
        self.data[k] = v

    def delete(self, k: str) -> None:
        self.data.pop(k, None)


def read_through(cache: Cache, db: dict[str, str], k: str) -> str:
    hit = cache.get(k)
    if hit is not None:
        return hit
    v = db[k]                 # 缓存未命中，回源
    cache.set(k, v)
    return v


def main() -> None:
    key = "user:1"

    print("== 1. 先删缓存，再写库 ==")
    db = {key: "v1"}
    cache = Cache()
    cache.set(key, "v1")
    cache.delete(key)                 # 写请求：先删缓存
    read_through(cache, db, key)      # 读请求：回源读到旧值
    db[key] = "v2"                    # 写请求：再写库
    print(f"缓存 = {cache.get(key)}，库 = {db[key]}，是否脏 = {cache.get(key) != db[key]}")

    print("== 2. 先写库，再删缓存 ==")
    db = {key: "v1"}
    cache = Cache()
    cache.set(key, "v1")
    old = cache.get(key)              # 读请求先拿到旧值
    db[key] = "v2"                    # 写请求：先写库
    cache.delete(key)                 # 再删缓存
    cache.set(key, old)               # 读请求迟到，把旧值写回
    print(f"缓存 = {cache.get(key)}，库 = {db[key]}，是否脏 = {cache.get(key) != db[key]}")

    print("== 3. 延迟双删 ==")
    db = {key: "v1"}
    cache = Cache()
    cache.set(key, "v1")
    old = cache.get(key)
    cache.delete(key)
    db[key] = "v2"
    cache.set(key, old)               # 迟到的读请求仍回填旧值
    cache.delete(key)                 # 延迟再删一次
    print(f"缓存命中 = {cache.get(key) is not None}，下一次读会回源拿到 {read_through(cache, db, key)}")

    print("== 4. 结论 ==")
    print("先删缓存有窗口，先写库再删缓存也有窗口，双删只能把窗口压小")
    print("要强一致就别用缓存读；要最终一致就配 binlog 订阅 + 重试删除")


if __name__ == "__main__":
    main()
