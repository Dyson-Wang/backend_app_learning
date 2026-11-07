package main

import "fmt"

type Account struct {
	balance int
	version int
}

func (a *Account) read() (int, int) { return a.balance, a.version }

// 乐观锁：只有 version 没变才写入
func (a *Account) cas(expectedVersion, newBalance int) bool {
	if a.version != expectedVersion {
		return false
	}
	a.balance = newBalance
	a.version++
	return true
}

func debit(acc *Account, amount int) (bool, int) {
	for attempt := 0; attempt <= 10; attempt++ {
		balance, version := acc.read()
		if acc.cas(version, balance-amount) {
			return true, attempt
		}
	}
	return false, 11
}

func main() {
	fmt.Println("== 1. 无锁：两个事务都基于旧值写回 ==")
	acc := &Account{balance: 100, version: 1}
	b1, v1 := acc.read()
	b2, v2 := acc.read()
	acc.balance, acc.version = b1-10, v1+1 // T1
	acc.balance, acc.version = b2-10, v2+1 // T2 覆盖 T1
	fmt.Printf("期望 80，实际 %d（丢失更新）\n", acc.balance)

	fmt.Println("== 2. 乐观锁：version 校验 + 重试 ==")
	acc2 := &Account{balance: 100, version: 1}
	staleBalance, staleVersion := acc2.read() // T2 先读到旧值
	ok1, att1 := debit(acc2, 10)              // T1 提交
	fmt.Printf("T1 ok=%v attempts=%d，balance=%d version=%d\n", ok1, att1, acc2.balance, acc2.version)

	if acc2.cas(staleVersion, staleBalance-10) {
		fmt.Println("T2 用旧 version 写入成功（不应发生）")
	} else {
		fmt.Println("T2 用旧 version 写入失败：version 已被 T1 改掉")
	}

	ok2, att2 := debit(acc2, 10)
	fmt.Printf("T2 重试 ok=%v attempts=%d，balance=%d version=%d\n", ok2, att2, acc2.balance, acc2.version)
	fmt.Printf("期望 80，实际 %d\n", acc2.balance)

	fmt.Println("== 3. 结论 ==")
	fmt.Println("悲观锁（SELECT FOR UPDATE）锁住整行，吞吐低；乐观锁读不加锁、写校验 version，冲突才重试")
	fmt.Println("隔离级别越高越安全，但锁等待和死锁概率也越高；读已提交 + 乐观锁是多数业务的折中")
}
