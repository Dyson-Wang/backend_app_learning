package main

import (
	"fmt"
	"sync"
	"time"
)

// 用内存 map 模拟 Redis
type FakeRedis struct {
	mu   sync.Mutex
	data map[string]string
	exp  map[string]time.Time
	now  func() time.Time
}

func NewFakeRedis() *FakeRedis {
	return &FakeRedis{data: map[string]string{}, exp: map[string]time.Time{}, now: time.Now}
}

func (r *FakeRedis) alive(key string) bool {
	t, ok := r.exp[key]
	return ok && r.now().Before(t)
}

func (r *FakeRedis) SetNX(key, val string, ttl time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.alive(key) {
		return false
	}
	r.data[key] = val
	r.exp[key] = r.now().Add(ttl)
	return true
}

// 错误解锁：不校验 token，会删掉别人的锁
func (r *FakeRedis) Del(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, key)
	delete(r.exp, key)
}

// 正确解锁：先比较 token 再删（等价于 Lua 脚本）
func (r *FakeRedis) DelIfMatch(key, token string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.alive(key) || r.data[key] != token {
		return false
	}
	delete(r.data, key)
	delete(r.exp, key)
	return true
}

func (r *FakeRedis) Get(key string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.alive(key) {
		return ""
	}
	return r.data[key]
}

func main() {
	clock := time.Unix(0, 0)
	redis := NewFakeRedis()
	redis.now = func() time.Time { return clock }
	const key = "lock:order:1"

	fmt.Println("== 1. A 抢到锁 ==")
	fmt.Printf("A SetNX = %v, holder=%s\n", redis.SetNX(key, "token-A", 100*time.Millisecond), redis.Get(key))

	fmt.Println("== 2. B 抢锁失败 ==")
	fmt.Printf("B SetNX = %v, holder=%s\n", redis.SetNX(key, "token-B", 100*time.Millisecond), redis.Get(key))

	fmt.Println("== 3. A 的业务超时，锁过期 ==")
	clock = clock.Add(150 * time.Millisecond)
	fmt.Printf("锁是否还有效 = %v\n", redis.Get(key) != "")

	fmt.Println("== 4. B 抢到锁 ==")
	fmt.Printf("B SetNX = %v, holder=%s\n", redis.SetNX(key, "token-B", 100*time.Millisecond), redis.Get(key))

	fmt.Println("== 5. A 醒来直接 DEL：删掉了 B 的锁（错误） ==")
	redis.Del(key)
	fmt.Printf("holder=%q（B 的锁被误删）\n", redis.Get(key))

	fmt.Println("== 6. 修正：DelIfMatch 校验 token ==")
	redis.SetNX(key, "token-B", 100*time.Millisecond)
	fmt.Printf("A DelIfMatch(token-A) = %v, holder=%s\n", redis.DelIfMatch(key, "token-A"), redis.Get(key))
	fmt.Printf("B DelIfMatch(token-B) = %v, holder=%q\n", redis.DelIfMatch(key, "token-B"), redis.Get(key))

	fmt.Println("== 7. 结论 ==")
	fmt.Println("加锁：SET key token NX PX ttl；解锁：比较 token 再 DEL，必须原子")
	fmt.Println("业务可能超过 ttl 时，还要有看门狗续期，否则锁提前释放")
}
