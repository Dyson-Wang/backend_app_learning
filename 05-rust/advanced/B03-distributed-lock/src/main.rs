use std::collections::HashMap;
use std::sync::Mutex;
use std::thread;
use std::time::{Duration, Instant};

#[derive(Clone)]
struct Entry {
    token: String,
    expire_at: Instant,
}

// 用 HashMap 模拟 Redis 的 SET NX PX / GET / DEL
struct FakeRedis {
    data: Mutex<HashMap<String, Entry>>,
}

impl FakeRedis {
    fn new() -> Self {
        Self {
            data: Mutex::new(HashMap::new()),
        }
    }

    fn set_nx_px(&self, key: &str, token: &str, ttl_ms: u64) -> bool {
        let mut g = self.data.lock().unwrap();
        let now = Instant::now();
        if let Some(cur) = g.get(key) {
            if cur.expire_at > now {
                return false;
            }
        }
        g.insert(
            key.to_string(),
            Entry {
                token: token.to_string(),
                expire_at: now + Duration::from_millis(ttl_ms),
            },
        );
        true
    }

    fn get(&self, key: &str) -> String {
        let g = self.data.lock().unwrap();
        match g.get(key) {
            Some(e) if e.expire_at > Instant::now() => e.token.clone(),
            _ => String::new(),
        }
    }

    fn delete(&self, key: &str) {
        self.data.lock().unwrap().remove(key);
    }

    // 真实环境用 Lua 保证「比较 + 删除」原子
    fn del_if_match(&self, key: &str, token: &str) -> bool {
        let mut g = self.data.lock().unwrap();
        match g.get(key) {
            Some(e) if e.token == token && e.expire_at > Instant::now() => {
                g.remove(key);
                true
            }
            _ => false,
        }
    }
}

fn main() {
    let r = FakeRedis::new();
    let key = "lock:order:42";

    println!("== 1. A 加锁 ==");
    println!("A SetNX = {}，holder={}", r.set_nx_px(key, "token-A", 50), r.get(key));

    println!("== 2. B 尝试加锁（A 还持有时） ==");
    println!("B SetNX = {}", r.set_nx_px(key, "token-B", 50));

    println!("== 3. A 的业务超时，锁过期 ==");
    thread::sleep(Duration::from_millis(60));
    println!("锁是否还有效 = {}", !r.get(key).is_empty());

    println!("== 4. B 抢到锁 ==");
    println!("B SetNX = {}，holder={}", r.set_nx_px(key, "token-B", 50), r.get(key));

    println!("== 5. A 醒来直接 DEL：删掉了 B 的锁（错误） ==");
    r.delete(key);
    println!("holder=\"{}\"（B 的锁被误删）", r.get(key));

    println!("== 6. 修正：del_if_match 校验 token ==");
    r.set_nx_px(key, "token-B", 50);
    println!(
        "A del_if_match(token-A) = {}，holder={}",
        r.del_if_match(key, "token-A"),
        r.get(key)
    );
    println!(
        "B del_if_match(token-B) = {}，holder={}",
        r.del_if_match(key, "token-B"),
        r.get(key)
    );

    println!("== 7. 结论 ==");
    println!("加锁 SET key token NX PX ttl；解锁必须比较 token 再删，且要原子");
    println!("业务可能超过 ttl 时要续期，否则锁提前释放");
}
