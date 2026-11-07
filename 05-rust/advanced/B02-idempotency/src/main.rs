use std::collections::{HashMap, HashSet};
use std::sync::{Arc, Condvar, Mutex};
use std::thread;
use std::time::Duration;

#[derive(Clone, Debug)]
struct Order {
    id: usize,
    key: String,
    amount: i64,
}

#[derive(Default)]
struct Inner {
    done: HashMap<String, Order>,
    busy: HashSet<String>,
    created: usize,
    replayed: usize,
    seq: usize,
}

struct IdemStore {
    inner: Mutex<Inner>,
    cv: Condvar,
}

impl IdemStore {
    fn new() -> Self {
        Self {
            inner: Mutex::new(Inner::default()),
            cv: Condvar::new(),
        }
    }

    fn create(&self, key: &str, amount: i64) -> (Order, bool) {
        let mut g = self.inner.lock().unwrap();
        if let Some(o) = g.done.get(key) {
            g.replayed += 1;
            return (o.clone(), true);
        }
        // 同 key 已有线程在建单：等它完成，直接拿结果
        while g.busy.contains(key) {
            g = self.cv.wait(g).unwrap();
            if let Some(o) = g.done.get(key) {
                g.replayed += 1;
                return (o.clone(), true);
            }
        }
        g.busy.insert(key.to_string());
        drop(g);                                  // 建单期间不持锁

        thread::sleep(Duration::from_millis(20)); // 模拟写库

        let mut g = self.inner.lock().unwrap();
        g.seq += 1;
        let order = Order {
            id: g.seq,
            key: key.to_string(),
            amount,
        };
        g.done.insert(key.to_string(), order.clone());
        g.created += 1;
        g.busy.remove(key);
        self.cv.notify_all();
        (order, false)
    }
}

fn main() {
    let store = Arc::new(IdemStore::new());
    let mut handles = Vec::new();

    println!("== 1. 20 个并发请求，同一个 Idempotency-Key ==");
    for _ in 0..20 {
        let s = Arc::clone(&store);
        handles.push(thread::spawn(move || s.create("idem-7f3a", 9900)));
    }
    let results: Vec<(Order, bool)> = handles.into_iter().map(|h| h.join().unwrap()).collect();

    let created = results.iter().filter(|(_, replayed)| !*replayed).count();
    let replayed = results.iter().filter(|(_, replayed)| *replayed).count();
    let ids: HashSet<usize> = results.iter().map(|(o, _)| o.id).collect();
    println!(
        "真实建单 = {}，回放 = {}，订单 id 去重后 = {}",
        created,
        replayed,
        ids.len()
    );

    println!("== 2. 换一个 key ==");
    let (o2, replayed2) = store.create("idem-other", 100);
    println!("新 key 建单 id = {}，replayed = {}", o2.id, replayed2);

    println!("== 3. 结论 ==");
    println!("幂等键必须落库并加唯一索引；Mutex<HashMap> 只是演示");
    println!("Condvar 做「等待别人建完」，比忙等更省 CPU");
}
