package main

import "fmt"

// 简化模型：所有请求同时到达，按到达顺序 FIFO 抢池子里的连接
func simulate(durations []int, pool, waitTimeout int) (int, int, int) {
	freeAt := make([]int, pool)
	ok, timeout, total := 0, 0, 0
	for _, d := range durations {
		best := 0
		for i := 1; i < pool; i++ {
			if freeAt[i] < freeAt[best] {
				best = i
			}
		}
		wait := freeAt[best]
		if wait > waitTimeout {
			timeout++
			continue
		}
		freeAt[best] = wait + d
		if freeAt[best] > total {
			total = freeAt[best]
		}
		ok++
	}
	return ok, timeout, total
}

func main() {
	normal := make([]int, 10)
	for i := range normal {
		normal[i] = 30
	}

	fmt.Println("== 1. 正常：全部 30ms ==")
	ok, to, total := simulate(normal, 2, 100)
	fmt.Printf("成功 = %d，超时 = %d，总耗时 = %dms\n", ok, to, total)

	fmt.Println("== 2. 混入 1 个 300ms 慢查询（等待上限 100ms） ==")
	slow := append([]int(nil), normal...)
	slow[0] = 300
	ok, to, total = simulate(slow, 2, 100)
	fmt.Printf("成功 = %d，超时 = %d，总耗时 = %dms\n", ok, to, total)

	fmt.Println("== 3. 结论 ==")
	fmt.Println("池子只有 2 个连接，1 个慢查询就占掉一半；等待超时把雪崩挡在入口")
	fmt.Println("慢查询必须配 statement_timeout + 索引，光加连接数只是把问题推后")
}
