package main

import "fmt"

type Event struct {
	ID   string
	Kind string // create / pay / ship
}

// 状态机：只接受合法迁移，顺序错了就报错
func apply(state string, e Event) (string, string) {
	switch {
	case e.Kind == "create" && state == "":
		return "CREATED", ""
	case e.Kind == "pay" && state == "CREATED":
		return "PAID", ""
	case e.Kind == "ship" && state == "PAID":
		return "SHIPPED", ""
	default:
		return state, fmt.Sprintf("%s 在 %s 状态下不合法", e.Kind, stateOrInit(state))
	}
}

func stateOrInit(s string) string {
	if s == "" {
		return "INIT"
	}
	return s
}

func consume(events []Event) (string, string) {
	state := ""
	for _, e := range events {
		var err string
		state, err = apply(state, e)
		if err != "" {
			return state, err
		}
	}
	return state, ""
}

func main() {
	create := Event{"c1", "create"}
	pay := Event{"c2", "pay"}
	ship := Event{"c3", "ship"}

	fmt.Println("== 1. 同一 key 被并发消费（顺序打乱） ==")
	state, err := consume([]Event{ship, create, pay})
	fmt.Printf("最终状态 = %s，错误 = %s\n", state, err)

	fmt.Println("== 2. 按 key 路由到同一分区、单消费者串行处理 ==")
	state, err = consume([]Event{create, pay, ship})
	fmt.Printf("最终状态 = %s，错误 = %q\n", state, err)

	fmt.Println("== 3. 分区函数：同一 order_id 必然落到同一分区 ==")
	for _, k := range []string{"order-42", "order-43", "order-42"} {
		h := 0
		for _, ch := range k {
			h = (h*31 + int(ch)) % 4
		}
		fmt.Printf("%s -> 分区 %d\n", k, h)
	}

	fmt.Println("== 4. 结论 ==")
	fmt.Println("顺序性不是全局的，是「同一分区内」的；跨 key 的顺序本来就没有意义")
	fmt.Println("要顺序就牺牲并发：单分区单消费者，或者消费端按 key 做串行队列")
}
