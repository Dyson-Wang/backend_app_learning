// Package toolkit 演示 Go 的包组织与 internal 可见性规则。
// internal/ 下的包只能被其父目录（这里是 01-go）内的代码 import，外部模块无法引用。
package toolkit

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

func NewRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
