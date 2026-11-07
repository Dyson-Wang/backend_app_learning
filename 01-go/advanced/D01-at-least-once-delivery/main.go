package main

import "fmt"

type Event struct {
	ID     string
	Amount int
}

// seen 为 nil 表示消费端没有幂等表
func handle(seen map[string]bool, balance *int, e Event) bool {
	if seen != nil {
		if seen[e.ID] {
			return false // 重复消息，直接 ACK 不重复处理
		}
		seen[e.ID] = true
	}
	*balance += e.Amount
	return true
}

func main() {
	// Broker 的投递序列：e2 被投递了两次（消费者 ACK 丢了）
	deliveries := []Event{
		{"e1", 10},
		{"e2", 10},
		{"e2", 10},
		{"e3", 10},
	}

	fmt.Println("== 1. 消费端没有幂等表 ==")
	balance := 100
	for _, e := range deliveries {
		handle(nil, &balance, e)
	}
	fmt.Printf("投递 %d 次，余额 100 -> %d（多加了 10）\n", len(deliveries), balance)

	fmt.Println("== 2. 消费端按消息 ID 去重 ==")
	balance = 100
	seen := map[string]bool{}
	applied := 0
	for _, e := range deliveries {
		if handle(seen, &balance, e) {
			applied++
		}
	}
	fmt.Printf("投递 %d 次，实际处理 %d 次，余额 100 -> %d\n", len(deliveries), applied, balance)

	fmt.Println("== 3. 结论 ==")
	fmt.Println("至少一次投递下，重复是常态；正确性必须由消费端幂等保证")
	fmt.Println("去重表要落库并和业务写在同一事务里，否则还是会重复")
}
