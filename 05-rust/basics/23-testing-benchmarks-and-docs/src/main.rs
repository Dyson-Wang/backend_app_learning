use std::time::Instant;

fn add(a: i32, b: i32) -> i32 {
    a + b
}

// 手写断言：失败时给出用例名和期望/实际
fn assert_eq_i32(actual: i32, expected: i32, case: &str) {
    if actual != expected {
        panic!("[{}] 期望 {}，实际 {}", case, expected, actual);
    }
    println!("  通过：{}", case);
}

// 表驱动测试
fn run_table_tests() {
    let cases = [(1, 2, 3), (0, 0, 0), (-1, 1, 0), (100, -50, 50)];
    for (a, b, expected) in cases {
        assert_eq_i32(add(a, b), expected, &format!("{} + {}", a, b));
    }
}

fn bench(label: &str, f: impl Fn() -> i32, rounds: u32) {
    let start = Instant::now();
    let mut acc = 0i64;
    for _ in 0..rounds {
        acc += f() as i64;
    }
    let elapsed = start.elapsed();
    println!(
        "  基准 {}：{} 轮，累计 {}，耗时 {:?}",
        label, rounds, acc, elapsed
    );
}

/// 文档注释：`cargo doc` 会把它渲染成 HTML。
///
/// # 例子
///
/// ```ignore
/// assert_eq!(add(1, 2), 3);
/// ```
fn documented_add(a: i32, b: i32) -> i32 {
    a + b
}

fn main() {
    println!("表驱动测试：");
    run_table_tests();

    println!("基准测试（仅演示计时，结果随机器变化）：");
    bench("add(1,2)", || add(1, 2), 100_000);

    println!("文档注释示例：documented_add(2,3) = {}", documented_add(2, 3));
    println!("真实项目用 cargo test / cargo bench，见 README");
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn add_works() {
        assert_eq!(add(1, 2), 3);
        assert_eq!(add(-1, 1), 0);
    }

    #[test]
    fn documented_add_matches() {
        assert_eq!(documented_add(2, 3), add(2, 3));
    }
}
