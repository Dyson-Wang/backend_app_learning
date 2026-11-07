# E01 JWT 刷新与吊销

## 问题复现

JWT 无状态，签发后收不回来。登出、改密码后旧 token 仍有效；refresh 不轮换则泄露一次能无限续期。

## 错误写法

```rust
if parts[2] == sign(parts[1]) { return Ok(payload); }   // 不验 exp、不查黑名单
```

## 正确写法

```rust
if parts[0] != "HS256" { return (String::new(), "bad format or alg".into()); }
if parts[2] != sign(parts[1]) { return (String::new(), "signature mismatch".into()); }
if now >= exp { return (String::new(), "token expired".into()); }
if revoked.contains(jti) { return (String::new(), "token revoked".into()); }
```

> 注意：本示例的 `sign()` 用 `DefaultHasher` 代替 HMAC-SHA256，因为 Rust 标准库不含加密算法。生产必须用 `hmac` + `sha2` + `base64` crate。

## 怎么验证

```bash
cargo run --quiet
```

预期：正常 `ok = true`；篡改 `signature mismatch`；旧 refresh `token revoked`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 签名 | std `crypto/hmac` | std `crypto` | std `hmac` | std `Mac` | **需第三方 crate** |
| 常量时间比较 | `hmac.Equal` | `timingSafeEqual` | `compare_digest` | `MessageDigest.isEqual` | `ConstantTimeEq` |
| 吊销 | jti 黑名单 | 同 | 同 | 同 | 同 |
