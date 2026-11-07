#[derive(Debug, Clone, PartialEq)]
enum Shape {
    Circle { radius: f64 },
    Rect { w: f64, h: f64 },
    Triangle(f64, f64),
}

fn area(s: &Shape) -> f64 {
    match s {
        Shape::Circle { radius } => std::f64::consts::PI * radius * radius,
        Shape::Rect { w, h } => w * h,
        Shape::Triangle(b, h) => 0.5 * b * h,
    }
}

fn describe(s: &Shape) -> String {
    match s {
        Shape::Circle { radius } if *radius > 10.0 => format!("大圆 r={}", radius),
        Shape::Circle { .. } => "小圆".to_string(),
        Shape::Rect { w, h } if (w - h).abs() < f64::EPSILON => "正方形".to_string(),
        Shape::Rect { .. } => "矩形".to_string(),
        Shape::Triangle(..) => "三角形".to_string(),
    }
}

fn main() {
    let shapes = vec![
        Shape::Circle { radius: 1.0 },
        Shape::Circle { radius: 20.0 },
        Shape::Rect { w: 3.0, h: 3.0 },
        Shape::Rect { w: 4.0, h: 2.0 },
        Shape::Triangle(6.0, 4.0),
    ];

    // match 必须穷尽所有变体，少一个编译不过
    for s in &shapes {
        println!("{:?} -> {}，面积 {:.2}", s, describe(s), area(s));
    }

    // Option / Result 本身就是 enum
    let maybe: Option<i32> = Some(7);
    match maybe {
        Some(n) if n % 2 == 1 => println!("Some({}) 是奇数", n),
        Some(n) => println!("Some({}) 是偶数", n),
        None => println!("None"),
    }

    // if let：只关心一个分支
    if let Some(v) = maybe {
        println!("if let 解构 = {}", v);
    }

    // let else：不匹配就提前返回，适合卫语句
    let Some(x) = maybe else {
        println!("不会走到这里");
        return;
    };
    println!("let else 解构 = {}", x);

    // @ 绑定：既匹配范围又拿到值
    for n in [1, 5, 10] {
        match n {
            v @ 1..=5 => println!("{} 在 1..=5 内（绑定 v={}）", n, v),
            _ => println!("{} 超出范围", n),
        }
    }

    // matches! 宏：把 match 压成 bool
    println!("matches!(Some(3), Some(_)) = {}", matches!(maybe, Some(_)));
}
