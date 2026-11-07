# E01 JWT 刷新与吊销

## 问题复现

JWT 无状态，签发后收不回来。登出、改密码后旧 token 仍有效；refresh 不轮换则泄露一次能无限续期。

## 错误写法

```java
if (sig.equals(sign(payload))) return payload;   // 不验 exp、不查黑名单
```

## 正确写法

```java
if (!parts[0].equals("HS256")) return new String[]{"", "bad format or alg"};
if (!parts[2].equals(sign(parts[1]))) return new String[]{"", "signature mismatch"};
if (now >= exp) return new String[]{"", "token expired"};
if (revoked.contains(jti)) return new String[]{"", "token revoked"};
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：正常 `ok = true`；篡改 `signature mismatch`；旧 refresh `token revoked`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 签名 | `crypto/hmac` | `createHmac` | `hmac` | `javax.crypto.Mac` | 需 crate |
| 常量时间比较 | `hmac.Equal` | `timingSafeEqual` | `compare_digest` | `MessageDigest.isEqual` | `ConstantTimeEq` |
| 吊销 | jti 黑名单 | 同 | 同 | 同 | 同 |
