#[derive(Debug, Clone)]
struct Account {
    balance: i64,
    version: u64,
}

impl Account {
    fn new(balance: i64) -> Self {
        Self { balance, version: 0 }
    }

    /// 乐观锁：只有 version 没被别人改过才写入
    fn cas(&mut self, expected_version: u64, new_balance: i64) -> bool {
        if self.version != expected_version {
            return false;
        }
        self.balance = new_balance;
        self.version += 1;
        true
    }
}

fn debit(acc: &mut Account, amount: i64) -> Option<u32> {
    for attempt in 0..=10u32 {
        let snapshot = acc.clone();               // 事务开始时的读视图
        if acc.cas(snapshot.version, snapshot.balance - amount) {
            return Some(attempt);
        }
    }
    None
}

fn main() {
    println!("== 1. 无锁：两个事务都基于旧值写回 ==");
    let mut naive = Account::new(100);
    let read1 = naive.balance;
    let read2 = naive.balance;
    naive.balance = read1 - 10;                   // 两个事务都用过期的 100
    naive.balance = read2 - 10;
    println!("期望 80，实际 {}（丢失更新）", naive.balance);

    println!("== 2. 乐观锁：version 校验 + 重试 ==");
    let mut acc = Account::new(100);
    let a1 = debit(&mut acc, 10);
    println!("T1 attempts={:?}，balance={} version={}", a1, acc.balance, acc.version);

    let stale = acc.version - 1;
    println!("T2 用旧 version={} 写入：{}", stale, acc.cas(stale, 70));

    let a2 = debit(&mut acc, 10);
    println!("T2 重试 attempts={:?}，balance={} version={}", a2, acc.balance, acc.version);
    println!("期望 80，实际 {}", acc.balance);

    println!("== 3. 结论 ==");
    println!("悲观锁锁整行，吞吐低；乐观锁读不加锁、写校验 version，冲突才重试");
    println!("Rust 里 borrow checker 能挡住「持有引用时又可变借用」，但挡不住并发丢更新");
}
