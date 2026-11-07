package main

import "fmt"

type Source struct {
	Value  string
	Origin string
}

// 优先级：环境变量 > 配置文件 > 默认值
func resolve(key string, def string, file map[string]string, env map[string]string) Source {
	if v, ok := env[key]; ok {
		return Source{v, "env"}
	}
	if v, ok := file[key]; ok {
		return Source{v, "file"}
	}
	return Source{def, "default"}
}

type Migration struct {
	Version int
	SQL     string
}

func migrate(applied map[int]bool, all []Migration) []int {
	var ran []int
	for _, m := range all {
		if applied[m.Version] {
			continue
		}
		applied[m.Version] = true
		ran = append(ran, m.Version)
	}
	return ran
}

func main() {
	file := map[string]string{"db.name": "brm", "port": "9090"}
	env := map[string]string{"port": "8080"}

	fmt.Println("== 1. 配置优先级 ==")
	for _, k := range []string{"port", "db.host", "db.name"} {
		s := resolve(k, "localhost", file, env)
		fmt.Printf("%s = %s（来源：%s）\n", k, s.Value, s.Origin)
	}

	fmt.Println("== 2. 密钥缺失拒绝启动 ==")
	secrets := map[string]string{}
	if _, ok := secrets["DATABASE_URL"]; !ok {
		fmt.Println("启动失败：DATABASE_URL 未设置（不要给密钥默认值）")
	}

	fmt.Println("== 3. migration ==")
	all := []Migration{
		{1, "create accounts"},
		{2, "add email"},
		{3, "add index"},
	}
	applied := map[int]bool{}
	for _, v := range migrate(applied, all) {
		fmt.Printf("应用 v%d %s\n", v, all[v-1].SQL)
	}
	fmt.Printf("当前版本 = %d\n", len(applied))
	again := migrate(applied, all)
	fmt.Printf("再次执行：%d 个待应用，当前版本 = %d\n", len(again), len(applied))

	fmt.Println("== 4. 结论 ==")
	fmt.Println("配置分层要能说出「这个值从哪来」，否则排查时全靠猜")
	fmt.Println("migration 必须可重复执行：靠版本表跳过，不靠人记得跑过没有")
}
