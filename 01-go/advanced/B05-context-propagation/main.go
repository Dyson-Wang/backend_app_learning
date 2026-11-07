package main

import (
	"context"
	"fmt"
	"time"
)

// repo 模拟慢下游：200ms，但尊重 ctx 取消
func repo(ctx context.Context) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond):
		return "data", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// service 把 ctx 原样往下传，不新建 context.Background()
func service(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("service 提前退出: %w", err)
	}
	return repo(ctx)
}

func handler(ctx context.Context) (string, error) {
	return service(ctx)
}

func main() {
	fmt.Println("== 1. 上游超时 50ms，下游要 200ms ==")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := handler(ctx)
	elapsed := time.Since(start)

	fmt.Printf("err = %v\n", err)
	fmt.Printf("耗时 = %v（远小于 200ms，说明取消已传播到 repo）\n", elapsed.Round(10*time.Millisecond))

	fmt.Println("== 2. 客户端主动取消 ==")
	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel2()
	}()

	start2 := time.Now()
	_, err2 := handler(ctx2)
	fmt.Printf("err = %v，耗时 = %v\n", err2, time.Since(start2).Round(10*time.Millisecond))

	fmt.Println("== 3. 反例：中间层用 context.Background() 截断传播 ==")
	badCtx, badCancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer badCancel()
	start3 := time.Now()
	_, err3 := service(context.Background()) // 丢掉了 badCtx
	fmt.Printf("err = %v，耗时 = %v（取消了上游也没用）\n", err3, time.Since(start3).Round(10*time.Millisecond))

	fmt.Println("== 4. 结论 ==")
	fmt.Println("ctx 必须一路显式传递；任何一层新建 Background 都会切断取消传播")
	fmt.Println("超时应该在最外层设置一次，中间层只传不改")
}
