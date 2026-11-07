use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::{Duration, Instant};

struct FakeDb {
    hits: Mutex<u64>,
}

impl FakeDb {
    fn new() -> Self {
        FakeDb {
            hits: Mutex::new(0),
        }
    }

    fn query(&self, key: &str) -> String {
        *self.hits.lock().unwrap() += 1;
        thread::sleep(Duration::from_millis(20));
        format!("value-of-{}", key)
    }

    fn hits(&self) -> u64 {
        *self.hits.lock().unwrap()
    }
}

fn unprotected_get(db: &FakeDb, cache: &Mutex<HashMap<String, String>>, key: &str) -> String {
    {
        let data = cache.lock().unwrap();
        if let Some(v) = data.get(key) {
            return v.clone();
        }
    }
    let value = db.query(key); // 并发请求全部走到这里
    cache.lock().unwrap().insert(key.to_string(), value.clone());
    value
}

struct GuardedCache {
    data: Mutex<HashMap<String, String>>,
    key_locks: Mutex<HashMap<String, Arc<Mutex<()>>>>,
}

impl GuardedCache {
    fn new() -> Self {
        GuardedCache {
            data: Mutex::new(HashMap::new()),
            key_locks: Mutex::new(HashMap::new()),
        }
    }

    fn get(&self, db: &FakeDb, key: &str) -> String {
        {
            let data = self.data.lock().unwrap();
            if let Some(v) = data.get(key) {
                return v.clone();
            }
        }

        let lock = {
            let mut map = self.key_locks.lock().unwrap();
            Arc::clone(
                map.entry(key.to_string())
                    .or_insert_with(|| Arc::new(Mutex::new(()))),
            )
        };

        let _guard = lock.lock().unwrap();

        {
            let data = self.data.lock().unwrap();
            if let Some(v) = data.get(key) {
                return v.clone();
            }
        }

        let value = db.query(key);
        self.data
            .lock()
            .unwrap()
            .insert(key.to_string(), value.clone());
        value
    }
}

fn run(n: usize, guarded: bool) -> (u64, u128) {
    let db = Arc::new(FakeDb::new());
    let start = Instant::now();

    let mut handles = Vec::with_capacity(n);
    if guarded {
        let cache = Arc::new(GuardedCache::new());
        for _ in 0..n {
            let db = Arc::clone(&db);
            let cache = Arc::clone(&cache);
            handles.push(thread::spawn(move || {
                cache.get(&db, "user:1");
            }));
        }
    } else {
        let cache = Arc::new(Mutex::new(HashMap::new()));
        for _ in 0..n {
            let db = Arc::clone(&db);
            let cache = Arc::clone(&cache);
            handles.push(thread::spawn(move || {
                unprotected_get(&db, &cache, "user:1");
            }));
        }
    }

    for h in handles {
        h.join().unwrap();
    }

    (db.hits(), start.elapsed().as_millis())
}

fn main() {
    let n = 50;

    let (hits, ms) = run(n, false);
    println!("无保护：{} 并发 -> 回源 {} 次，耗时 {}ms", n, hits, ms);

    let (hits, ms) = run(n, true);
    println!(
        "按 key 加锁 + 双重检查：{} 并发 -> 回源 {} 次，耗时 {}ms",
        n, hits, ms
    );

    println!("结论：缓存击穿是「同一热 key 失效瞬间的并发回源」，无锁读写挡不住");
}
