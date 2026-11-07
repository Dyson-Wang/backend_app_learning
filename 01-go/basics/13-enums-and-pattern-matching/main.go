package main

import "fmt"

// Status 用 iota 定义枚举：Go 没有 enum 关键字，靠自定义类型 + const 组模拟
type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusDone
	StatusFailed
)

func (s Status) String() string {
	names := [...]string{"pending", "running", "done", "failed"}
	if int(s) < 0 || int(s) >= len(names) {
		return "unknown"
	}
	return names[int(s)]
}

// Expr 是「代数数据类型」的近似：接口 + 未导出标记方法实现封闭
// sealedExpr 未导出 → 包外类型无法实现 Expr，从而分支集合可控
type Expr interface {
	Eval() int
	sealedExpr()
}

type Lit struct{ V int }
type Add struct{ L, R Expr }
type Mul struct{ L, R Expr }

func (Lit) sealedExpr() {}
func (Add) sealedExpr() {}
func (Mul) sealedExpr() {}

func (l Lit) Eval() int { return l.V }
func (a Add) Eval() int { return a.L.Eval() + a.R.Eval() }
func (m Mul) Eval() int { return m.L.Eval() * m.R.Eval() }

// describe 用 type switch 做「模式匹配」：按具体类型分支 + 解构字段
func describe(e Expr) string {
	switch v := e.(type) {
	case Lit:
		return fmt.Sprintf("字面量 %d", v.V)
	case Add:
		return fmt.Sprintf("加法(%s + %s)", describe(v.L), describe(v.R))
	case Mul:
		return fmt.Sprintf("乘法(%s * %s)", describe(v.L), describe(v.R))
	default:
		// Go 的 type switch 没有编译期穷尽性检查，只能靠 default 运行时兜底
		return "未知表达式类型"
	}
}

func main() {
	// 1. iota 枚举：取值、自增、String() 输出
	fmt.Printf("枚举: pending=%d running=%d done=%d failed=%d\n",
		StatusPending, StatusRunning, StatusDone, StatusFailed)
	fmt.Printf("枚举字符串: %s / %s\n", StatusRunning, Status(99))

	// 2. 用 iota 做位掩码枚举
	const (
		PermRead = 1 << iota
		PermWrite
		PermExec
	)
	fmt.Printf("位掩码: read=%d write=%d exec=%d all=%d\n",
		PermRead, PermWrite, PermExec, PermRead|PermWrite|PermExec)

	// 3. ADT + type switch 模式匹配：构造 (1+2)*(3+4) 并求值/描述
	expr := Mul{L: Add{L: Lit{1}, R: Lit{2}}, R: Add{L: Lit{3}, R: Lit{4}}}
	fmt.Printf("ADT 描述: %s\n", describe(expr))
	fmt.Printf("ADT 求值: %d\n", expr.Eval())

	// 4. 解构：type assertion 取出具体类型并访问字段
	if add, ok := expr.L.(Add); ok {
		fmt.Printf("解构: 左子树是 Add, L=%d, R=%d\n", add.L.(Lit).V, add.R.(Lit).V)
	}

	// 5. 穷尽性局限：未知类型落到 default，编译器不会报错
	fmt.Printf("穷尽性: default 兜底 → %s\n", describe(nil))
}
