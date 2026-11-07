package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	// 1. 不可变性：字符串不能改 s[i]；转成 []byte 改副本，原串不变
	ascii := "hello"
	bs := []byte(ascii)
	bs[0] = 'H'
	fmt.Printf("不可变: 原串 = %q, []byte 改后转回 = %q\n", ascii, string(bs))

	// 2. 字节 vs 码点：中文一个字通常 3 字节
	s := "Go语言"
	fmt.Printf("长度: %q 字节数 = %d, 码点数 = %d, 合法 UTF-8 = %t\n",
		s, len(s), utf8.RuneCountInString(s), utf8.ValidString(s))

	// 3. 遍历差异：按下标取到字节，range 取到 rune
	fmt.Print("按字节遍历: ")
	for i := 0; i < len(s); i++ {
		fmt.Printf("%02x ", s[i])
	}
	fmt.Println()
	fmt.Print("按 rune 遍历(字节偏移:字符): ")
	for idx, r := range s {
		fmt.Printf("%d:%c ", idx, r)
	}
	fmt.Println()

	// 4. 编码解码：rune <-> UTF-8 字节，以及 base64 / hex 这类字节编码
	r := '语'
	buf := make([]byte, utf8.RuneLen(r))
	n := utf8.EncodeRune(buf, r)
	decoded, size := utf8.DecodeRune(buf)
	fmt.Printf("编解码: rune %q = %d 字节 %v, 解码回 %q(占 %d 字节)\n", r, n, buf, decoded, size)
	enc := base64.StdEncoding.EncodeToString([]byte(s))
	plain, _ := base64.StdEncoding.DecodeString(enc)
	fmt.Printf("base64: %s -> %q\n", enc, plain)
	fmt.Printf("hex: %s\n", hex.EncodeToString([]byte("Go")))

	// 5. 拼接: + / fmt.Sprintf / strings.Builder
	var b strings.Builder
	b.WriteString("Go")
	b.WriteString("语言")
	fmt.Printf("拼接: + = %q, Sprintf = %q, Builder = %q\n",
		"Go"+"语言", fmt.Sprintf("%s%s", "Go", "语言"), b.String())

	// 6. 转义字符串 vs 原始字符串(反引号，不解释转义)
	fmt.Printf("转义: %q\n", "第一行\n第二行\t带制表")
	fmt.Printf("原始: %q\n", `第一行\n第二行`)

	// 7. 大小写 / 查找 / 切分 / 连接
	fmt.Printf("大小写: ToUpper = %s, ToLower = %s\n", strings.ToUpper("Go"), strings.ToLower("GO"))
	fmt.Printf("查找: Contains = %t, Index = %d, HasPrefix = %t\n",
		strings.Contains("backend", "end"), strings.Index("backend", "end"), strings.HasPrefix("backend", "back"))
	parts := strings.Split("a,b,c", ",")
	fmt.Printf("切分与连接: Split = %v, Join = %q, TrimSpace = %q\n",
		parts, strings.Join(parts, "-"), strings.TrimSpace("  hi  "))
}
