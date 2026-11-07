package main

import "fmt"

// 接口：隐式实现，类型不需要写 implements
type Storer interface {
	Get(key string) (string, bool)
	Put(key, value string)
}

type MemoryStore struct {
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

// 指针接收者：需要修改接收者时用
func (m *MemoryStore) Put(key, value string) {
	m.data[key] = value
}

// 值接收者：只读时可以用值接收者
func (m *MemoryStore) Get(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

type Logger struct {
	Prefix string
}

func (l Logger) Log(msg string) {
	fmt.Printf("[%s] %s\n", l.Prefix, msg)
}

// 嵌入（组合）：Service 自动获得 Logger 的方法集
type Service struct {
	Logger
	store Storer
	name  string
}

func (s *Service) Set(key, value string) {
	s.store.Put(key, value)
	s.Log("set " + key)
}

func describe(x interface{}) {
	switch v := x.(type) {
	case string:
		fmt.Println("type switch: string =", v)
	case int:
		fmt.Println("type switch: int =", v)
	case Storer:
		fmt.Println("type switch: implements Storer")
	default:
		fmt.Printf("type switch: other %T\n", v)
	}
}

func main() {
	svc := &Service{
		Logger: Logger{Prefix: "svc"},
		store:  NewMemoryStore(),
		name:   "demo",
	}
	svc.Set("user:1", "alice")

	// 面向接口编程：变量类型是 Storer，不关心具体实现
	var s Storer = NewMemoryStore()
	s.Put("k", "v")
	if v, ok := s.Get("k"); ok {
		fmt.Println("interface call:", v)
	}

	describe("hello")
	describe(42)
	describe(s)
	describe(3.14)
}
