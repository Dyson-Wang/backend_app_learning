struct Msg {
    id: &'static str,
    failures: u32,   // 前 N 次投递失败（可重试）
    permanent: bool, // 永久错误，重试也没用
}

fn deliver(m: &Msg, max_attempts: u32) -> (bool, u32) {
    for attempt in 1..=max_attempts {
        if m.permanent {
            println!("{} -> 永久错误，直接进死信（不重试）", m.id);
            return (false, attempt);
        }
        if attempt > m.failures {
            println!("{} -> 成功（第 {} 次）", m.id, attempt);
            return (true, attempt);
        }
        if attempt < max_attempts {
            let backoff = 10u32 << (attempt - 1);
            println!("{} -> 可重试错误，退避 {}ms 后重试", m.id, backoff);
            continue;
        }
        println!("{} -> 重试耗尽，进死信", m.id);
    }
    (false, max_attempts)
}

fn main() {
    let msgs = [
        Msg { id: "m1", failures: 0, permanent: false },
        Msg { id: "m2", failures: 2, permanent: false },
        Msg { id: "m3", failures: 0, permanent: true },
        Msg { id: "m4", failures: 99, permanent: false },
        Msg { id: "m5", failures: 0, permanent: false },
    ];

    println!("== 1. 逐条投递（最多 3 次） ==");
    let mut ok = 0;
    let mut dead = 0;
    let mut total = 0;
    for m in &msgs {
        let (success, attempts) = deliver(m, 3);
        total += attempts;
        if success {
            ok += 1;
        } else {
            dead += 1;
        }
    }

    println!("== 2. 统计 ==");
    println!("成功 = {}，死信 = {}，总投递 = {}", ok, dead, total);

    println!("== 3. 结论 ==");
    println!("永久错误不要重试，重试只会拖慢整个队列");
    println!("死信队列必须有人看：它进 DLQ 说明要么代码有 bug，要么数据有问题");
}
