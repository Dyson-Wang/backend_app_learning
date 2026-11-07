// 03 · 控制流：分支、循环、迭代
// if / match 是表达式，loop / while / for 三种循环，配合 break 带值与标签跳出。

fn classify(n: i32) -> &'static str {
    // if 是表达式：每个分支类型必须一致，且必须覆盖所有情况
    if n < 0 {
        "负数"
    } else if n == 0 {
        "零"
    } else {
        "正数"
    }
}

fn main() {
    // 1) if 作为表达式直接赋值
    let n = -5;
    println!("if 表达式: {} 属于 {}", n, classify(n));

    // 2) match：穷尽性由编译器检查，支持区间和「或」模式
    let code = 404;
    let desc = match code {
        200 => "OK",
        400 | 401 | 403 => "客户端错误",
        404 => "未找到",
        500..=599 => "服务端错误",
        _ => "未知状态",
    };
    println!("match 分支: {} -> {}", code, desc);

    // 3) loop + break 带值：break 的值就是整个 loop 表达式的结果
    let mut i = 0;
    let result = loop {
        i += 1;
        if i * i > 50 {
            break i * i;
        }
    };
    println!("loop break 带值: 首个平方超过 50 的是 {}", result);

    // 4) while 循环求和
    let mut sum = 0;
    let mut k = 1;
    while k <= 5 {
        sum += k;
        k += 1;
    }
    println!("while 求和 1..=5 = {}", sum);

    // 5) for + 区间 + 迭代器适配
    let evens: Vec<i32> = (1..=10).filter(|x| *x % 2 == 0).collect();
    println!("for/区间偶数: {:?}", evens);

    // 6) 带标签的循环：break 'outer 直接跳出外层
    let mut found = (0, 0);
    'outer: for a in 1..5 {
        for b in 1..5 {
            if a * b == 6 {
                found = (a, b);
                break 'outer;
            }
        }
    }
    println!("标签 break: 找到 a*b==6 -> {:?}", found);

    // 7) if let：只关心一个模式时的简写
    let maybe = Some(42);
    if let Some(v) = maybe {
        println!("if let: 取出 Some({})", v);
    }
}
