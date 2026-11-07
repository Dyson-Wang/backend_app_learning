use std::collections::HashSet;

#[derive(Clone)]
struct Event {
    id: &'static str,
    amount: i32,
}

/// seen 为 None 表示消费端没有幂等表
fn handle(seen: &mut Option<HashSet<&'static str>>, balance: &mut i32, e: &Event) -> bool {
    if let Some(s) = seen {
        if s.contains(e.id) {
            return false; // 重复消息，直接 ACK
        }
        s.insert(e.id);
    }
    *balance += e.amount;
    true
}

fn main() {
    // e2 被投递两次（消费者 ACK 丢了）
    let deliveries = [
        Event { id: "e1", amount: 10 },
        Event { id: "e2", amount: 10 },
        Event { id: "e2", amount: 10 },
        Event { id: "e3", amount: 10 },
    ];

    println!("== 1. 消费端没有幂等表 ==");
    let mut balance = 100;
    let mut none: Option<HashSet<&'static str>> = None;
    for e in &deliveries {
        handle(&mut none, &mut balance, e);
    }
    println!(
        "投递 {} 次，余额 100 -> {}（多加了 10）",
        deliveries.len(),
        balance
    );

    println!("== 2. 消费端按消息 ID 去重 ==");
    let mut balance = 100;
    let mut seen = Some(HashSet::new());
    let mut applied = 0;
    for e in &deliveries {
        if handle(&mut seen, &mut balance, e) {
            applied += 1;
        }
    }
    println!(
        "投递 {} 次，实际处理 {} 次，余额 100 -> {}",
        deliveries.len(),
        applied,
        balance
    );

    println!("== 3. 结论 ==");
    println!("至少一次投递下，重复是常态；正确性必须由消费端幂等保证");
    println!("去重表要落库并和业务写在同一事务里，否则还是会重复");
}
