// 09 · 并发原语入门
// thread::spawn + join、move 闭包、Arc<Mutex<T>> 共享计数、mpsc 通道传消息。

use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    // 1) 最简单的子线程：move 把所有权移交线程，join 等待并取回返回值
    let handle = thread::spawn(move || {
        let sum: i32 = (1..=5).sum();
        println!("[子线程] 1..=5 求和 = {}", sum);
        sum
    });
    match handle.join() {
        Ok(v) => println!("[主线程] 子线程返回 {}", v),
        Err(_) => println!("[主线程] 子线程 panic 了"),
    }

    // 2) Arc<Mutex<T>>：多线程安全地共享并修改同一个计数器
    let counter = Arc::new(Mutex::new(0i32));
    let mut handles = Vec::new();
    for _ in 0..4 {
        let c = Arc::clone(&counter);
        let h = thread::spawn(move || {
            // 锁作用域手动收窄：改完立即释放，避免持锁做无关操作
            {
                let mut num = c.lock().unwrap();
                *num += 1;
            }
        });
        handles.push(h);
    }
    for h in handles {
        h.join().unwrap();
    }
    let final_count = *counter.lock().unwrap();
    println!("Arc<Mutex> 计数结果 = {}（预期 4）", final_count);

    // 3) mpsc 通道：两个生产者发送，主线程收集后排序输出（保证可复现）
    let (tx, rx) = std::sync::mpsc::channel();
    let tx2 = tx.clone();
    let h1 = thread::spawn(move || {
        tx.send("来自线程A").unwrap();
    });
    let h2 = thread::spawn(move || {
        tx2.send("来自线程B").unwrap();
    });
    drop(tx); // 丢掉主发送者，rx 才能在所有发送者消失后结束
    h1.join().unwrap();
    h2.join().unwrap();
    let mut messages: Vec<&str> = rx.iter().collect();
    messages.sort();
    for msg in messages {
        println!("收到消息: {}", msg);
    }
    println!("通道已关闭，主线程退出");
}
