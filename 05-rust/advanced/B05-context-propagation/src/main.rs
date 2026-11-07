use std::thread;
use std::time::{Duration, Instant};

struct Deadline {
    at: Instant,
}

impl Deadline {
    fn new(budget_ms: u64) -> Self {
        Self {
            at: Instant::now() + Duration::from_millis(budget_ms),
        }
    }

    fn remaining_ms(&self) -> i64 {
        let now = Instant::now();
        if self.at > now {
            (self.at - now).as_millis() as i64
        } else {
            0
        }
    }
}

/// deadline 为 None 表示中间层没传，下游完全不看上游预算
fn downstream(deadline: Option<&Deadline>, work_ms: u64) -> &'static str {
    match deadline {
        None => {
            thread::sleep(Duration::from_millis(work_ms));
            "ok"
        }
        Some(d) => {
            let left = d.remaining_ms();
            if (left as u64) < work_ms {
                if left > 0 {
                    thread::sleep(Duration::from_millis(left as u64));   // 跑到截止就放弃
                }
                "deadline exceeded"
            } else {
                thread::sleep(Duration::from_millis(work_ms));
                "ok"
            }
        }
    }
}

fn main() {
    println!("== 1. 上游预算 50ms，下游要 200ms ==");
    let d = Deadline::new(50);
    let t0 = Instant::now();
    let r1 = downstream(Some(&d), 200);
    println!(
        "res = {}，耗时 = {}ms（远小于 200ms）",
        r1,
        t0.elapsed().as_millis()
    );

    println!("== 2. 中间层丢掉 deadline ==");
    let t1 = Instant::now();
    let r2 = downstream(None, 200);
    println!(
        "res = {}，耗时 = {}ms（上游再急也没用）",
        r2,
        t1.elapsed().as_millis()
    );

    println!("== 3. 结论 ==");
    println!("deadline 必须一路显式传；Rust 里没有隐式上下文，忘记传就是编译期可见的类型缺失");
    println!("生产上用 tokio::time::timeout，跨服务用 X-Request-Deadline 头");
}
