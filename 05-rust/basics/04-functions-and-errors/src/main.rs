// 04 · 函数与错误处理
// 函数用表达式返回；没有异常，失败用 Result<T, E> / Option<T> 当值传递，? 负责传播。

use std::num::ParseIntError;

// 函数体最后一个表达式即返回值（不加分号）；return 只用于提前返回
fn add(a: i32, b: i32) -> i32 {
    a + b
}

// 返回 Result：? 在 Err 时把错误原样提前返回
fn parse_and_double(s: &str) -> Result<i32, ParseIntError> {
    let n = s.parse::<i32>()?;
    Ok(n * 2)
}

// 用 match 自己消费 Result，转成可读文本
fn describe(r: Result<i32, ParseIntError>) -> String {
    match r {
        Ok(v) => format!("成功: {}", v),
        Err(e) => format!("失败: {}", e),
    }
}

// Option 表达「可能没有」，这里找第一个偶数
fn find_even(nums: &[i32]) -> Option<i32> {
    nums.iter().copied().find(|x| *x % 2 == 0)
}

fn main() {
    // 1) 普通函数 + 表达式返回
    println!("add(3, 4) = {}", add(3, 4));

    // 2) Result 成功路径
    let ok = parse_and_double("21");
    println!("parse_and_double(\"21\"): {}", describe(ok));

    // 3) Result 失败路径（不 panic，错误是一个值）
    let bad = parse_and_double("abc");
    println!("parse_and_double(\"abc\"): {}", describe(bad));

    // 4) Option + 模式匹配
    let nums = [1, 3, 4, 7];
    match find_even(&nums) {
        Some(v) => println!("find_even: 找到偶数 {}", v),
        None => println!("find_even: 没有偶数"),
    }

    // 5) unwrap_or 给「没有」一个默认值
    let empty: &[i32] = &[];
    let v = find_even(empty).unwrap_or(-1);
    println!("空切片 unwrap_or(-1) = {}", v);

    // 6) ? 与 match 的对照：两者都是把错误当值处理
    match parse_and_double("50") {
        Ok(n) => println!("? 链路成功: {}", n),
        Err(e) => println!("? 链路失败: {}", e),
    }
}
