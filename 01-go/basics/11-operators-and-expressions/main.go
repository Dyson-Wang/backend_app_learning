package main

import (
	"fmt"
	"strconv"
)

// iif 模拟三目表达式：Go 没有 ?:，但用泛型函数可以写出等价的表达式
func iif[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func main() {
	// 1. 算术与整数除法：整数相除截断，负数向零取整
	a, b := 7, 2
	fmt.Printf("算术: %d + %d = %d, %d / %d = %d, %d %% %d = %d\n", a, b, a+b, a, b, a/b, a, b, a%b)
	fmt.Printf("负数整除向零取整: -7/2 = %d, -7%%2 = %d\n", -7/2, -7%2)

	// 2. 优先级：先乘除后加减，可用括号改变
	fmt.Printf("优先级: 2+3*4 = %d, (2+3)*4 = %d\n", 2+3*4, (2+3)*4)

	// 3. 比较与逻辑：结果都是 bool
	fmt.Printf("比较/逻辑: 3>2 = %t, (3>2)&&(2<1) = %t, !(1==1) = %t\n",
		3 > 2, (3 > 2) && (2 < 1), !(1 == 1))

	// 4. 短路求值：右侧带副作用的函数不会被调用
	calls := 0
	side := func() bool { calls++; return true }
	_ = false && side()
	_ = true || side()
	fmt.Printf("短路求值: 右侧函数实际调用次数 = %d（应为 0）\n", calls)

	// 5. 位运算
	x := 0b1100
	fmt.Printf("位运算: %04b & %04b = %04b, | = %04b, ^ = %04b, <<1 = %04b\n",
		x, 0b1010, x&0b1010, x|0b1010, x^0b1010, x<<1)

	// 6. 赋值与复合赋值
	n := 10
	n += 5
	n <<= 1
	n &= 0b11111
	fmt.Printf("复合赋值: n 最终 = %d\n", n)

	// 7. 整数溢出：按二进制补码回绕，是确定的、不 panic
	var u8 uint8 = 255
	u8++
	fmt.Printf("无符号溢出: 255+1 (uint8) = %d\n", u8)
	var i8 int8 = 127
	i8++
	fmt.Printf("有符号溢出: 127+1 (int8) = %d\n", i8)

	// 8. 类型转换：必须显式，不做隐式提升；float->int 截断，int<->string 走 strconv
	f := 3.99
	fmt.Printf("显式转换: float64(%d) = %v, int(%v) = %d（截断）\n", a, float64(a), f, int(f))
	s := strconv.Itoa(42)
	back, err := strconv.Atoi(s)
	fmt.Printf("int<->string: Itoa(42) = %q, Atoi 回读 = %d, err = %v\n", s, back, err)

	// 9. 没有三目运算符，用 iif 泛型函数或普通 if 表达
	fmt.Printf("等价三目: iif(5>3, \"大\", \"小\") = %s\n", iif(5 > 3, "大", "小"))
}
