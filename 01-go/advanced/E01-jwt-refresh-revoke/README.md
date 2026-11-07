# E01 JWT 刷新与吊销

## 问题复现

JWT 是无状态的，签发后服务端没法「收回」。用户登出、改密码、被盗号，旧 token 仍然有效到过期。另一个坑：refresh token 不轮换，泄露一次就能无限续期。

## 错误写法

```go
// 只验签名，不验 exp、不查黑名单
if hmac.Equal(sig, sign(payload)) { return payload, nil }
```

## 正确写法

三层校验，缺一不可：

```go
if parts[0] != "HS256" { return "", "bad format or alg" }   // 防 alg=none
if !hmac.Equal(...) { return "", "signature mismatch" }
if now >= exp { return "", "token expired" }
if revoked[jti] { return "", "token revoked" }              // 黑名单
```

refresh 轮换：每次刷新生成新 jti，旧的立刻进黑名单。

## 怎么验证

```bash
go run .
```

预期：正常 `ok = true`；篡改 `signature mismatch`；旧 refresh `token revoked`；吊销后 access 也 `token revoked`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 签名 | `crypto/hmac` | `crypto.createHmac` | `hmac` | `javax.crypto.Mac` | 无 std，需 crate |
| 黑名单 | Redis SET + TTL | 同 | 同 | 同 | 同 |
| 吊销替代 | 短 TTL + token 版本号 | 同 | 同 | 同 | 同 |
