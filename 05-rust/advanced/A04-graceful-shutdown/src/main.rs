use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::thread;
use std::time::{Duration, Instant};

fn handle(mut stream: TcpStream) {
    let mut buf = [0u8; 1024];
    let _ = stream.read(&mut buf);
    thread::sleep(Duration::from_millis(120));       // 模拟在途慢请求
    let body = "slow done";
    let resp = format!(
        "HTTP/1.1 200 OK\r\nContent-Length: {}\r\nConnection: close\r\n\r\n{}",
        body.len(),
        body
    );
    let _ = stream.write_all(resp.as_bytes());
}

fn main() {
    let listener = TcpListener::bind("127.0.0.1:0").expect("bind");
    let addr = listener.local_addr().expect("local_addr");
    listener.set_nonblocking(true).expect("nonblocking");

    let stop = Arc::new(AtomicBool::new(false));
    let stop_accept = Arc::clone(&stop);

    // accept 循环：非阻塞轮询，这样 stop 能及时生效
    let acceptor = thread::spawn(move || {
        loop {
            if stop_accept.load(Ordering::SeqCst) {
                break;
            }
            match listener.accept() {
                Ok((s, _)) => {
                    thread::spawn(move || handle(s));
                }
                Err(ref e) if e.kind() == std::io::ErrorKind::WouldBlock => {
                    thread::sleep(Duration::from_millis(5));
                }
                Err(_) => break,
            }
        }
    });

    println!("== 1. 发起一个 120ms 的在途请求 ==");
    let mut conn = TcpStream::connect(addr).expect("connect");
    conn.write_all(b"GET /slow HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n")
        .expect("write");
    thread::sleep(Duration::from_millis(30));

    println!("== 2. 触发优雅关闭 ==");
    let t0 = Instant::now();
    stop.store(true, Ordering::SeqCst);
    acceptor.join().expect("join acceptor");          // 停止 accept
    let cost = t0.elapsed().as_millis();

    println!("== 3. 在途请求的结果 ==");
    let mut resp = String::new();
    conn.read_to_string(&mut resp).expect("read");
    let body = resp.rsplit("\r\n\r\n").next().unwrap_or("").to_string();
    println!("body = {}，关闭耗时 = {}ms", body, cost);

    println!("== 4. 关闭后新请求 ==");
    match TcpStream::connect(addr) {
        Ok(_) => println!("新请求 -> 连接成功（不应该发生）"),
        Err(e) => println!("新请求 -> 被拒绝 {}", e.kind()),
    }

    println!("== 5. 结论 ==");
    println!("停止 accept 之后，在途 handler 仍能写完响应；新连接直接被拒");
    println!("生产上用 tokio 的 graceful shutdown + CancellationToken，别手写轮询");
}
