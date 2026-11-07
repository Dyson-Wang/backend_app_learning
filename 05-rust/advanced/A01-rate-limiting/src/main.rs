use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::Arc;

struct TokenBucket {
    rate: f64,
    burst: f64,
    tokens: f64,
    last_ms: u64,
    clock: Arc<AtomicU64>,
}

impl TokenBucket {
    fn new(rate: f64, burst: f64) -> Self {
        TokenBucket {
            rate,
            burst,
            tokens: burst,
            last_ms: 0,
            clock: Arc::new(AtomicU64::new(0)),
        }
    }

    fn advance(&self, ms: u64) {
        self.clock.fetch_add(ms, Ordering::SeqCst);
    }

    fn allow(&mut self) -> bool {
        let now = self.clock.load(Ordering::SeqCst);
        let elapsed_sec = (now - self.last_ms) as f64 / 1000.0;
        self.tokens = self.burst.min(self.tokens + elapsed_sec * self.rate);
        self.last_ms = now;
        if self.tokens >= 1.0 {
            self.tokens -= 1.0;
            true
        } else {
            false
        }
    }
}

fn main() {
    let mut bucket = TokenBucket::new(2.0, 5.0);

    println!("== t=0ms：burst=5，连续 6 次 ==");
    for i in 1..=6 {
        let ok = bucket.allow();
        println!("  #{} -> allow={}, tokens={:.2}", i, ok, bucket.tokens);
    }

    println!("== t=500ms：补 1 个令牌 ==");
    bucket.advance(500);
    for i in 7..=8 {
        let ok = bucket.allow();
        println!("  #{} -> allow={}, tokens={:.2}", i, ok, bucket.tokens);
    }

    println!("== t=1500ms：补 2 个令牌 ==");
    bucket.advance(1000);
    for i in 9..=11 {
        let ok = bucket.allow();
        println!("  #{} -> allow={}, tokens={:.2}", i, ok, bucket.tokens);
    }

    println!("结论：令牌桶允许突发（burst），但长期速率被 rate 卡死");
}
