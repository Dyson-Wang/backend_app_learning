#[derive(Clone)]
struct Event {
    id: &'static str,
    kind: &'static str,
}

fn apply(state: &str, e: &Event) -> (String, String) {
    let ok = match (e.kind, state) {
        ("create", "") => Some("CREATED"),
        ("pay", "CREATED") => Some("PAID"),
        ("ship", "PAID") => Some("SHIPPED"),
        _ => None,
    };
    match ok {
        Some(next) => (next.to_string(), String::new()),
        None => {
            let shown = if state.is_empty() { "INIT" } else { state };
            (
                state.to_string(),
                format!("{} 在 {} 状态下不合法", e.kind, shown),
            )
        }
    }
}

fn consume(events: &[Event]) -> (String, String) {
    let mut state = String::new();
    for e in events {
        let (next, err) = apply(&state, e);
        state = next;
        if !err.is_empty() {
            return (state, err);
        }
    }
    (state, String::new())
}

fn partition(key: &str, n: u32) -> u32 {
    let mut h: u32 = 0;
    for ch in key.chars() {
        h = (h.wrapping_mul(31).wrapping_add(ch as u32)) % n;
    }
    h
}

fn main() {
    let create = Event { id: "c1", kind: "create" };
    let pay = Event { id: "c2", kind: "pay" };
    let ship = Event { id: "c3", kind: "ship" };

    println!("== 1. 同一 key 被并发消费（顺序打乱） ==");
    let (state, err) = consume(&[ship.clone(), create.clone(), pay.clone()]);
    let shown = if state.is_empty() { "INIT" } else { &state };
    println!("最终状态 = {}，错误 = {}", shown, err);

    println!("== 2. 按 key 路由到同一分区、单消费者串行处理 ==");
    let (state, err) = consume(&[create, pay, ship]);
    println!(
        "最终状态 = {}，错误 = {}",
        state,
        if err.is_empty() { "无" } else { &err }
    );

    println!("== 3. 分区函数：同一 order_id 必然落到同一分区 ==");
    for k in ["order-42", "order-43", "order-42"] {
        println!("{} -> 分区 {}", k, partition(k, 4));
    }

    println!("== 4. 结论 ==");
    println!("顺序性不是全局的，是「同一分区内」的；跨 key 的顺序本来就没有意义");
    println!("要顺序就牺牲并发：单分区单消费者，或者消费端按 key 做串行队列");
}
