use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Arc, Mutex};

#[derive(Debug, Clone, Copy, PartialEq)]
enum State {
    Closed,
    Open,
    HalfOpen,
}

impl State {
    fn label(self) -> &'static str {
        match self {
            State::Closed => "closed",
            State::Open => "open",
            State::HalfOpen => "half-open",
        }
    }
}

struct Inner {
    state: State,
    failures: u32,
    threshold: u32,
    cooldown_ms: u64,
    opened_at: u64,
    probe_in_flight: bool,
}

struct Breaker {
    inner: Mutex<Inner>,
    clock: Arc<AtomicU64>,
}

impl Breaker {
    fn new(threshold: u32, cooldown_ms: u64) -> Self {
        Breaker {
            inner: Mutex::new(Inner {
                state: State::Closed,
                failures: 0,
                threshold,
                cooldown_ms,
                opened_at: 0,
                probe_in_flight: false,
            }),
            clock: Arc::new(AtomicU64::new(0)),
        }
    }

    fn advance(&self, ms: u64) {
        self.clock.fetch_add(ms, Ordering::SeqCst);
    }

    fn allow(&self) -> Result<(), &'static str> {
        let now = self.clock.load(Ordering::SeqCst);
        let mut inner = self.inner.lock().unwrap();
        match inner.state {
            State::Open => {
                if now.saturating_sub(inner.opened_at) >= inner.cooldown_ms {
                    inner.state = State::HalfOpen;
                    inner.probe_in_flight = true;
                    Ok(())
                } else {
                    Err("circuit open: fast fail")
                }
            }
            State::HalfOpen => {
                if inner.probe_in_flight {
                    Err("circuit half-open: probe already in flight")
                } else {
                    inner.probe_in_flight = true;
                    Ok(())
                }
            }
            State::Closed => Ok(()),
        }
    }

    fn success(&self) {
        let mut inner = self.inner.lock().unwrap();
        inner.failures = 0;
        inner.probe_in_flight = false;
        inner.state = State::Closed;
    }

    fn failure(&self) {
        let now = self.clock.load(Ordering::SeqCst);
        let mut inner = self.inner.lock().unwrap();
        inner.failures += 1;
        inner.probe_in_flight = false;
        if inner.state == State::HalfOpen || inner.failures >= inner.threshold {
            inner.state = State::Open;
            inner.opened_at = now;
        }
    }

    fn state(&self) -> State {
        self.inner.lock().unwrap().state
    }
}

fn call(b: &Breaker, downstream_ok: bool) -> String {
    if let Err(reason) = b.allow() {
        return format!("rejected: {}", reason);
    }
    if downstream_ok {
        b.success();
        "ok".to_string()
    } else {
        b.failure();
        "downstream error".to_string()
    }
}

fn main() {
    let b = Breaker::new(5, 200);

    println!("== 1. 正常请求 ==");
    println!("call -> {}, state={}", call(&b, true), b.state().label());

    println!("== 2. 连续 5 次失败触发熔断 ==");
    for i in 1..=5 {
        println!("failure #{} -> {}, state={}", i, call(&b, false), b.state().label());
    }

    println!("== 3. 熔断期间快速失败 ==");
    println!("call -> {}, state={}", call(&b, true), b.state().label());

    println!("== 4. 冷却 200ms 后进入半开 ==");
    b.advance(200);
    println!("allow -> {}, state={}", b.allow().is_ok(), b.state().label());

    println!("== 5. 探测成功恢复 closed ==");
    b.success();
    println!("state={}, 后续请求 -> {}", b.state().label(), call(&b, true));

    println!("== 6. 半开探测失败则重新 open ==");
    for _ in 0..5 {
        b.failure();
    }
    b.advance(200);
    let _ = b.allow();
    b.failure();
    println!("state={}（探测失败后立即回到 open）", b.state().label());
}
