fn check(name: &str, ok: bool, detail: &str) -> bool {
    if ok {
        println!("PASS {}", name);
        return true;
    }
    println!("FAIL {}: {}", name, detail);
    false
}

fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn divide(a: i32, b: i32) -> Result<i32, &'static str> {
    if b == 0 {
        return Err("division by zero");
    }
    Ok(a / b)
}

fn parse_age(s: &str) -> Result<i32, &'static str> {
    if s.is_empty() {
        return Err("empty");
    }
    Ok(18) // 故意不校验非数字，用来触发失败用例
}

fn main() {
    println!("== 1. 单元测试 ==");
    let mut failures = 0;
    if !check("add(1,2)=3", add(1, 2) == 3, "期望 3") {
        failures += 1;
    }
    if !check("divide(1,0) 返回 Err", divide(1, 0).is_err(), "期望 Err") {
        failures += 1;
    }
    if !check(
        "parse_age(\"abc\") 返回 Err",
        parse_age("abc").is_err(),
        "期望 Err，实际 Ok(18)",
    ) {
        failures += 1;
    }

    println!("== 2. 汇总 ==");
    println!("3 个用例：{} 通过，{} 失败", 3 - failures, failures);

    println!("== 3. 结论 ==");
    println!("最后一条是故意失败的用例，用来看失败输出长什么样");
    println!("Rust 内建 #[test]，cargo test 直接跑；集成测试放 tests/ 目录");
    println!("容器化：静态二进制 + distroless，镜像可以做到 5MB 以内");
}
