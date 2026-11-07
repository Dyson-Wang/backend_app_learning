package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// State 是熔断器的三种状态。
type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

func (s State) String() string {
	switch s {
	case Closed:
		return "closed"
	case Open:
		return "open"
	case HalfOpen:
		return "half-open"
	}
	return "unknown"
}

// Breaker 用连续失败次数触发熔断，冷却后放一个探测请求。
type Breaker struct {
	mu        sync.Mutex
	state     State
	failures  int
	threshold int
	cooldown  time.Duration
	openedAt  time.Time
	now       func() time.Time // 可注入时钟，测试不用真睡
}

func NewBreaker(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{threshold: threshold, cooldown: cooldown, now: time.Now}
}

// Allow 决定这次请求能不能打到下游。
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case Open:
		if b.now().Sub(b.openedAt) >= b.cooldown {
			b.state = HalfOpen
			return nil // 放一个探测请求过去
		}
		return errors.New("circuit open: fast fail")
	case HalfOpen:
		return errors.New("circuit half-open: probe already in flight")
	default:
		return nil
	}
}

func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = Closed
}

func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.state == HalfOpen || b.failures >= b.threshold {
		b.state = Open
		b.openedAt = b.now()
	}
}

func (b *Breaker) Current() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

func call(b *Breaker, ok bool) string {
	if err := b.Allow(); err != nil {
		return "rejected: " + err.Error()
	}
	if ok {
		b.Success()
		return "ok"
	}
	b.Failure()
	return "downstream error"
}

func main() {
	clock := time.Unix(0, 0)
	b := NewBreaker(5, 200*time.Millisecond)
	b.now = func() time.Time { return clock }

	fmt.Println("== 1. 正常请求 ==")
	fmt.Printf("call -> %s, state=%s\n", call(b, true), b.Current())

	fmt.Println("== 2. 连续 5 次失败触发熔断 ==")
	for i := 1; i <= 5; i++ {
		fmt.Printf("failure #%d -> %s, state=%s\n", i, call(b, false), b.Current())
	}

	fmt.Println("== 3. 熔断期间快速失败，不再打下游 ==")
	fmt.Printf("call -> %s, state=%s\n", call(b, true), b.Current())

	fmt.Println("== 4. 冷却 200ms 后进入半开 ==")
	clock = clock.Add(200 * time.Millisecond)
	fmt.Printf("allow -> %v, state=%s\n", b.Allow(), b.Current())

	fmt.Println("== 5. 探测成功恢复 closed ==")
	b.Success()
	fmt.Printf("state=%s, 后续请求 -> %s\n", b.Current(), call(b, true))

	fmt.Println("== 6. 半开探测失败则重新 open ==")
	for i := 0; i < 5; i++ {
		b.Failure()
	}
	clock = clock.Add(200 * time.Millisecond)
	_ = b.Allow() // 进入 half-open
	b.Failure()
	fmt.Printf("state=%s（探测失败后立即回到 open）\n", b.Current())
}
