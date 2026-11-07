package main

import "fmt"

type Cache struct {
	data map[string]string
}

func NewCache() *Cache { return &Cache{data: map[string]string{}} }

func (c *Cache) Get(k string) (string, bool) {
	v, ok := c.data[k]
	return v, ok
}

func (c *Cache) Set(k, v string) { c.data[k] = v }
func (c *Cache) Del(k string)    { delete(c.data, k) }

func readThrough(c *Cache, db map[string]string, k string) string {
	if v, ok := c.Get(k); ok {
		return v
	}
	v := db[k]          // 缓存未命中，回源
	c.Set(k, v)
	return v
}

func main() {
	const key = "user:1"

	fmt.Println("== 1. 先删缓存，再写库 ==")
	db := map[string]string{key: "v1"}
	c := NewCache()
	c.Set(key, "v1")
	c.Del(key)                       // 写请求：先删缓存
	_ = readThrough(c, db, key)      // 读请求：缓存没了，回源读到旧值 v1
	db[key] = "v2"                   // 写请求：再写库
	got, _ := c.Get(key)
	fmt.Printf("缓存 = %s，库 = %s，是否脏 = %v\n", got, db[key], got != db[key])

	fmt.Println("== 2. 先写库，再删缓存 ==")
	db = map[string]string{key: "v1"}
	c = NewCache()
	c.Set(key, "v1")
	old, _ := c.Get(key)             // 读请求先拿到旧值
	db[key] = "v2"                   // 写请求：先写库
	c.Del(key)                       // 再删缓存
	c.Set(key, old)                  // 读请求迟到，把旧值写回缓存
	got, _ = c.Get(key)
	fmt.Printf("缓存 = %s，库 = %s，是否脏 = %v\n", got, db[key], got != db[key])

	fmt.Println("== 3. 延迟双删 ==")
	db = map[string]string{key: "v1"}
	c = NewCache()
	c.Set(key, "v1")
	old, _ = c.Get(key)
	c.Del(key)
	db[key] = "v2"
	c.Set(key, old)                  // 迟到的读请求仍然回填旧值
	c.Del(key)                       // 延迟一段时间后再删一次，兜住窗口
	got, ok := c.Get(key)
	fmt.Printf("缓存命中 = %v，下一次读会回源拿到 %s\n", ok, readThrough(c, db, key))

	fmt.Println("== 4. 结论 ==")
	fmt.Println("先删缓存有窗口，先写库再删缓存也有窗口，双删只能把窗口压小")
	fmt.Println("要强一致就别用缓存读；要最终一致就配 binlog 订阅 + 重试删除")
}
