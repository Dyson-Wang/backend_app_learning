# E02 密码哈希与 RBAC

## 问题复现

存明文 → 拖库全泄露；MD5 无盐 → 彩虹表；全站一个盐 → 撞库全中。

## 错误写法

```rust
let h = format!("{:x}", md5(password));     // 无盐
if h == stored { ... }                      // 逐字节比较
```

## 正确写法

```rust
// 每用户独立盐
let salt = random_hex(16);
let h = hash_password(&password, &salt);
// 固定时间比较
let ok = safe_equal(&h, &stored);
```

> 本示例的 `hash_password` 用 `DefaultHasher` 迭代 1000 次代替 PBKDF2，因为标准库不含 KDF。生产用 `argon2` crate。

## 怎么验证

```bash
cargo run --quiet
```

预期：同盐一致、异盐不同、正确密码 `true`、错误密码 `false`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| KDF | `x/crypto/bcrypt` | `crypto.pbkdf2` | `hashlib.pbkdf2_hmac` | `PBKDF2WithHmacSHA256` | `argon2` crate |
| 固定时间比较 | `hmac.Equal` | `timingSafeEqual` | `compare_digest` | `MessageDigest.isEqual` | 手写 `xor` 累积 |
| 生态特点 | 半标准 | 内建齐全 | 内建齐全 | 内建齐全 | **一切靠 crate** |
