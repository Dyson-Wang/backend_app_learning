use std::fmt::Display;

// 泛型函数 + trait bound：T 必须能比较且能复制
fn largest<T: PartialOrd + Copy>(list: &[T]) -> T {
    let mut max = list[0];
    for &item in list.iter() {
        if item > max {
            max = item;
        }
    }
    max
}

// 泛型结构体
struct Pair<T> {
    first: T,
    second: T,
}

// 条件实现：只有满足约束的 T 才有这个方法
impl<T: Display + PartialOrd + Copy> Pair<T> {
    fn cmp_display(&self) {
        if self.first >= self.second {
            println!("较大的是 first = {}", self.first);
        } else {
            println!("较大的是 second = {}", self.second);
        }
    }
}

// 自定义 trait 作为约束
trait Summary {
    fn summarize(&self) -> String;
}

impl Summary for Pair<i32> {
    fn summarize(&self) -> String {
        format!("Pair({}, {})", self.first, self.second)
    }
}

// impl Trait 作为返回类型：调用方只看到「是个 Fn」
fn make_adder(n: i32) -> impl Fn(i32) -> i32 {
    move |x| x + n
}

fn main() {
    println!("largest i32  = {}", largest(&[3, 9, 2, 7]));
    println!("largest char = {}", largest(&['a', 'z', 'm']));
    println!("largest f64  = {}", largest(&[1.5, 0.5, 2.5]));

    let p = Pair { first: 10, second: 4 };
    p.cmp_display();
    println!("自定义 trait = {}", p.summarize());

    let add5 = make_adder(5);
    println!("impl Fn 返回闭包：add5(10) = {}", add5(10));

    // 单态化：同一个泛型函数，编译器为每种 T 生成一份机器码
    println!("泛型在编译期单态化，运行时无擦除、无虚调用");
}
