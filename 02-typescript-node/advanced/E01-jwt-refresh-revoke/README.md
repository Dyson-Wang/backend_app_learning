# E01 JWT 刷新与吊销

## 问题复现

JWT 无状态，签发后收不回来。登出、改密码、被盗号，旧 token 仍然有效到过期。refresh 不轮换时，泄露一次能无限续期。

## 错误写法

```js
if (sig === sign(payload)) return payload;   // 不验 exp、不查黑名单
```

## 正确写法

```js
if (parts[0] !== 'HS256') return ['', 'bad format or alg'];   // 防 alg=none
if (parts[2] !== sign(parts[1])) return ['', 'signature mismatch'];
if (now >= exp) return ['', 'token expired'];
if (revoked.has(jti)) return ['', 'token revoked'];
```

## 怎么验证

```bash
node index.mjs
```

预期：正常 `ok = true`；篡改 `signature mismatch`；旧 refresh `token revoked`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 签名 | `crypto/hmac` | `crypto.createHmac` | `hmac` | `Mac` | 需 crate |
| 黑名单 | Redis SET + TTL | 同 | 同 | 同 | 同 |
| 吊销替代 | 短 TTL + 版本号 | 同 | 同 | 同 | 同 |
