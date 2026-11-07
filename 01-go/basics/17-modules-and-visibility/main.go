package main

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	mrand "math/rand"
	"net/url"
	"reflect"
)

// Config 演示导出/私有字段：大写导出，小写私有
type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	secret  string // 未导出
	version string // 未导出
}

// Exported 首字母大写 → 导出；buildTag 小写 → 包内可见
func Exported() string { return "可被其他包访问" }
func buildTag() string  { return "仅本包可见" }

func main() {
	// 1. 导出/私有规则：encoding/json 只序列化导出字段
	cfg := Config{Host: "localhost", Port: 8080, secret: "s3cr3t", version: "v1"}
	data, _ := json.Marshal(cfg)
	fmt.Printf("导出规则: json.Marshal = %s（secret/version 未导出，被忽略）\n", data)
	fmt.Printf("包内访问: Exported() = %s, buildTag() = %s（同包内私有也可访问）\n",
		Exported(), buildTag())

	// 2. 命名冲突与别名：math/rand 与 crypto/rand 都叫 rand，用别名区分
	r := mrand.New(mrand.NewSource(42)) // 固定种子 → 输出可复现
	fmt.Printf("别名 mrand(种子42): %d, %d\n", r.Intn(100), r.Intn(100))
	buf := make([]byte, 8)
	n, err := crand.Read(buf)
	fmt.Printf("别名 crand: 读取随机字节数 = %d, err = %v（只演示 API，不打印随机值）\n", n, err)

	// 3. 目录即包：包由导入路径标识，包名是路径最后一段
	t := reflect.TypeOf(url.URL{})
	fmt.Printf("包标识: url.URL 的 PkgPath = %q, Name = %q（包名 = 路径末段 url）\n",
		t.PkgPath(), t.Name())
	fmt.Printf("本包类型: Config 的 Name = %q（主包路径由 go.mod 决定）\n",
		reflect.TypeOf(Config{}).Name())

	// 4. 反射看字段导出性：StructField.PkgPath 非空即未导出
	ct := reflect.TypeOf(Config{})
	for i := 0; i < ct.NumField(); i++ {
		f := ct.Field(i)
		if f.PkgPath != "" {
			fmt.Printf("字段可见性: %s 未导出 → 包外无法访问/序列化\n", f.Name)
		}
	}
}
