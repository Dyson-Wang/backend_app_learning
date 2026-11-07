use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::thread;

// 用 enum 表达路由结果，match 穷尽处理
enum Route {
    Health,
    Users,
    NotFound,
}

fn parse_route(path: &str) -> Route {
    match path {
        "/healthz" => Route::Health,
        "/users" => Route::Users,
        _ => Route::NotFound,
    }
}

fn handle(mut stream: TcpStream) {
    let mut buf = [0u8; 1024];
    let n = stream.read(&mut buf).unwrap_or(0);
    let req = String::from_utf8_lossy(&buf[..n]);

    let path = req
        .lines()
        .next()
        .and_then(|line| line.split_whitespace().nth(1))
        .unwrap_or("/");

    let (status, body) = match parse_route(path) {
        Route::Health => ("200 OK", r#"{"status":"ok"}"#.to_string()),
        Route::Users => (
            "200 OK",
            r#"[{"id":1,"name":"ada"},{"id":2,"name":"lin"}]"#.to_string(),
        ),
        Route::NotFound => ("404 Not Found", r#"{"error":"not found"}"#.to_string()),
    };

    let resp = format!(
        "HTTP/1.1 {}\r\nContent-Type: application/json\r\nContent-Length: {}\r\nConnection: close\r\n\r\n{}",
        status,
        body.len(),
        body
    );
    let _ = stream.write_all(resp.as_bytes());
    let _ = stream.flush();
}

fn main() {
    // 端口 0 = 让内核分配空闲端口，避免和本机服务冲突
    let listener = TcpListener::bind("127.0.0.1:0").expect("bind");
    let addr = listener.local_addr().unwrap();
    println!("内存服务监听 {}", addr);

    let server = thread::spawn(move || {
        for _ in 0..3 {
            if let Ok((stream, _)) = listener.accept() {
                handle(stream);
            }
        }
    });

    for path in ["/healthz", "/users", "/nope"] {
        let mut s = TcpStream::connect(addr).expect("connect");
        let req = format!(
            "GET {} HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n",
            path
        );
        s.write_all(req.as_bytes()).unwrap();

        let mut resp = String::new();
        s.read_to_string(&mut resp).unwrap();

        let status = resp.lines().next().unwrap_or("");
        let body = resp.split("\r\n\r\n").nth(1).unwrap_or("");
        println!("{} -> {} {}", path, status, body);
    }

    server.join().unwrap();
    println!("服务已关闭，进程退出");
}
