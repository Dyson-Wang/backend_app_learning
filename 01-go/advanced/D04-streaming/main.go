package main

import "fmt"

// 用逻辑时钟代替真实 sleep，保证输出可复现
type Clock struct{ t int }

func (c *Clock) tick(ms int) int { c.t += ms; return c.t }

func main() {
	const chunks = 5
	const perChunk = 20

	fmt.Println("== 1. 缓冲模式：全部生成完再发 ==")
	c := &Clock{}
	firstByte := 0
	for i := 0; i < chunks; i++ {
		c.tick(perChunk) // 生成第 i 块
	}
	firstByte = c.t // 最后一块生成完才开始发
	fmt.Printf("首字节 = %dms，总耗时 = %dms\n", firstByte, c.t)

	fmt.Println("== 2. 流式模式：生成一块发一块 ==")
	c = &Clock{}
	firstByte = 0
	for i := 0; i < chunks; i++ {
		c.tick(perChunk)
		if i == 0 {
			firstByte = c.t // 第一块生成完立刻发
		}
	}
	fmt.Printf("首字节 = %dms，总耗时 = %dms\n", firstByte, c.t)

	fmt.Println("== 3. 对比 ==")
	fmt.Println("总耗时一样，首字节差了 4 倍：用户感知的是首字节，不是总耗时")
	fmt.Println("SSE / chunked / WebSocket 的价值就在这里；代价是不能压缩、不能算 Content-Length")

	fmt.Println("== 4. 结论 ==")
	fmt.Println("大响应先想清楚：客户端要等全量，还是能边收边渲染")
	fmt.Println("流式要处理背压：客户端读得慢，服务端要么阻塞要么丢，不能无限缓冲")
}
