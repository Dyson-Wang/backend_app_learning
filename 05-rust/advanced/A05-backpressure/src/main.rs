use std::sync::mpsc::{sync_channel, TrySendError};
use std::sync::{Arc, Barrier};
use std::thread;

fn main() {
    const CAPACITY: usize = 5;
    const TOTAL: usize = 20;

    let (tx, rx) = sync_channel::<usize>(CAPACITY);
    let gate = Arc::new(Barrier::new(2));
    let gate_worker = Arc::clone(&gate);

    let worker = thread::spawn(move || {
        gate_worker.wait();
        let mut processed = 0;
        for _ in rx {
            processed += 1;
        }
        processed
    });

    let mut accepted = 0;
    let mut rejected = 0;
    for i in 0..TOTAL {
        match tx.try_send(i) {
            Ok(()) => accepted += 1,
            Err(TrySendError::Full(_)) => rejected += 1, // 满则快速拒绝
            Err(TrySendError::Disconnected(_)) => break,
        }
    }

    gate.wait();
    drop(tx);
    let processed = worker.join().unwrap();

    println!("容量={}，提交={}", CAPACITY, TOTAL);
    println!("accepted={} rejected={} processed={}", accepted, rejected, processed);
    println!("结论：拒绝发生在入口，内存和延迟都是常数；无界队列只是把 OOM 推迟");

    let unbounded: Vec<usize> = (0..TOTAL).collect();
    println!("无界队列长度={}（全部接收，没有上限）", unbounded.len());
}
