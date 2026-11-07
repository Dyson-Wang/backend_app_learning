use std::cell::RefCell;
use std::rc::Rc;

#[derive(Debug)]
struct Droppable(&'static str);

impl Drop for Droppable {
    fn drop(&mut self) {
        println!("  drop: {}", self.0);
    }
}

fn take_ownership(s: String) -> String {
    s
}

// 生命周期标注：返回的引用至少和两个输入活得一样久
fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() >= y.len() {
        x
    } else {
        y
    }
}

fn main() {
    // 栈 vs 堆
    let stack_val = 42;
    let heap_val = Box::new(42);
    println!("栈值 = {}，堆值解引用 = {}", stack_val, *heap_val);

    // Copy：赋值是复制，原变量仍可用
    let a = 1;
    let b = a;
    println!("Copy：a={} b={}（a 仍可用）", a, b);

    // Move：String 不是 Copy，赋值后原变量失效
    let s1 = String::from("hello");
    let s2 = s1;
    println!("Move：s2 = {}", s2);

    // Clone：显式深拷贝
    let s3 = s2.clone();
    println!("Clone：s2={} s3={}（两份独立数据）", s2, s3);

    // 函数传参也是移动
    let s4 = String::from("传入后所有权转移");
    let s5 = take_ownership(s4);
    println!("所有权转移：{}", s5);

    // 借用：不转移所有权
    let s6 = String::from("borrow");
    let len = s6.len();
    println!("借用读取：{} 长度 {}", s6, len);

    // Rc：单线程多所有者，引用计数
    let shared = Rc::new(vec![1, 2, 3]);
    let clone1 = Rc::clone(&shared);
    println!(
        "Rc 引用计数 = {}，长度 = {}",
        Rc::strong_count(&shared),
        clone1.len()
    );
    drop(clone1);
    println!("drop 一个克隆后引用计数 = {}", Rc::strong_count(&shared));

    // Rc<RefCell<T>>：共享 + 内部可变
    let cell = Rc::new(RefCell::new(0));
    let c2 = Rc::clone(&cell);
    *c2.borrow_mut() += 1;
    println!("RefCell 内部可变性：{}", *cell.borrow());

    // Drop 时机：离开作用域立即析构，没有 GC 停顿
    {
        let _d = Droppable("离开作用域");
        println!("块内，Drop 尚未发生");
    }
    println!("块结束后 Drop 已发生");

    let (p, q) = ("long string", "short");
    println!("生命周期标注 longest = {}", longest(p, q));
}
