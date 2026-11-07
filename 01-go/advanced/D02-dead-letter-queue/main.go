package main

import "fmt"

type Msg struct {
	ID        string
	Failures  int  // 前 N 次投递失败（可重试）
	Permanent bool // 永久错误，重试也没用
}

func deliver(m Msg, maxAttempts int) (bool, int) {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if m.Permanent {
			fmt.Printf("%s -> 永久错误，直接进死信（不重试）\n", m.ID)
			return false, attempt
		}
		if attempt > m.Failures {
			fmt.Printf("%s -> 成功（第 %d 次）\n", m.ID, attempt)
			return true, attempt
		}
		if attempt < maxAttempts {
			backoff := 10 << (attempt - 1) // 10ms, 20ms, 40ms
			fmt.Printf("%s -> 可重试错误，退避 %dms 后重试\n", m.ID, backoff)
			continue
		}
		fmt.Printf("%s -> 重试耗尽，进死信\n", m.ID)
	}
	return false, maxAttempts
}

func main() {
	msgs := []Msg{
		{"m1", 0, false},
		{"m2", 2, false},
		{"m3", 0, true},
		{"m4", 99, false},
		{"m5", 0, false},
	}

	fmt.Println("== 1. 逐条投递（最多 3 次） ==")
	ok, dead, total := 0, 0, 0
	for _, m := range msgs {
		success, attempts := deliver(m, 3)
		total += attempts
		if success {
			ok++
		} else {
			dead++
		}
	}

	fmt.Println("== 2. 统计 ==")
	fmt.Printf("成功 = %d，死信 = %d，总投递 = %d\n", ok, dead, total)

	fmt.Println("== 3. 结论 ==")
	fmt.Println("永久错误不要重试，重试只会拖慢整个队列")
	fmt.Println("死信队列必须有人看：它进 DLQ 说明要么代码有 bug，要么数据有问题")
}
