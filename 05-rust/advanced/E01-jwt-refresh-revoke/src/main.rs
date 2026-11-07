use std::collections::hash_map::DefaultHasher;
use std::collections::HashSet;
use std::hash::{Hash, Hasher};

const SECRET: &str = "s3cr3t";

// Rust 标准库没有 HMAC/SHA-256，这里用 DefaultHasher 演示流程。
// 生产必须用 hmac + sha2 crate，别拿这个当签名。
fn sign(payload: &str) -> String {
    let mut h = DefaultHasher::new();
    format!("HS256.{}.{}", SECRET, payload).hash(&mut h);
    format!("{:016x}{:016x}", h.finish(), h.finish())
}

fn issue(sub: &str, jti: &str, exp: i64) -> String {
    let payload = format!("sub={};jti={};exp={}", sub, jti, exp);
    format!("HS256.{}.{}", payload, sign(&payload))
}

fn field<'a>(payload: &'a str, key: &str) -> &'a str {
    for kv in payload.split(';') {
        if let Some(rest) = kv.strip_prefix(&format!("{}=", key)) {
            return rest;
        }
    }
    ""
}

fn verify(token: &str, now: i64, revoked: &HashSet<String>) -> (String, String) {
    let parts: Vec<&str> = token.split('.').collect();
    if parts.len() != 3 || parts[0] != "HS256" {
        return (String::new(), "bad format or alg".to_string());
    }
    if parts[2] != sign(parts[1]) {
        return (String::new(), "signature mismatch".to_string());
    }
    let exp: i64 = field(parts[1], "exp").parse().unwrap_or(0);
    if now >= exp {
        return (String::new(), "token expired".to_string());
    }
    if revoked.contains(field(parts[1], "jti")) {
        return (String::new(), "token revoked".to_string());
    }
    (field(parts[1], "sub").to_string(), String::new())
}

fn main() {
    let now: i64 = 1000;
    let mut revoked: HashSet<String> = HashSet::new();

    println!("== 1. 签发 ==");
    let access = issue("u1", "jti-a1", now + 60);
    let refresh = issue("u1", "jti-r1", now + 86400);
    println!("access.jti = jti-a1，exp = +60s");
    println!("refresh.jti = jti-r1，exp = +86400s");

    println!("== 2. 校验正常 token ==");
    let (sub, err) = verify(&access, now, &revoked);
    println!("ok = {}，sub = {}", err.is_empty(), sub);

    println!("== 3. 篡改 payload ==");
    let tampered = access.replace("sub=u1", "sub=u2");
    let (_, err) = verify(&tampered, now, &revoked);
    println!("ok = {}，reason = {}", err.is_empty(), err);

    println!("== 4. 刷新：轮换 refresh token ==");
    revoked.insert("jti-r1".to_string());
    let _new_refresh = issue("u1", "jti-r2", now + 86400);
    println!("新 refresh.jti = jti-r2，旧 jti-r1 进黑名单");
    let (_, err) = verify(&refresh, now, &revoked);
    println!("用旧 refresh 再刷：ok = {}，reason = {}", err.is_empty(), err);

    println!("== 5. 吊销 access token ==");
    revoked.insert("jti-a1".to_string());
    let (_, err) = verify(&access, now, &revoked);
    println!("ok = {}，reason = {}", err.is_empty(), err);

    println!("== 6. 结论 ==");
    println!("签名只保证没被篡改，吊销必须靠 jti 黑名单（或短 TTL + 版本号）");
    println!("Rust 无 std 加密库，生产用 hmac + sha2 + base64 crate");
}
