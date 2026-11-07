fn main() {
    // 算术：整数除法会截断，浮点除法要显式转换
    let a: i32 = 7;
    let b: i32 = 2;
    println!("整数 {} / {} = {}，取余 = {}", a, b, a / b, a % b);
    println!("浮点 {:.3}", a as f64 / b as f64);

    // 整数溢出有四种显式策略，生产代码必须选一种
    let max = i32::MAX;
    println!("MAX + 1 wrapping    = {}", max.wrapping_add(1));
    println!("MAX + 1 checked     = {:?}", max.checked_add(1));
    println!("MAX + 1 saturating  = {}", max.saturating_add(1));
    println!("MAX + 1 overflowing = {:?}", max.overflowing_add(1));

    // 位运算
    let flags: u8 = 0b1010_0110;
    println!("flags = {:#010b}", flags);
    println!("低 4 位 = {:#06b}，左移 1 位 = {:#010b}", flags & 0b0000_1111, flags << 1);
    println!("右移 2 位 = {:#010b}，按位取反 = {:#010b}", flags >> 2, !flags);

    // 逻辑运算短路求值：右侧可能根本不执行
    let mut calls = 0;
    {
        let mut probe = |v: bool| { calls += 1; v };
        let _ = false && probe(true);
        let _ = true || probe(true);
    }
    println!("短路求值后右侧被调用次数 = {}", calls);

    // 表达式块：块的值就是最后一个表达式（没有 return）
    let tax = {
        let base = 100;
        base * 13 / 100
    };
    println!("表达式块 tax = {}", tax);

    // 类型转换必须显式 as，可能截断且不报错
    let big: i64 = 300;
    println!("i64 300 as u8 = {}（截断而非报错）", big as u8);

    // 复合赋值
    let mut n = 10;
    n += 5;
    n *= 2;
    n -= 4;
    n /= 2;
    println!("复合赋值后 n = {}", n);

    // 比较返回 bool，Rust 不允许把整数当条件
    println!("3 < 5 = {}，3 == 3 = {}", 3 < 5, 3 == 3);
}
