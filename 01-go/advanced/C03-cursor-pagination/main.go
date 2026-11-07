package main

import (
	"fmt"
	"sort"
)

type Item struct {
	ID    int
	Score int
}

func ordered(items []Item) []Item {
	out := append([]Item(nil), items...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Score != out[j].Score {
			return out[i].Score > out[j].Score
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func offsetPage(items []Item, offset, limit int) []Item {
	s := ordered(items)
	if offset >= len(s) {
		return nil
	}
	end := offset + limit
	if end > len(s) {
		end = len(s)
	}
	return s[offset:end]
}

// 游标：只取「排在上次最后一条之后」的数据
func cursorPage(items []Item, after *Item, limit int) []Item {
	s := ordered(items)
	out := make([]Item, 0, limit)
	for _, it := range s {
		if after != nil && !(it.Score < after.Score || (it.Score == after.Score && it.ID > after.ID)) {
			continue
		}
		out = append(out, it)
		if len(out) == limit {
			break
		}
	}
	return out
}

func ids(items []Item) []int {
	out := make([]int, len(items))
	for i, it := range items {
		out[i] = it.ID
	}
	return out
}

func overlap(a, b []Item) int {
	seen := map[int]bool{}
	for _, it := range a {
		seen[it.ID] = true
	}
	n := 0
	for _, it := range b {
		if seen[it.ID] {
			n++
		}
	}
	return n
}

func main() {
	base := make([]Item, 10)
	for i := 0; i < 10; i++ {
		base[i] = Item{ID: i + 1, Score: 100 - i*5}
	}

	fmt.Println("== 1. OFFSET 分页：翻页期间有新数据插入 ==")
	p1 := offsetPage(base, 0, 3)
	grown := append(append([]Item(nil), base...), Item{ID: 11, Score: 97})
	p2 := offsetPage(grown, 3, 3)
	fmt.Printf("第一页 = %v\n", ids(p1))
	fmt.Printf("插入一条 Score=97 后，第二页 = %v，与第一页重复 %d 条\n", ids(p2), overlap(p1, p2))

	fmt.Println("== 2. 游标分页：用 (score, id) 做锚点 ==")
	c1 := cursorPage(base, nil, 3)
	last := c1[len(c1)-1]
	c2 := cursorPage(grown, &last, 3)
	fmt.Printf("第一页 = %v，锚点 = (score=%d, id=%d)\n", ids(c1), last.Score, last.ID)
	fmt.Printf("插入同一条后，第二页 = %v，与第一页重复 %d 条\n", ids(c2), overlap(c1, c2))

	fmt.Println("== 3. 结论 ==")
	fmt.Println("OFFSET 按「位置」翻页，前面插入数据会让整体后移，出现重复或漏读")
	fmt.Println("游标按「值」翻页，锚点不变结果就不变；代价是不能跳页，只能上一页/下一页")
}
