package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
)

func main() {
	// 1. 时间点 / 时区 / 格式化 / 解析：用固定时间，不依赖当前时间
	layout := "2006-01-02 15:04:05"
	t := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	fmt.Printf("时间点(UTC): %s, Unix = %d\n", t.Format(layout), t.Unix())
	cst := time.FixedZone("CST", 8*3600)
	fmt.Printf("时区转换: 同一时刻在东八区 = %s\n", t.In(cst).Format(layout))
	parsed, err := time.Parse(layout, "2024-06-01 12:30:00")
	fmt.Printf("解析: err = %v, 解析结果 = %s\n", err, parsed.Format(layout))

	// 2. 时长：两个固定时间点相减，演示 Duration
	t2 := time.Date(2024, 1, 3, 6, 4, 5, 0, time.UTC)
	d := t2.Sub(t)
	fmt.Printf("时长: %s = %.2f 小时 = %.0f 分钟\n", d, d.Hours(), d.Minutes())
	fmt.Printf("时长构造: 90*time.Minute = %s, 舍入到小时 = %s\n",
		(90 * time.Minute).String(), (90 * time.Minute).Round(time.Hour).String())

	// 3. 正则：匹配、捕获组、命名捕获、替换、查找全部
	re := regexp.MustCompile(`(?P<user>[a-z]+)@(?P<host>[a-z]+)\.(?P<tld>[a-z]+)`)
	m := re.FindStringSubmatch("联系 a@b.com 或 x@y.org")
	fmt.Printf("正则: 是否匹配 = %t\n", re.MatchString("a@b.com"))
	if m != nil {
		fmt.Printf("正则: 完整匹配 = %s, 捕获组 = %v\n", m[0], m[1:])
		fmt.Printf("正则: 命名组 user=%s host=%s tld=%s\n",
			m[re.SubexpIndex("user")], m[re.SubexpIndex("host")], m[re.SubexpIndex("tld")])
	}
	fmt.Printf("正则替换: %s\n", re.ReplaceAllString("mail a@b.com", "[隐藏邮箱]"))
	fmt.Printf("正则查找全部: %v\n", re.FindAllString("a@b.com, x@y.org", -1))

	// 4. 随机数：固定种子保证输出可复现
	r := rand.New(rand.NewSource(7))
	fmt.Printf("随机数(种子7): Intn(100) = %d, Perm(5) = %v\n", r.Intn(100), r.Perm(5))
	nums := []int{5, 3, 9, 1, 7}
	r.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	fmt.Printf("随机洗牌(可复现): %v\n", nums)

	// 5. 排序：Ints / Strings / Slice 自定义比较
	ints := []int{4, 1, 3, 2}
	sort.Ints(ints)
	fmt.Printf("排序: sort.Ints = %v, sort.SearchInts(3) = %d\n", ints, sort.SearchInts(ints, 3))
	words := []string{"go", "rust", "java", "c"}
	sort.Strings(words)
	fmt.Printf("排序: sort.Strings = %v\n", words)
	people := []string{"Ann", "Bob", "Zoe", "Cid"}
	sort.Slice(people, func(i, j int) bool {
		return strings.ToLower(people[i]) < strings.ToLower(people[j])
	})
	fmt.Printf("排序: sort.Slice(忽略大小写) = %v\n", people)
	fmt.Printf("工具函数: strings.Repeat/Contains/ReplaceAll = %s %t %s\n",
		strings.Repeat("ab", 3), strings.Contains("backend", "end"),
		strings.ReplaceAll("a-b-c", "-", "/"))
}
