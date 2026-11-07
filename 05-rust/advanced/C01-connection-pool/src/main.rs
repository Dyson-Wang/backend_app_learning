use std::sync::{Arc, Condvar, Mutex};
use std::thread;
use std::time::{Duration, Instant};

struct Pool {
    size: usize,
    avail: Mutex<usize>,
    cv: Condvar,
    peak: Mutex<usize>,
    in_use: Mutex<usize>,
}

impl Pool {
    fn new(size: usize) -> Self {
        Self {
            size,
            avail: Mutex::new(size),
            cv: Condvar::new(),
            peak: Mutex::new(0),
            in_use: Mutex::new(0),
        }
    }

    fn acquire(&self, wait: Option<Duration>) -> bool {
        let mut avail = self.avail.lock().unwrap();
        if *avail == 0 {
            match wait {
                None => {
                    while *avail == 0 {
                        avail = self.cv.wait(avail).unwrap();
                    }
                }
                Some(d) => {
                    let deadline = Instant::now() + d;
                    while *avail == 0 {
                        let now = Instant::now();
                        if now >= deadline {
                            return false;
                        }
                        let (g, timeout) = self.cv.wait_timeout(avail, deadline - now).unwrap();
                        avail = g;
                        if timeout.timed_out() && *avail == 0 {
                            return false;
                        }
                    }
                }
            }
        }
        *avail -= 1;
        drop(avail);

        let mut n = self.in_use.lock().unwrap();
        *n += 1;
        let mut p = self.peak.lock().unwrap();
        if *n > *p {
            *p = *n;
        }
        true
    }

    fn release(&self) {
        *self.in_use.lock().unwrap() -= 1;
        *self.avail.lock().unwrap() += 1;
        self.cv.notify_one();
    }

    fn peak(&self) -> usize {
        *self.peak.lock().unwrap()
    }

    fn size(&self) -> usize {
        self.size
    }
}

fn run(pool: Arc<Pool>, n: usize, hold_ms: u64, wait: Option<Duration>) -> (usize, usize, u128) {
    let ok = Arc::new(Mutex::new(0usize));
    let timeout = Arc::new(Mutex::new(0usize));
    let t0 = Instant::now();
    let mut handles = Vec::new();

    for _ in 0..n {
        let pool = Arc::clone(&pool);
        let ok = Arc::clone(&ok);
        let timeout = Arc::clone(&timeout);
        handles.push(thread::spawn(move || {
            if pool.acquire(wait) {
                thread::sleep(Duration::from_millis(hold_ms));
                *ok.lock().unwrap() += 1;
                pool.release();
            } else {
                *timeout.lock().unwrap() += 1;
            }
        }));
    }
    for h in handles {
        h.join().unwrap();
    }
    let cost = t0.elapsed().as_millis();
    let o = *ok.lock().unwrap();
    let t = *timeout.lock().unwrap();
    (o, t, cost)
}

fn main() {
    println!("== 1. 池大小 3，10 个请求，无等待超时 ==");
    let p1 = Arc::new(Pool::new(3));
    let (ok1, to1, cost1) = run(Arc::clone(&p1), 10, 120, None);
    println!(
        "ok={} timeout={} 耗时={}ms peak={}",
        ok1,
        to1,
        cost1,
        p1.peak()
    );

    println!("== 2. 同样的池，等待超过 40ms 就快速失败 ==");
    let p2 = Arc::new(Pool::new(3));
    let (ok2, to2, cost2) = run(Arc::clone(&p2), 10, 120, Some(Duration::from_millis(40)));
    println!(
        "ok={} timeout={} 耗时={}ms peak={}",
        ok2,
        to2,
        cost2,
        p2.peak()
    );

    println!("== 3. 结论 ==");
    println!("连接池/并发上限必须配等待超时：没有超时，故障时请求无限排队");
    println!("sqlx 的 acquire_timeout、bb8 的 connection_timeout 都是这个作用");
    println!("（本示例池大小 {}，同一时刻最多 3 个在跑）", p1.size());
}
