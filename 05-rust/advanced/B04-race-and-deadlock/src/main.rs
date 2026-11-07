use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::{mpsc, Arc, Mutex};
use std::thread;
use std::time::Duration;

fn main() {
    println!("== 1. 死锁：两个线程反序加锁 ==");
    let a = Arc::new(Mutex::new(()));
    let b = Arc::new(Mutex::new(()));
    let (tx, rx) = mpsc::channel::<&'static str>();

    let (a1, b1, tx1) = (Arc::clone(&a), Arc::clone(&b), tx.clone());
    thread::spawn(move || {
        let _ga = a1.lock().unwrap();
        thread::sleep(Duration::from_millis(10));
        let _gb = b1.lock().unwrap();
        let _ = tx1.send("thread-1 完成");
    });

    let (a2, b2, tx2) = (Arc::clone(&a), Arc::clone(&b), tx.clone());
    thread::spawn(move || {
        let _gb = b2.lock().unwrap();
        thread::sleep(Duration::from_millis(10));
        let _ga = a2.lock().unwrap();
        let _ = tx2.send("thread-2 完成");
    });

    match rx.recv_timeout(Duration::from_millis(100)) {
        Ok(msg) => println!("未死锁：{}", msg),
        Err(_) => println!("检测到死锁：100ms 内无人完成（各持一把锁等对方）"),
    }

    println!("== 2. 原子递增：4 线程 × 1000 次 ==");
    let counter = Arc::new(AtomicUsize::new(0));
    let mut handles = Vec::new();
    for _ in 0..4 {
        let c = Arc::clone(&counter);
        handles.push(thread::spawn(move || {
            for _ in 0..1000 {
                c.fetch_add(1, Ordering::SeqCst);
            }
        }));
    }
    for h in handles {
        h.join().unwrap();
    }
    println!("期望 4000，实际 {}", counter.load(Ordering::SeqCst));

    println!("== 3. 统一加锁顺序：两把锁按地址排序 ==");
    let x = Arc::new(Mutex::new(0));
    let y = Arc::new(Mutex::new(0));
    let mut handles = Vec::new();
    for _ in 0..2 {
        let x = Arc::clone(&x);
        let y = Arc::clone(&y);
        handles.push(thread::spawn(move || {
            let _gx = x.lock().unwrap();
            thread::sleep(Duration::from_millis(2));
            let _gy = y.lock().unwrap();
        }));
    }
    for h in handles {
        h.join().unwrap();
    }
    println!("两个线程按同一顺序加锁，未死锁");

    println!("结论：Rust 在编译期禁止数据竞争，但死锁仍然可能；靠统一加锁顺序或 try_lock 避免");
}
