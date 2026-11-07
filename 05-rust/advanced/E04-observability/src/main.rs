use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::Mutex;

struct Metrics {
    requests: u64,
    errors: u64,
    total_ms: u64,
    max_ms: u64,
}

impl Metrics {
    fn new() -> Self {
        Metrics {
            requests: 0,
            errors: 0,
            total_ms: 0,
            max_ms: 0,
        }
    }

    fn observe(&mut self, status: u16, ms: u64) {
        self.requests += 1;
        if status >= 500 {
            self.errors += 1;
        }
        self.total_ms += ms;
        if ms > self.max_ms {
            self.max_ms = ms;
        }
    }
}

static TRACE_SEQ: AtomicU64 = AtomicU64::new(0);

fn next_trace_id() -> String {
    let n = TRACE_SEQ.fetch_add(1, Ordering::SeqCst) + 1;
    format!("trace-{:04}", n)
}

fn log_line(level: &str, trace_id: &str, method: &str, path: &str, status: u16, ms: u64) {
    let msg = if status >= 500 {
        r#","msg":"upstream failed""#
    } else {
        ""
    };
    println!(
        r#"{{"level":"{}","trace_id":"{}","method":"{}","path":"{}","status":{},"latency_ms":{}{}}}"#,
        level, trace_id, method, path, status, ms, msg
    );
}

fn handle(metrics: &Mutex<Metrics>, method: &str, path: &str, status: u16, ms: u64) {
    let trace_id = next_trace_id();
    metrics.lock().unwrap().observe(status, ms);
    let level = if status >= 500 { "error" } else { "info" };
    log_line(level, &trace_id, method, path, status, ms);
}

fn main() {
    let metrics = Mutex::new(Metrics::new());

    let reqs: [(&str, &str, u16, u64); 10] = [
        ("GET", "/healthz", 200, 1),
        ("GET", "/users", 200, 12),
        ("POST", "/users", 201, 18),
        ("GET", "/users/999", 404, 3),
        ("POST", "/users", 400, 2),
        ("GET", "/orders", 500, 240),
        ("GET", "/users", 200, 9),
        ("GET", "/users/1", 200, 7),
        ("GET", "/orders", 500, 260),
        ("DELETE", "/users/1", 204, 5),
    ];

    for (method, path, status, ms) in reqs {
        handle(&metrics, method, path, status, ms);
    }

    let m = metrics.lock().unwrap();
    let avg = if m.requests == 0 {
        0.0
    } else {
        m.total_ms as f64 / m.requests as f64
    };

    println!("---- RED 指标 ----");
    println!(
        "requests={} errors={} error_rate={:.2}% avg_ms={:.1} max_ms={}",
        m.requests,
        m.errors,
        m.errors as f64 / m.requests as f64 * 100.0,
        avg,
        m.max_ms
    );
    println!("排查顺序：先看 error_rate 与 max_ms，再用 trace_id 串起同一次请求的全部日志");
}
