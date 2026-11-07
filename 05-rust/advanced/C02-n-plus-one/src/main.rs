use std::collections::{BTreeSet, HashMap};

#[derive(Clone)]
struct Order {
    id: usize,
    account_id: usize,
    amount: i64,
}

struct FakeDb {
    orders: Vec<Order>,
    accounts: HashMap<usize, String>,
    queries: usize,
}

impl FakeDb {
    fn new() -> Self {
        let mut accounts = HashMap::new();
        accounts.insert(1, "alice".to_string());
        accounts.insert(2, "bob".to_string());
        accounts.insert(3, "carol".to_string());

        let orders = (1..=10)
            .map(|i| Order {
                id: i,
                account_id: i % 3 + 1,
                amount: 100 + i as i64,
            })
            .collect();

        Self {
            orders,
            accounts,
            queries: 0,
        }
    }

    fn list_orders(&mut self) -> Vec<Order> {
        self.queries += 1;
        self.orders.clone()
    }

    fn get_account(&mut self, id: usize) -> String {
        self.queries += 1;
        self.accounts.get(&id).cloned().unwrap_or_default()
    }

    fn get_accounts(&mut self, ids: &[usize]) -> HashMap<usize, String> {
        self.queries += 1;
        ids.iter()
            .filter_map(|id| self.accounts.get(id).map(|n| (*id, n.clone())))
            .collect()
    }
}

fn n_plus_one(db: &mut FakeDb) -> Vec<String> {
    let mut out = Vec::new();
    for o in db.list_orders() {
        // 1 次
        let name = db.get_account(o.account_id); // N 次
        out.push(format!("{}:{}", o.id, name));
    }
    out
}

fn batched(db: &mut FakeDb) -> Vec<String> {
    let orders = db.list_orders(); // 1 次
    let ids: Vec<usize> = orders
        .iter()
        .map(|o| o.account_id)
        .collect::<BTreeSet<_>>()
        .into_iter()
        .collect();
    let mapping = db.get_accounts(&ids); // 1 次
    orders
        .iter()
        .map(|o| format!("{}:{}", o.id, mapping.get(&o.account_id).cloned().unwrap_or_default()))
        .collect()
}

fn main() {
    let mut db1 = FakeDb::new();
    let a = n_plus_one(&mut db1);
    println!("N+1：{} 条订单 -> 查询 {} 次", a.len(), db1.queries);

    let mut db2 = FakeDb::new();
    let b = batched(&mut db2);
    println!("批量：{} 条订单 -> 查询 {} 次", b.len(), db2.queries);

    println!("结果一致：{}", a == b);
    println!("结论：N+1 的代价是 1+N 次往返，批量把 N 压成 1");
}
