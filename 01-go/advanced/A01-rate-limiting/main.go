package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"
)

// TokenBucket: 固定速率补充令牌，允许突发流量但限制长期平均速率
type TokenBucket struct {
	mu     sync.Mutex
	rate   float64 // 每秒补充的令牌数
	burst  float64 // 桶容量，也是允许的瞬时突发
	tokens float64
	last   time.Time
}

func NewTokenBucket(rate, burst float64) *TokenBucket {
	return &TokenBucket{
		rate:   rate,
		burst:  burst,
		tokens: burst,
		last:   time.Now(),
	}
}

// Allow 必须在锁内完成「补充 + 判断 + 扣减」，否则高并发下会超发
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now
	b.tokens += elapsed * b.rate
	if b.tokens > b.burst {
		b.tokens = b.burst
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func RateLimit(next http.Handler, bucket *TokenBucket) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !bucket.Allow() {
			w.Header().Set("Retry-After", "1")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintln(w, `{"code":"RATE_LIMITED","message":"too many requests"}`)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newHandler(bucket *TokenBucket) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	return RateLimit(mux, bucket)
}

func runSelfTest() {
	// rate=1/s, burst=5：瞬间最多放 5 个，之后每秒只补 1 个
	bucket := NewTokenBucket(1, 5)
	srv := httptest.NewServer(newHandler(bucket))
	defer srv.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	var wg sync.WaitGroup
	var allowed, limited int32

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Get(srv.URL + "/api")
			if err != nil {
				return
			}
			defer resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt32(&allowed, 1)
			case http.StatusTooManyRequests:
				atomic.AddInt32(&limited, 1)
			}
		}()
	}
	wg.Wait()

	fmt.Printf("burst=5 rate=1/s, 20 concurrent -> 200:%d 429:%d\n",
		atomic.LoadInt32(&allowed), atomic.LoadInt32(&limited))
	if allowed < 5 || allowed > 6 || allowed+limited != 20 {
		log.Fatal("FAIL: burst should admit about 5 requests")
	}

	// 等令牌补充后应能再次放行
	time.Sleep(1100 * time.Millisecond)
	resp, err := client.Get(srv.URL + "/api")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("after 1.1s refill ->", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("FAIL: token refill did not work")
	}
	fmt.Println("selftest PASS")
}

func main() {
	serve := flag.Bool("serve", false, "start a real server on :8081 instead of selftest")
	flag.Parse()

	if !*serve {
		runSelfTest()
		return
	}

	log.Println("listening on :8081, try: for i in $(seq 1 10); do curl -s -o /dev/null -w '%{http_code}\\n' http://localhost:8081/api; done")
	log.Fatal(http.ListenAndServe(":8081", newHandler(NewTokenBucket(5, 5))))
}
