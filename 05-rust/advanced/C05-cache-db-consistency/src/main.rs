use std::collections::HashMap;

#[derive(Default)]
struct Cache {
    data: HashMap<String, String>,
}

impl Cache {
    fn get(&self, k: &str) -> Option<String> {
        self.data.get(k).cloned()
    }
    fn set(&mut self, k: &str, v: &str) {
        self.data.insert(k.to_string(), v.to_string());
    }
    fn delete(&mut self, k: &str) {
        self.data.remove(k);
    }
}

fn read_through(cache: &mut Cache, db: &HashMap<String, String>, k: &str) -> String {
    if let Some(v) = cache.get(k) {
        return v;
    }
    let v = db.get(k).cloned().unwrap_or_default();   // 缓存未命中，回源
    cache.set(k, &v);
    v
}

fn main() {
    let key = "user:1";
    let k = key;

    println!("== 1. 先删缓存，再写库 ==");
    let mut db: HashMap<String, String> = HashMap::new();
    db.insert(k.to_string(), "v1".to_string());
    let mut cache = Cache::default();
    cache.set(k, "v1");
    cache.delete(k);                       // 写请求：先删缓存
    read_through(&mut cache, &db, k);      // 读请求：回源读到旧值
    db.insert(k.to_string(), "v2".to_string());   // 写请求：再写库
    let got = cache.get(k).unwrap_or_default();
    println!(
        "缓存 = {}，库 = {}，是否脏 = {}",
        got,
        db[k],
        got != db[k]
    );

    println!("== 2. 先写库，再删缓存 ==");
    let mut db: HashMap<String, String> = HashMap::new();
    db.insert(k.to_string(), "v1".to_string());
    let mut cache = Cache::default();
    cache.set(k, "v1");
    let old = cache.get(k).unwrap_or_default();   // 读请求先拿到旧值
    db.insert(k.to_string(), "v2".to_string());   // 写请求：先写库
    cache.delete(k);                              // 再删缓存
    cache.set(k, &old);                           // 读请求迟到，把旧值写回
    let got = cache.get(k).unwrap_or_default();
    println!(
        "缓存 = {}，库 = {}，是否脏 = {}",
        got,
        db[k],
        got != db[k]
    );

    println!("== 3. 延迟双删 ==");
    let mut db: HashMap<String, String> = HashMap::new();
    db.insert(k.to_string(), "v1".to_string());
    let mut cache = Cache::default();
    cache.set(k, "v1");
    let old = cache.get(k).unwrap_or_default();
    cache.delete(k);
    db.insert(k.to_string(), "v2".to_string());
    cache.set(k, &old);                    // 迟到的读请求仍回填旧值
    cache.delete(k);                       // 延迟再删一次
    let hit = cache.get(k).is_some();
    let next = read_through(&mut cache, &db, k);
    println!("缓存命中 = {}，下一次读会回源拿到 {}", hit, next);

    println!("== 4. 结论 ==");
    println!("先删缓存有窗口，先写库再删缓存也有窗口，双删只能把窗口压小");
    println!("要强一致就别用缓存读；要最终一致就配 binlog 订阅 + 重试删除");
}
