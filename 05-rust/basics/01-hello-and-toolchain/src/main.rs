use std::thread;

fn main() {
    println!("hello, backend roadmap");
    println!("os:   {}", std::env::consts::OS);
    println!("arch: {}", std::env::consts::ARCH);
    match thread::available_parallelism() {
        Ok(n) => println!("cpus: {}", n),
        Err(_) => println!("cpus: unknown"),
    }

    // 变量默认不可变，要改必须显式 mut
    let x = 5;
    // x = 6; // 取消注释会编译失败：cannot assign twice to immutable variable
    let mut y = 5;
    y += 1;
    println!("immutable x={}, mutable y={}", x, y);

    // 所有权移动：s 在这一行之后失效
    let s = String::from("hello");
    let t = s;
    // println!("{}", s); // 编译失败：borrow of moved value
    println!("moved value: {}", t);

    // 借用不转移所有权
    let borrowed = &t;
    println!("borrowed: {}", borrowed);
    println!("t is still usable: {}", t);
}
