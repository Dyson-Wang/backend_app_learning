use std::collections::{HashMap, HashSet};

fn resolve<'a>(
    key: &str,
    def: &'a str,
    file: &'a HashMap<String, String>,
    env: &'a HashMap<String, String>,
) -> (&'a str, &'static str) {
    if let Some(v) = env.get(key) {
        return (v, "env");
    }
    if let Some(v) = file.get(key) {
        return (v, "file");
    }
    (def, "default")
}

fn migrate(applied: &mut HashSet<u32>, all: &[(u32, &str)]) -> Vec<u32> {
    let mut ran = Vec::new();
    for (version, _) in all {
        if applied.contains(version) {
            continue;
        }
        applied.insert(*version);
        ran.push(*version);
    }
    ran
}

fn main() {
    let mut file: HashMap<String, String> = HashMap::new();
    file.insert("db.name".to_string(), "brm".to_string());
    file.insert("port".to_string(), "9090".to_string());
    let mut env: HashMap<String, String> = HashMap::new();
    env.insert("port".to_string(), "8080".to_string());

    println!("== 1. 配置优先级 ==");
    for k in ["port", "db.host", "db.name"] {
        let (value, origin) = resolve(k, "localhost", &file, &env);
        println!("{} = {}（来源：{}）", k, value, origin);
    }

    println!("== 2. 密钥缺失拒绝启动 ==");
    let secrets: HashMap<String, String> = HashMap::new();
    if !secrets.contains_key("DATABASE_URL") {
        println!("启动失败：DATABASE_URL 未设置（不要给密钥默认值）");
    }

    println!("== 3. migration ==");
    let all = [
        (1u32, "create accounts"),
        (2u32, "add email"),
        (3u32, "add index"),
    ];
    let mut applied: HashSet<u32> = HashSet::new();
    for v in migrate(&mut applied, &all) {
        println!("应用 v{} {}", v, all[(v - 1) as usize].1);
    }
    println!("当前版本 = {}", applied.len());
    let again = migrate(&mut applied, &all);
    println!("再次执行：{} 个待应用，当前版本 = {}", again.len(), applied.len());

    println!("== 4. 结论 ==");
    println!("配置分层要能说出「这个值从哪来」，否则排查时全靠猜");
    println!("migration 必须可重复执行：靠版本表跳过，不靠人记得跑过没有");
}
