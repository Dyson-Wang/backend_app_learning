// 返回闭包：必须 move，否则返回的闭包会引用已销毁的栈变量
fn make_counter() -> impl FnMut() -> i32 {
    let mut count = 0;
    move || {
        count += 1;
        count
    }
}

fn apply_twice<F: Fn(i32) -> i32>(f: F, x: i32) -> i32 {
    f(f(x))
}

fn call_once<F: FnOnce() -> String>(f: F) -> String {
    f()
}

fn main() {
    let mut c = make_counter();
    let a1 = c();
    let a2 = c();
    let a3 = c();
    println!("闭包持有状态：{} {} {}", a1, a2, a3);

    // 捕获不可变引用
    let base = 10;
    let add_base = |x: i32| x + base;
    println!("捕获不可变引用：add_base(5) = {}", add_base(5));

    // 捕获可变引用
    let mut total = 0;
    {
        let mut acc = |x: i32| total += x;
        acc(3);
        acc(4);
    }
    println!("捕获可变引用：total = {}", total);

    // move：把所有权搬进闭包
    let owned = String::from("moved");
    let consume = move || owned.len();
    println!("move 闭包：len = {}", consume());

    println!("高阶函数 apply_twice(|x| x*3, 2) = {}", apply_twice(|x| x * 3, 2));

    // 三种闭包 trait：Fn / FnMut / FnOnce
    println!("FnOnce = {}", call_once(|| "只能调用一次".to_string()));

    // 无捕获闭包可以强转成函数指针
    let f: fn(i32) -> i32 = |x| x + 1;
    println!("函数指针调用 f(41) = {}", f(41));

    // 闭包 + 迭代器组合
    let nums = vec![1, 2, 3, 4, 5];
    let sum: i32 = nums.iter().filter(|&&x| x % 2 == 1).map(|&x| x * 10).sum();
    println!("filter + map + sum = {}", sum);
}
