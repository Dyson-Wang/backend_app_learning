package main

import "fmt"

func main() {
	const backlog = 1000
	const perSec = 10 // 单消费者每秒处理 10 条
	const producers = 20

	fmt.Println("== 1. 单消费者，逐条处理 ==")
	fmt.Printf("%d 条 / %d 条每秒 = %d 秒\n", backlog, perSec, backlog/perSec)

	fmt.Println("== 2. 扩到 3 个消费者 ==")
	three := (backlog + perSec*3 - 1) / (perSec * 3)
	fmt.Printf("%d 条 / %d 条每秒 = %d 秒\n", backlog, perSec*3, three)

	fmt.Println("== 3. 批量 10 条一次（单消费者） ==")
	batch := backlog / (perSec * 10)
	fmt.Printf("%d 条 / %d 条每秒 = %d 秒\n", backlog, perSec*10, batch)

	fmt.Println("== 4. 生产仍在继续（每秒 20 条） ==")
	net := perSec*10 - producers
	clear := (backlog + net - 1) / net
	fmt.Printf("净消化 %d 条每秒，清空 %d 条需要 %d 秒\n", net, backlog, clear)

	fmt.Println("== 5. 结论 ==")
	fmt.Println("积压先算「净消化速率」，不是「消费速率」——生产端还在往里灌")
	fmt.Println("扩消费者最快但受分区数限制；批量提升单机吞吐，代价是延迟和失败重放")
}
