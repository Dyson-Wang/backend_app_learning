# E02 密码哈希与 RBAC

## 问题复现

数据库存明文密码，一次拖库全泄露；或者用 MD5 无盐，彩虹表一查就出来；再或者全站共用一个盐，撞库一次全中。

## 错误写法

```go
db.Save(user, md5(password))          // 无盐、快速哈希
if hash == stored { ... }             // 逐字节比较，时序泄露
```

## 正确写法

```go
salt := randomHex(16)                 // 每用户独立
hash := pbkdf2(password, salt, 1000, 16)
// 校验用固定时间比较
hmac.Equal([]byte(hash), []byte(stored))
```

生产直接用 `bcrypt.GenerateFromPassword`（`golang.org/x/crypto`）或 argon2id，成本参数调到单次 50-100ms。

RBAC 用「角色 → 权限集合」查表，别在 handler 里散落 `if role == "admin"`。

## 怎么验证

```bash
go run .
```

预期：同盐一致、异盐不同、正确密码 `true`、错误密码 `false`；admin 有写权限、viewer 没有。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| KDF | `x/crypto/bcrypt` | `crypto.pbkdf2` / `bcrypt` | `hashlib.pbkdf2_hmac` | `PBKDF2WithHmacSHA256` | `argon2` crate |
| 固定时间比较 | `hmac.Equal` | `timingSafeEqual` | `hmac.compare_digest` | `MessageDigest.isEqual` | `subtle::ConstantTimeEq` |
| 常见坑 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 |
