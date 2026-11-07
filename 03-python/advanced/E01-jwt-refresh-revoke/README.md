# E01 JWT 刷新与吊销

## 问题复现

JWT 无状态，签发后收不回来。登出、改密码后旧 token 仍有效；refresh 不轮换则泄露一次能无限续期。

## 错误写法

```python
if hmac.compare_digest(sig, sign(payload)):
    return payload      # 不验 exp、不查黑名单
```

## 正确写法

```python
if parts[0] != "HS256": return "", "bad format or alg"
if not hmac.compare_digest(parts[2], sign(parts[1])): return "", "signature mismatch"
if now >= int(field(parts[1], "exp")): return "", "token expired"
if field(parts[1], "jti") in revoked: return "", "token revoked"
```

## 怎么验证

```bash
python3 main.py
```

预期：正常 `ok = True`；篡改 `signature mismatch`；旧 refresh `token revoked`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 签名 | `crypto/hmac` | `createHmac` | `hmac` | `Mac` | 需 crate |
| 常量时间比较 | `hmac.Equal` | `timingSafeEqual` | `compare_digest` | `isEqual` | `ConstantTimeEq` |
| 吊销 | jti 黑名单 | 同 | 同 | 同 | 同 |
