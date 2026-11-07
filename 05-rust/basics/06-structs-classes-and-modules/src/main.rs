// 06 · 结构体 / 类、接口、模块
// struct 定义数据，impl 定义方法，trait 充当接口，mod 划分命名空间。

// 模块：pub 的项才能被模块外访问
mod geometry {
    #[derive(Debug)]
    pub struct Point {
        pub x: f64,
        pub y: f64,
    }

    // trait 就是其它语言的「接口」，可带默认方法
    pub trait Shape {
        fn area(&self) -> f64;   // 必须实现
        fn name(&self) -> &str;  // 必须实现
        fn describe(&self) -> String {
            format!("{} 面积={:.2}", self.name(), self.area())
        }
    }

    #[derive(Debug)]
    pub struct Circle {
        pub center: Point,
        pub radius: f64,
    }

    impl Circle {
        // 关联函数（无 self），相当于静态方法 / 构造器
        pub fn new(x: f64, y: f64, r: f64) -> Circle {
            Circle {
                center: Point { x, y },
                radius: r,
            }
        }
    }

    impl Shape for Circle {
        fn area(&self) -> f64 {
            std::f64::consts::PI * self.radius * self.radius
        }
        fn name(&self) -> &str {
            "圆"
        }
    }

    pub struct Pair(pub i32, pub i32); // 元组结构体
    pub struct Marker;                 // 单元结构体
}

use geometry::{Circle, Marker, Pair, Point, Shape};

fn main() {
    // 1) 具名结构体 + 字段初始化
    let p = Point { x: 1.0, y: 2.0 };
    println!("Point: {:?}", p);

    // 2) 关联函数构造 + 方法调用
    let c = Circle::new(0.0, 0.0, 2.0);
    println!("Circle: {:?}", c);
    println!("方法 area = {:.2}", c.area());

    // 3) trait 默认方法（接口默认实现）
    println!("trait 默认方法: {}", c.describe());

    // 4) trait 对象 dyn Shape：运行时多态
    let shapes: Vec<&dyn Shape> = vec![&c];
    for s in shapes {
        println!("dyn Shape: {}", s.describe());
    }

    // 5) 元组结构体 / 单元结构体
    let pair = Pair(3, 7);
    let _m = Marker;
    println!("元组结构体: Pair({}, {})", pair.0, pair.1);
}
