package main

import (
	"fmt"
	"time"
)

type Order struct {
	ID     int
	UserID int
}

type User struct {
	ID   int
	Name string
}

type FakeDB struct {
	queries int
}

func (db *FakeDB) Orders() []Order {
	db.queries++
	orders := make([]Order, 0, 10)
	for i := 1; i <= 10; i++ {
		orders = append(orders, Order{ID: i, UserID: i%3 + 1})
	}
	return orders
}

func (db *FakeDB) UserByID(id int) User {
	db.queries++
	time.Sleep(time.Millisecond)
	return User{ID: id, Name: fmt.Sprintf("user-%d", id)}
}

// N+1：先查订单，再为每个订单查一次用户
func (db *FakeDB) OrdersNPlusOne() []string {
	orders := db.Orders()
	out := make([]string, 0, len(orders))
	for _, o := range orders {
		u := db.UserByID(o.UserID)
		out = append(out, fmt.Sprintf("order-%d by %s", o.ID, u.Name))
	}
	return out
}

// 批量：一次查订单 + 一次 IN 查询
func (db *FakeDB) OrdersBatch() []string {
	orders := db.Orders()

	db.queries++ // 模拟 SELECT * FROM users WHERE id IN (...)
	users := map[int]User{}
	for i := 1; i <= 3; i++ {
		users[i] = User{ID: i, Name: fmt.Sprintf("user-%d", i)}
	}

	out := make([]string, 0, len(orders))
	for _, o := range orders {
		out = append(out, fmt.Sprintf("order-%d by %s", o.ID, users[o.UserID].Name))
	}
	return out
}

func main() {
	db1 := &FakeDB{}
	start := time.Now()
	res1 := db1.OrdersNPlusOne()
	fmt.Printf("N+1：%d 条订单 -> 查询 %d 次，耗时 %v\n", len(res1), db1.queries, time.Since(start).Round(time.Millisecond))

	db2 := &FakeDB{}
	start = time.Now()
	res2 := db2.OrdersBatch()
	fmt.Printf("批量：%d 条订单 -> 查询 %d 次，耗时 %v\n", len(res2), db2.queries, time.Since(start).Round(time.Millisecond))

	fmt.Printf("第一条结果一致：%v\n", res1[0] == res2[0])
	fmt.Println("结论：N+1 的代价是「1 + N 次往返」，批量把 N 压成 1，网络延迟越高收益越大")
}
