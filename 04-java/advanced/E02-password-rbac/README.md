# E02 密码哈希与 RBAC

## 问题复现

存明文 → 拖库全泄露；MD5 无盐 → 彩虹表；全站一个盐 → 撞库全中。

## 错误写法

```java
db.save(md5(password));
if (hash.equals(stored)) { ... }      // 逐字节比较
```

## 正确写法

```java
PBEKeySpec spec = new PBEKeySpec(password.toCharArray(), salt.getBytes(), 1000, 128);
SecretKeyFactory f = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256");
MessageDigest.isEqual(a.getBytes(), b.getBytes());
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：同盐一致、异盐不同、正确密码 `true`、错误密码 `false`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| KDF | `x/crypto/bcrypt` | `crypto.pbkdf2` | `hashlib.pbkdf2_hmac` | `PBKDF2WithHmacSHA256` | `argon2` crate |
| 固定时间比较 | `hmac.Equal` | `timingSafeEqual` | `compare_digest` | `MessageDigest.isEqual` | `ConstantTimeEq` |
| 常见坑 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 |
