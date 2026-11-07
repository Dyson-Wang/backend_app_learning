mod geometry {
    pub struct Point {
        pub x: f64,
        pub y: f64,
    }

    impl Point {
        pub fn new(x: f64, y: f64) -> Self {
            Point { x, y }
        }

        pub fn distance_to_origin(&self) -> f64 {
            (self.x * self.x + self.y * self.y).sqrt()
        }
    }

    pub mod shapes {
        pub fn area_of_square(side: f64) -> f64 {
            side * side
        }

        // pub(crate)：整个 crate 可见，但不出现在对外 API
        pub(crate) fn internal_helper() -> &'static str {
            "crate 内可见"
        }
    }

    // 默认私有：只有本模块及其子模块能调用
    fn private_helper() -> &'static str {
        "模块内私有"
    }

    pub fn describe() -> String {
        format!("{} / {}", private_helper(), shapes::internal_helper())
    }
}

mod outer {
    pub mod inner {
        pub fn hello() -> &'static str {
            "inner::hello"
        }

        pub fn call_super() -> &'static str {
            super::sibling()
        }
    }

    pub fn sibling() -> &'static str {
        "outer::sibling"
    }
}

use geometry::shapes::area_of_square as square_area;
use geometry::Point;

fn main() {
    let p = Point::new(3.0, 4.0);
    println!("Point 距离原点 = {}", p.distance_to_origin());

    // 别名导入避免命名冲突
    println!("别名导入 square_area(5) = {}", square_area(5.0));

    // 嵌套模块用路径访问
    println!("嵌套模块 = {}", outer::inner::hello());
    println!("super 访问同级 = {}", outer::inner::call_super());
    println!("模块内私有函数只能被 describe 调用 = {}", geometry::describe());

    // 下面这行会编译失败，所以注释掉：
    // geometry::private_helper();
    println!("私有项不可从外部访问：private_helper 只在 geometry 内可见");
}
