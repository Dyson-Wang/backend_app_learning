"""N+1 查询与批量写。"""
from dataclasses import dataclass


@dataclass
class Order:
    id: int
    account_id: int
    amount: int


@dataclass
class Account:
    id: int
    name: str


class FakeDB:
    def __init__(self) -> None:
        self.orders = [Order(i, i % 3 + 1, 100 + i) for i in range(1, 11)]
        self.accounts = {1: Account(1, "alice"), 2: Account(2, "bob"), 3: Account(3, "carol")}
        self.queries = 0

    def list_orders(self) -> list[Order]:
        self.queries += 1
        return list(self.orders)

    def get_account(self, account_id: int) -> Account:
        self.queries += 1
        return self.accounts[account_id]

    def get_accounts(self, ids: list[int]) -> dict[int, Account]:
        self.queries += 1
        return {i: self.accounts[i] for i in ids}


def n_plus_one(db: FakeDB) -> list[tuple[int, str]]:
    out = []
    for o in db.list_orders():            # 1 次
        acc = db.get_account(o.account_id)  # N 次
        out.append((o.id, acc.name))
    return out


def batched(db: FakeDB) -> list[tuple[int, str]]:
    orders = db.list_orders()             # 1 次
    ids = sorted({o.account_id for o in orders})
    mapping = db.get_accounts(ids)        # 1 次
    return [(o.id, mapping[o.account_id].name) for o in orders]


def main() -> None:
    db = FakeDB()
    a = n_plus_one(db)
    print(f"N+1：{len(a)} 条订单 -> 查询 {db.queries} 次")

    db2 = FakeDB()
    b = batched(db2)
    print(f"批量：{len(b)} 条订单 -> 查询 {db2.queries} 次")

    print(f"结果一致：{a == b}")
    print("结论：N+1 的代价是 1+N 次往返，批量把 N 压成 1")


if __name__ == "__main__":
    main()
