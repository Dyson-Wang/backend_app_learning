use std::sync::{mpsc, Arc, Mutex};
use std::thread;
use std::time::Duration;

fn main() {
    // 1. 基本 channel：一个生产者 -> 一个消费者
    let (tx, rx) = mpsc::channel();
    thread::spawn(move || {
        for i in 1..=3 {
            tx.send(i).unwrap();
        }
    });
    let got: Vec<i32> = rx.iter().collect();
    println!("单生产者 channel 收到 = {:?}", got);

    // 2. 多生产者：clone sender
    let (tx, rx) = mpsc::channel();
    let mut handles = vec![];
    for id in 0..3 {
        let tx = tx.clone();
        handles.push(thread::spawn(move || {
            tx.send(format!("worker-{}", id)).unwrap();
        }));
    }
    drop(tx); // 必须丢掉原始 sender，否则 rx.iter() 永远不结束
    let mut msgs: Vec<String> = rx.iter().collect();
    msgs.sort();
    println!("多生产者消息 = {:?}", msgs);
    for h in handles {
        h.join().unwrap();
    }

    // 3. 有界队列：sync_channel 满时发送方阻塞，形成背压
    let (tx, rx) = mpsc::sync_channel(2);
    let producer = thread::spawn(move || {
        for i in 0..5 {
            tx.send(i).unwrap();
        }
    });
    let mut received = vec![];
    for v in rx {
        received.push(v);
    }
    producer.join().unwrap();
    println!("有界队列收到 = {:?}", received);

    // 4. 共享状态：Arc<Mutex<T>>
    let counter = Arc::new(Mutex::new(0));
    let mut hs = vec![];
    for _ in 0..4 {
        let c = Arc::clone(&counter);
        hs.push(thread::spawn(move || {
            for _ in 0..100 {
                *c.lock().unwrap() += 1;
            }
        }));
    }
    for h in hs {
        h.join().unwrap();
    }
    println!("Arc<Mutex> 计数 = {}", *counter.lock().unwrap());

    // 5. 关闭信号传播
    let (stop_tx, stop_rx) = mpsc::channel::<()>();
    let worker = thread::spawn(move || {
        let mut stopped = false;
        for _ in 0..10 {
            if stop_rx.try_recv().is_ok() {
                stopped = true;
                break;
            }
            thread::sleep(Duration::from_millis(1));
        }
        stopped
    });
    thread::sleep(Duration::from_millis(2));
    stop_tx.send(()).unwrap();
    println!("worker 是否响应关闭信号 = {}", worker.join().unwrap());
}
