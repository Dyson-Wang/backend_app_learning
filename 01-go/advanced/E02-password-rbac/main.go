package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// PBKDF2-HMAC-SHA256，标准库没有，手写 12 行
func pbkdf2(password, salt []byte, iter, keyLen int) []byte {
	var out []byte
	block := 1
	for len(out) < keyLen {
		m := hmac.New(sha256.New, password)
		m.Write(salt)
		m.Write([]byte{byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block)})
		u := m.Sum(nil)
		t := append([]byte(nil), u...)
		for i := 1; i < iter; i++ {
			m.Reset()
			m.Write(u)
			u = m.Sum(nil)
			for j := range t {
				t[j] ^= u[j]
			}
		}
		out = append(out, t...)
		block++
	}
	return out[:keyLen]
}

func hashPassword(password, salt string) string {
	return hex.EncodeToString(pbkdf2([]byte(password), []byte(salt), 1000, 16))
}

// 固定时间比较，避免按字节提前返回泄露信息
func safeEqual(a, b string) bool {
	return hmac.Equal([]byte(a), []byte(b))
}

var rolePerms = map[string][]string{
	"admin":  {"order:read", "order:write", "user:delete"},
	"viewer": {"order:read"},
}

func allowed(role, perm string) bool {
	for _, p := range rolePerms[role] {
		if p == perm {
			return true
		}
	}
	return false
}

func main() {
	salt := "salt-8f2a"
	hash := hashPassword("correct-horse", salt)
	fmt.Println("== 1. 加盐哈希 ==")
	fmt.Printf("salt = %s，hash = %s...\n", salt, hash[:16])
	fmt.Printf("同一密码 + 同一盐 结果一致 = %v\n", hash == hashPassword("correct-horse", salt))
	fmt.Printf("同一密码 + 不同盐 结果不同 = %v\n", hash != hashPassword("correct-horse", "salt-other"))

	fmt.Println("== 2. 校验 ==")
	fmt.Printf("正确密码 = %v\n", safeEqual(hash, hashPassword("correct-horse", salt)))
	fmt.Printf("错误密码 = %v\n", safeEqual(hash, hashPassword("wrong", salt)))

	fmt.Println("== 3. RBAC ==")
	fmt.Printf("admin 能 order:write = %v\n", allowed("admin", "order:write"))
	fmt.Printf("viewer 能 order:write = %v\n", allowed("viewer", "order:write"))
	fmt.Printf("viewer 能 user:delete = %v\n", allowed("viewer", "user:delete"))

	fmt.Println("== 4. 结论 ==")
	fmt.Println("密码永远不存明文，也不要可逆加密；盐必须每用户独立随机")
	fmt.Println("1000 轮只是演示，生产用 bcrypt/argon2id，成本参数按机器调")
}
