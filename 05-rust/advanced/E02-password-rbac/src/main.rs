use std::collections::hash_map::DefaultHasher;
use std::hash::{Hash, Hasher};

// 标准库没有 PBKDF2/bcrypt/argon2，这里用「加盐 + 多次哈希」演示流程。
// 生产必须用 argon2 或 bcrypt crate，并且成本参数按机器调到 50-100ms。
fn hash_password(password: &str, salt: &str) -> String {
    let mut acc = format!("{}{}", salt, password);
    for _ in 0..1000 {
        let mut h = DefaultHasher::new();
        acc.hash(&mut h);
        acc = format!("{:016x}", h.finish());
    }
    acc
}

// 固定时间比较：不因为前缀相同就提前返回
fn safe_equal(a: &str, b: &str) -> bool {
    if a.len() != b.len() {
        return false;
    }
    let mut diff = 0u8;
    for (x, y) in a.bytes().zip(b.bytes()) {
        diff |= x ^ y;
    }
    diff == 0
}

fn role_perms(role: &str) -> Vec<&'static str> {
    match role {
        "admin" => vec!["order:read", "order:write", "user:delete"],
        "viewer" => vec!["order:read"],
        _ => vec![],
    }
}

fn allowed(role: &str, perm: &str) -> bool {
    role_perms(role).contains(&perm)
}

fn main() {
    let salt = "salt-8f2a";
    let hashed = hash_password("correct-horse", salt);

    println!("== 1. 加盐哈希 ==");
    println!("salt = {}，hash = {}...", salt, &hashed[..16]);
    println!(
        "同一密码 + 同一盐 结果一致 = {}",
        hashed == hash_password("correct-horse", salt)
    );
    println!(
        "同一密码 + 不同盐 结果不同 = {}",
        hashed != hash_password("correct-horse", "salt-other")
    );

    println!("== 2. 校验 ==");
    println!(
        "正确密码 = {}",
        safe_equal(&hashed, &hash_password("correct-horse", salt))
    );
    println!("错误密码 = {}", safe_equal(&hashed, &hash_password("wrong", salt)));

    println!("== 3. RBAC ==");
    println!("admin 能 order:write = {}", allowed("admin", "order:write"));
    println!("viewer 能 order:write = {}", allowed("viewer", "order:write"));
    println!("viewer 能 user:delete = {}", allowed("viewer", "user:delete"));

    println!("== 4. 结论 ==");
    println!("密码永远不存明文，也不要可逆加密；盐必须每用户独立随机");
    println!("Rust 无 std KDF，生产用 argon2 crate；手写 KDF 只在教学场景可接受");
}
