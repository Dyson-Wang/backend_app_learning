package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const secret = "s3cr3t"

func sign(payload string) string {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte("HS256." + payload))
	return hex.EncodeToString(m.Sum(nil))[:32]
}

func issue(sub, jti string, exp int64) string {
	payload := fmt.Sprintf("sub=%s;jti=%s;exp=%d", sub, jti, exp)
	return "HS256." + payload + "." + sign(payload)
}

func field(payload, key string) string {
	for _, kv := range strings.Split(payload, ";") {
		if strings.HasPrefix(kv, key+"=") {
			return kv[len(key)+1:]
		}
	}
	return ""
}

func verify(token string, now int64, revoked map[string]bool) (string, string) {
	parts := strings.SplitN(token, ".", 3)
	if len(parts) != 3 || parts[0] != "HS256" {
		return "", "bad format or alg"
	}
	if !hmac.Equal([]byte(parts[2]), []byte(sign(parts[1]))) {
		return "", "signature mismatch"
	}
	exp, _ := strconv.ParseInt(field(parts[1], "exp"), 10, 64)
	if now >= exp {
		return "", "token expired"
	}
	if revoked[field(parts[1], "jti")] {
		return "", "token revoked"
	}
	return field(parts[1], "sub"), ""
}

func main() {
	now := int64(1000)
	revoked := map[string]bool{}

	fmt.Println("== 1. 签发 ==")
	access := issue("u1", "jti-a1", now+60)
	refresh := issue("u1", "jti-r1", now+86400)
	fmt.Println("access.jti = jti-a1，exp = +60s")
	fmt.Println("refresh.jti = jti-r1，exp = +86400s")

	fmt.Println("== 2. 校验正常 token ==")
	sub, err := verify(access, now, revoked)
	fmt.Printf("ok = %v，sub = %s\n", err == "", sub)

	fmt.Println("== 3. 篡改 payload ==")
	tampered := strings.Replace(access, "sub=u1", "sub=u2", 1)
	_, err = verify(tampered, now, revoked)
	fmt.Printf("ok = %v，reason = %s\n", err == "", err)

	fmt.Println("== 4. 刷新：轮换 refresh token ==")
	revoked["jti-r1"] = true
	newRefresh := issue("u1", "jti-r2", now+86400)
	fmt.Printf("新 refresh.jti = jti-r2，旧 jti-r1 进黑名单\n")
	_, err = verify(refresh, now, revoked)
	fmt.Printf("用旧 refresh 再刷：ok = %v，reason = %s\n", err == "", err)
	_ = newRefresh

	fmt.Println("== 5. 吊销 access token ==")
	revoked["jti-a1"] = true
	_, err = verify(access, now, revoked)
	fmt.Printf("ok = %v，reason = %s\n", err == "", err)

	fmt.Println("== 6. 结论 ==")
	fmt.Println("签名只保证没被篡改，吊销必须靠 jti 黑名单（或短 TTL + 版本号）")
	fmt.Println("refresh 必须轮换：旧的用过一次就作废，否则泄露后可长期续命")
}
