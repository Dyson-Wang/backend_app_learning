package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var ErrPermanent = errors.New("permanent failure")

type RetryableFunc func(ctx context.Context) error

// DoWithRetry: 指数退避 + 全抖动 + 区分可重试错误 + 遵守 ctx 取消
func DoWithRetry(ctx context.Context, maxAttempts int, baseDelay time.Duration, fn RetryableFunc) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn(ctx)
		if lastErr == nil {
			fmt.Printf("  attempt %d: success\n", attempt)
			return nil
		}
		if errors.Is(lastErr, ErrPermanent) {
			fmt.Printf("  attempt %d: permanent error, stop immediately\n", attempt)
			return lastErr
		}
		if attempt == maxAttempts {
			break
		}

		// 指数退避 ceiling = base * 2^(attempt-1)，实际延迟取 [0, ceiling] 随机值
		// 全抖动让大量客户端的重试时刻错开，避免同步化的重试风暴
		ceiling := baseDelay * time.Duration(1<<(attempt-1))
		delay := time.Duration(rand.Int63n(int64(ceiling) + 1))
		fmt.Printf("  attempt %d: retryable err=%v, sleep %s\n", attempt, lastErr, delay)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", maxAttempts, lastErr)
}

func runSelfTest() {
	fmt.Println("case 1: transient failure, retry until success")
	calls := 0
	err := DoWithRetry(context.Background(), 5, 10*time.Millisecond,
		func(ctx context.Context) error {
			calls++
			if calls < 3 {
				return errors.New("downstream 503")
			}
			return nil
		})
	if err != nil || calls != 3 {
		log.Fatalf("FAIL case1: calls=%d err=%v", calls, err)
	}

	fmt.Println("case 2: permanent error must NOT be retried")
	calls = 0
	err = DoWithRetry(context.Background(), 5, 10*time.Millisecond,
		func(ctx context.Context) error {
			calls++
			return fmt.Errorf("bad request: %w", ErrPermanent)
		})
	if !errors.Is(err, ErrPermanent) || calls != 1 {
		log.Fatalf("FAIL case2: calls=%d err=%v", calls, err)
	}
	fmt.Println("  calls:", calls, "(correct: no retry on 4xx)")

	fmt.Println("case 3: context cancellation stops retry immediately")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	calls = 0
	start := time.Now()
	err = DoWithRetry(ctx, 10, 50*time.Millisecond,
		func(ctx context.Context) error {
			calls++
			return errors.New("still down")
		})
	elapsed := time.Since(start)
	fmt.Printf("  calls=%d elapsed=%s err=%v\n", calls, elapsed.Round(time.Millisecond), err)
	if !errors.Is(err, context.DeadlineExceeded) {
		log.Fatal("FAIL case3: expected context deadline exceeded")
	}
	if elapsed > 200*time.Millisecond {
		log.Fatal("FAIL case3: retry ignored cancellation")
	}

	fmt.Println("selftest PASS")
}

func main() {
	runSelfTest()
}
