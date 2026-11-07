"""游标分页 vs OFFSET 分页。"""
from dataclasses import dataclass


@dataclass
class Item:
    id: int
    score: int


def ordered(items: list[Item]) -> list[Item]:
    return sorted(items, key=lambda i: (-i.score, i.id))


def offset_page(items: list[Item], offset: int, limit: int) -> list[Item]:
    return ordered(items)[offset:offset + limit]


def cursor_page(items: list[Item], after: Item | None, limit: int) -> list[Item]:
    out: list[Item] = []
    for it in ordered(items):
        if after is not None and not ((-it.score, it.id) > (-after.score, after.id)):
            continue
        out.append(it)
        if len(out) == limit:
            break
    return out


def ids(items: list[Item]) -> list[int]:
    return [i.id for i in items]


def overlap(a: list[Item], b: list[Item]) -> int:
    seen = {i.id for i in a}
    return sum(1 for i in b if i.id in seen)


def main() -> None:
    base = [Item(i + 1, 100 - i * 5) for i in range(10)]

    print("== 1. OFFSET 分页：翻页期间有新数据插入 ==")
    p1 = offset_page(base, 0, 3)
    grown = base + [Item(11, 97)]
    p2 = offset_page(grown, 3, 3)
    print(f"第一页 = {ids(p1)}")
    print(f"插入一条 Score=97 后，第二页 = {ids(p2)}，与第一页重复 {overlap(p1, p2)} 条")

    print("== 2. 游标分页：用 (score, id) 做锚点 ==")
    c1 = cursor_page(base, None, 3)
    last = c1[-1]
    c2 = cursor_page(grown, last, 3)
    print(f"第一页 = {ids(c1)}，锚点 = (score={last.score}, id={last.id})")
    print(f"插入同一条后，第二页 = {ids(c2)}，与第一页重复 {overlap(c1, c2)} 条")

    print("== 3. 结论 ==")
    print("OFFSET 按位置翻页，前面插入数据会让整体后移，出现重复或漏读")
    print("游标按值翻页，锚点不变结果就不变；代价是不能跳页")


if __name__ == "__main__":
    main()
