package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

var ErrConflict = errors.New("request with this key is still in flight")

type Order struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}

type entry struct {
	done  bool
	order Order
}

// IdempotencyStore 用互斥锁模拟数据库唯一索引的原子性。
// 生产里这层必须是数据库的唯一约束，进程内锁在多实例下无效。
type IdempotencyStore struct {
	mu      sync.Mutex
	entries map[string]*entry
	created int
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{entries: make(map[string]*entry)}
}

// Claim 把「查 + 占位」合成一个原子操作，对应 INSERT ... ON CONFLICT DO NOTHING
func (s *IdempotencyStore) Claim(key string) (Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[key]; ok {
		if !e.done {
			return Order{}, false, ErrConflict
		}
		return e.order, true, nil // 已完成：直接回放上次结果
	}
	s.entries[key] = &entry{}
	return Order{}, false, nil
}

func (s *IdempotencyStore) Complete(key string, amount int) Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.created++
	order := Order{ID: fmt.Sprintf("ORD-%d", s.created), Amount: amount}
	e := s.entries[key]
	e.order = order
	e.done = true
	return order
}

func (s *IdempotencyStore) CreatedCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.created
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func createOrderHandler(store *IdempotencyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Idempotency-Key header is required", http.StatusBadRequest)
			return
		}

		order, cached, err := store.Claim(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		if cached {
			writeJSON(w, http.StatusOK, order) // 回放，不重复创建
			return
		}

		time.Sleep(5 * time.Millisecond) // 放大并发窗口，制造重复提交
		order = store.Complete(key, 100)
		writeJSON(w, http.StatusCreated, order)
	}
}

func post(client *http.Client, url, key string) (int, string) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return 0, ""
	}
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, ""
	}
	defer resp.Body.Close()
	var o Order
	json.NewDecoder(resp.Body).Decode(&o)
	return resp.StatusCode, o.ID
}

func runSelfTest() {
	store := NewIdempotencyStore()
	srv := httptest.NewServer(createOrderHandler(store))
	defer srv.Close()
	client := &http.Client{Timeout: 3 * time.Second}

	// case 1: 20 个并发携带相同 key，遇到 409 就短暂重试，最终应只创建 1 单
	const n = 20
	var wg sync.WaitGroup
	statuses := make(chan int, n)
	ids := make(chan string, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for attempt := 0; attempt < 5; attempt++ {
				status, id := post(client, srv.URL, "same-key")
				if status == http.StatusConflict {
					time.Sleep(10 * time.Millisecond)
					continue
				}
				statuses <- status
				ids <- id
				return
			}
			statuses <- 0
			ids <- ""
		}()
	}
	wg.Wait()
	close(statuses)
	close(ids)

	created, replay := 0, 0
	unique := map[string]bool{}
	for s := range statuses {
		switch s {
		case http.StatusCreated:
			created++
		case http.StatusOK:
			replay++
		}
	}
	for id := range ids {
		if id != "" {
			unique[id] = true
		}
	}
	fmt.Printf("20 concurrent POST same key -> 201:%d 200:%d distinct orders:%d\n",
		created, replay, len(unique))
	if created != 1 || replay != n-1 || len(unique) != 1 {
		log.Fatal("FAIL: idempotency broken")
	}

	// case 2: 不同 key 应创建不同订单
	_, id1 := post(client, srv.URL, "another-key")
	_, id2 := post(client, srv.URL, "yet-another-key")
	fmt.Printf("different keys -> %s / %s\n", id1, id2)
	if id1 == id2 || id1 == "" || id2 == "" {
		log.Fatal("FAIL: different keys should create different orders")
	}

	// case 3: 缺 key 直接拒绝
	status, _ := post(client, srv.URL, "")
	fmt.Println("missing key ->", status)
	if status != http.StatusBadRequest {
		log.Fatal("FAIL: missing key should be rejected")
	}

	fmt.Println("total orders created:", store.CreatedCount())
	fmt.Println("selftest PASS")
}

func main() {
	runSelfTest()
}
