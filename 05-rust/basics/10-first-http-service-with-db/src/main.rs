use std::collections::HashMap;
use std::io::{BufRead, BufReader, Read, Write};
use std::net::{TcpListener, TcpStream};
use std::sync::{Arc, Mutex};

// 零第三方依赖：手写 TCP + HTTP/1.1，只为看清「框架替你做了什么」
type Users = Arc<Mutex<HashMap<u64, String>>>;

fn json_escape(s: &str) -> String {
    let mut out = String::new();
    for c in s.chars() {
        match c {
            '"' => out.push_str("\\\""),
            '\\' => out.push_str("\\\\"),
            _ => out.push(c),
        }
    }
    out
}

fn extract_name(body: &str) -> String {
    let key = "\"name\"";
    let Some(pos) = body.find(key) else { return String::new() };
    let rest = &body[pos + key.len()..];
    let Some(colon) = rest.find(':') else { return String::new() };
    let rest = rest[colon + 1..].trim_start();
    let Some(rest) = rest.strip_prefix('"') else { return String::new() };
    let Some(end) = rest.find('"') else { return String::new() };
    rest[..end].trim().to_string()
}

fn user_json(id: u64, name: &str) -> String {
    format!("{{\"id\":{},\"name\":\"{}\"}}", id, json_escape(name))
}

fn err_json(code: &str, message: &str) -> String {
    format!("{{\"error\":{{\"code\":\"{}\",\"message\":\"{}\"}}}}", code, message)
}

fn route(method: &str, path: &str, body: &str, users: &Users) -> (u16, &'static str, String) {
    if path == "/healthz" {
        return (200, "OK", "{\"status\":\"ok\"}".to_string());
    }
    if path == "/readyz" {
        return (200, "OK", "{\"status\":\"ready\"}".to_string());
    }

    if path == "/users" {
        return match method {
            "GET" => {
                let map = users.lock().unwrap();
                let items: Vec<String> = map
                    .iter()
                    .map(|(id, name)| user_json(*id, name))
                    .collect();
                (200, "OK", format!("{{\"data\":[{}]}}", items.join(",")))
            }
            "POST" => {
                let name = extract_name(body);
                if name.is_empty() {
                    return (400, "Bad Request", err_json("VALIDATION", "name is required"));
                }
                let mut map = users.lock().unwrap();
                let id = map.keys().max().copied().unwrap_or(0) + 1;
                map.insert(id, name.clone());
                (201, "Created", format!("{{\"data\":{}}}", user_json(id, &name)))
            }
            _ => (405, "Method Not Allowed", err_json("METHOD_NOT_ALLOWED", "use GET or POST")),
        };
    }

    if let Some(id_str) = path.strip_prefix("/users/") {
        if method != "GET" {
            return (405, "Method Not Allowed", err_json("METHOD_NOT_ALLOWED", "use GET"));
        }
        return match id_str.parse::<u64>() {
            Ok(id) => {
                let map = users.lock().unwrap();
                match map.get(&id) {
                    Some(name) => (200, "OK", format!("{{\"data\":{}}}", user_json(id, name))),
                    None => (404, "Not Found", err_json("NOT_FOUND", "user not found")),
                }
            }
            Err(_) => (400, "Bad Request", err_json("BAD_ID", "id must be an integer")),
        };
    }

    (404, "Not Found", err_json("NOT_FOUND", "route not found"))
}

fn respond(stream: &mut TcpStream, status: u16, reason: &str, body: &str) -> std::io::Result<()> {
    let response = format!(
        "HTTP/1.1 {} {}\r\nContent-Type: application/json; charset=utf-8\r\nContent-Length: {}\r\nConnection: close\r\n\r\n{}",
        status,
        reason,
        body.as_bytes().len(), // 中文按字节数，不能用 body.len() 的语义混淆
        body
    );
    stream.write_all(response.as_bytes())?;
    stream.flush()
}

fn handle(mut stream: TcpStream, users: Users) -> std::io::Result<()> {
    let mut reader = BufReader::new(stream.try_clone()?);

    let mut request_line = String::new();
    reader.read_line(&mut request_line)?;
    if request_line.is_empty() {
        return Ok(());
    }
    let mut parts = request_line.split_whitespace();
    let method = parts.next().unwrap_or("").to_string();
    let path = parts.next().unwrap_or("/").to_string();

    let mut content_length = 0usize;
    loop {
        let mut line = String::new();
        reader.read_line(&mut line)?;
        if line == "\r\n" || line.is_empty() {
            break;
        }
        let lower = line.to_ascii_lowercase();
        if let Some(rest) = lower.strip_prefix("content-length:") {
            content_length = rest.trim().parse().unwrap_or(0);
        }
    }

    let mut body = vec![0u8; content_length];
    if content_length > 0 {
        reader.read_exact(&mut body)?;
    }
    let body = String::from_utf8_lossy(&body).to_string();

    let (status, reason, payload) = route(&method, &path, &body, &users);
    println!("{} {} {} ", method, path, status);
    respond(&mut stream, status, reason, &payload)
}

fn main() -> std::io::Result<()> {
    let port = std::env::var("PORT").unwrap_or_else(|_| "8084".to_string());
    let addr = format!("127.0.0.1:{}", port);
    let listener = TcpListener::bind(&addr)?;
    let users: Users = Arc::new(Mutex::new(HashMap::new()));

    println!("listening on http://{}", addr);
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                let users = Arc::clone(&users);
                std::thread::spawn(move || {
                    if let Err(e) = handle(stream, users) {
                        eprintln!("handle error: {}", e);
                    }
                });
            }
            Err(e) => eprintln!("accept error: {}", e),
        }
    }
    Ok(())
}
