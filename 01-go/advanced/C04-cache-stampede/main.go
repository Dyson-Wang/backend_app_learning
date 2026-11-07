package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 模拟数据库：一次查询 20ms，并统计回源次数
type DB struct {
	hits int64
}

func (d *DB) Query(key string) string {
	atomic.AddInt64(&d.hits, 1)
	time.Sleep(20 * time.Millisecond)
	return "value-of-" + key
}

func (d *DB) Hits() int64 { return atomic.LoadInt64(&d.hits) }

// 错误写法：缓存未命中就回源，并发时全部穿透
func unprotectedGet(db *DB, cache map[string]string, mu *sync.Mutex, key string) string {
	mu.Lock()
	v, ok := cache[key]
	mu.Unlock()
	if ok {
		return v
	}
	v = db.Query(key) // 100 个并发请求会全部走到这里
	mu.Lock()
	cache[key] = v
	mu.Unlock()
	return v
}

// 正确写法：互斥 + 双重检查，只放一个请求回源
type GuardedCache struct {
	mu    sync.Mutex
	data  map[string]string
	keyMu map[string]*sync.Mutex
}

func NewGuardedCache() *GuardedCache {
	return &GuardedCache{data: map[string]string{}, keyMu: map[string]*sync.Mutex{}}
}

func (c *GuardedCache) Get(db *DB, key string) string {
	c.mu.Lock()
	if v, ok := c.data[key]; ok {
		c.mu.Unlock()
		return v
	}
	km, ok := c.keyMu[key]
	if !ok {
		km = &sync.Mutex{}
		c.keyMu[key] = km
	}
	c.mu.Unlock()

	km.Lock()
	defer km.Unlock()

	// 双重检查：等锁期间可能已经有人回源完成
	c.mu.Lock()
	if v, ok := c.data[key]; ok {
		c.mu.Unlock()
		return v
	}
	c.mu.Unlock()

	v := db.Query(key)
	c.mu.Lock()
	c.data[key] = v
	c.mu.Unlock()
	return v
}

func run(concurrency int, guarded bool) (int64, time.Duration) {
	db := &DB{}
	var wg sync.WaitGroup
	start := time.Now()

	if guarded {
		c := NewGuardedCache()
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); c.Get(db, "user:1") }()
		}
	} else {
		cache := map[string]string{}
		var mu sync.Mutex
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); unprotectedGet(db, cache, &mu, "user:1") }()
		}
	}
	wg.Wait()
	return db.Hits(), time.Since(start)
}

func main() {
	const n = 50
	hits1, d1 := run(n, false)
	fmt.Printf("无保护：%d 并发 -> 回源 %d 次，耗时 %v\n", n, hits1, d1.Round(time.Millisecond))

	hits2, d2 := run(n, true)
	fmt.Printf("互斥+双重检查：%d 并发 -> 回源 %d 次，耗时 %v\n", n, hits2, d2.Round(time.Millisecond))

	fmt.Println("结论：缓存击穿的本质是「同一热 key 失效瞬间的并发回源」，不是缓存本身没命中")
}
