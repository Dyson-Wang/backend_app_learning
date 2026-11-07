use std::any::{Any, TypeId};

trait Describe {
    fn name(&self) -> &'static str;
}

impl Describe for i32 {
    fn name(&self) -> &'static str {
        "i32"
    }
}

impl Describe for String {
    fn name(&self) -> &'static str {
        "String"
    }
}

// 运行时类型判定：Any + TypeId
fn type_name_of(v: &dyn Any) -> &'static str {
    if v.is::<i32>() {
        "i32"
    } else if v.is::<String>() {
        "String"
    } else if v.is::<Vec<i32>>() {
        "Vec<i32>"
    } else {
        "unknown"
    }
}

// macro_rules!：编译期展开的代码生成，替代「运行时反射找字段」
macro_rules! make_getter {
    ($name:ident, $field:ident) => {
        fn $name(v: &Record) -> i64 {
            v.$field
        }
    };
}

struct Record {
    id: i64,
    amount: i64,
}

make_getter!(get_id, id);
make_getter!(get_amount, amount);

// derive 是最常见的编译期元编程
#[derive(Debug, Clone, PartialEq)]
struct Order {
    no: String,
    amount: i64,
}

fn main() {
    // TypeId + Any：运行时可查类型
    let vals: Vec<Box<dyn Any>> = vec![
        Box::new(1i32),
        Box::new(String::from("x")),
        Box::new(vec![1i32]),
    ];
    for v in &vals {
        println!(
            "TypeId = {:?}，判定 = {}",
            v.type_id(),
            type_name_of(&**v)
        );
    }

    // downcast 取回具体类型
    let boxed: Box<dyn Any> = Box::new(42i32);
    if let Ok(n) = boxed.downcast::<i32>() {
        println!("downcast 成功 = {}", *n);
    }

    // trait 对象做动态分发（运行时多态）
    let items: Vec<Box<dyn Describe>> = vec![Box::new(1i32), Box::new(String::from("s"))];
    for it in &items {
        println!("动态分发 name = {}", it.name());
    }

    // 宏：编译期展开，零运行时开销
    let r = Record { id: 7, amount: 199 };
    println!("宏生成的 getter：id={} amount={}", get_id(&r), get_amount(&r));

    let o = Order {
        no: "A-1".into(),
        amount: 99,
    };
    let o2 = o.clone();
    println!("derive 生成 Debug/Clone/PartialEq：{:?}，相等 = {}", o, o == o2);

    println!("Rust 没有运行时反射：字段名、注解、动态调用都在编译期确定");
}
