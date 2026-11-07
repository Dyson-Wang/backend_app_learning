package main

import "fmt"

func main() {
	const clients = 100
	const downstreamCap = 100

	fmt.Println("== 1. 无退避：失败立即重试 3 次 ==")
	noBackoff := clients * (1 + 3)
	fmt.Printf("下游收到 %d QPS（容量 %d，超载 %d 倍）\n", noBackoff, downstreamCap, noBackoff/downstreamCap)

	fmt.Println("== 2. 指数退避 + 全抖动 ==")
	// 3 次重试分散到 1s / 2s / 4s 三个时间窗，每个窗口只承接 1/3 的客户端
	peak := clients + clients/3
	fmt.Printf("下游峰值 ≈ %d QPS（接近容量，不再放大）\n", peak)

	fmt.Println("== 3. 再加熔断 ==")
	fmt.Printf("下游请求 = 0 QPS（open 状态本地快速失败，直接返回降级结果）\n")

	fmt.Println("== 4. 结论 ==")
	fmt.Println("重试必须有上限、退避和抖动；三个都没有，就是给故障加杠杆")
	fmt.Println("熔断才是止损的关键：下游已经挂了，再发请求只会让它更晚恢复")
}
