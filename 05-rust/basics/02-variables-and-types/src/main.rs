// 02 · 变量、类型、常量、零值
// 只用 std；演示标量类型、类型推断、常量，以及「Rust 没有隐式零值」这一特点。

const MAX_RETRY: u32 = 3;           // 编译期常量，使用处内联
static APP_NAME: &str = "brm-rust"; // 静态变量，有固定内存地址，'static 生命周期

fn main() {
    // 1) 默认不可变 + 类型推断；后缀可指定字面量类型
    let port = 8080u16;
    let ratio = 0.75f64;
    let enabled = true;
    println!("推断类型: port={} u16, ratio={} f64, enabled={} bool", port, ratio, enabled);

    // 2) 显式类型标注 + mut 才能改
    let mut count: i64 = 0;
    count += 1;
    count = count * 10 + 5;
    println!("显式 i64: count={}", count);

    // 3) 遮蔽（shadowing）：允许换类型重新绑定同名变量
    let count = "十五";
    println!("遮蔽后新类型: count={}", count);

    // 4) 没有「未初始化即零值」：Default trait 才是 Rust 的零值来源
    //  let x: i32;  println!("{}", x); // 编译错误：use of possibly-uninitialized
    let z_i32 = i32::default();
    let z_bool = bool::default();
    let z_string = String::default();
    println!("Default 零值: i32={}, bool={}, String=\"{}\"", z_i32, z_bool, z_string);

    // 5) 类型转换必须显式写 as，Rust 不做隐式数值提升
    let a: i32 = 7;
    let b: i64 = a as i64;
    let c = b + 100;
    println!("显式转换: i32 {} -> i64 {} + 100 = {}", a, b, c);

    // 6) 复合类型：元组与定长数组
    let pair: (u8, &str) = (1, "ok");
    let arr: [i32; 4] = [10, 20, 30, 40];
    println!("元组: ({}, {}), 数组: {:?}, 长度 {}", pair.0, pair.1, arr, arr.len());

    // 7) 常量与静态变量
    println!("const MAX_RETRY={}, static APP_NAME={}", MAX_RETRY, APP_NAME);
}
