# E02 密码哈希与 RBAC

## 问题复现

存明文 → 拖库全泄露；MD5 无盐 → 彩虹表；全站一个盐 → 撞库全中。

## 错误写法

```python
db.save(hashlib.md5(password.encode()).hexdigest())
if h == stored: ...      # 逐字节比较
```

## 正确写法

```python
salt = secrets.token_hex(16)
h = hashlib.pbkdf2_hmac("sha256", password.encode(), salt.encode(), 1000, 16)
hmac.compare_digest(h, stored)
```

## 怎么验证

```bash
python3 main.py
```

预期：同盐一致、异盐不同、正确密码 `True`、错误密码 `False`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| KDF | `x/crypto/bcrypt` | `crypto.pbkdf2` | `hashlib.pbkdf2_hmac` | `PBKDF2WithHmacSHA256` | `argon2` crate |
| 固定时间比较 | `hmac.Equal` | `timingSafeEqual` | `compare_digest` | `isEqual` | `ConstantTimeEq` |
| 常见坑 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 | 忘记加盐 |
