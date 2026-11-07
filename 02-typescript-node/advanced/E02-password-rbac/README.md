# E02 密码哈希与 RBAC

## 问题复现

存明文 → 拖库全泄露；MD5 无盐 → 彩虹表；全站共用一个盐 → 撞库一次全中。

## 错误写法

```js
db.save(md5(password));          // 无盐快速哈希
if (hash === stored) { ... }     // 逐字节比较，时序泄露
```

## 正确写法

```js
const salt = randomBytes(16).toString('hex');
const hash = pbkdf2Sync(password, salt, 1000, 16, 'sha256').toString('hex');
timingSafeEqual(Buffer.from(hash, 'hex'), Buffer.from(stored, 'hex'));
```

生产用 `bcrypt` / `argon2`，成本调到单次 50-100ms。

## 怎么验证

```bash
node index.mjs
```

预期：同盐一致、异盐不同、正确密码 `true`、错误密码 `false`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| KDF | `x/crypto/bcrypt` | `crypto.pbkdf2` | `hashlib.pbkdf2_hmac` | `PBKDF2WithHmacSHA256` | `argon2` crate |
| 固定时间比较 | `hmac.Equal` | `timingSafeEqual` | `hmac.compare_digest` | `MessageDigest.isEqual` | `ConstantTimeEq` |
| 常见坑 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 |
