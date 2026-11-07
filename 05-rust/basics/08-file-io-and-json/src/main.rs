// 08 · 文件 I/O 与 JSON
// std::fs 读写临时文件；零依赖手写极简 JSON 序列化 / 取字段，说明为何生产用 serde。

use std::fs;
use std::io::{self, BufRead, BufReader, Write};
use std::path::PathBuf;

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

// 手写序列化：真实项目请用 serde_json
fn to_json(name: &str, age: u32) -> String {
    format!("{{\"name\":\"{}\",\"age\":{}}}", json_escape(name), age)
}

// 手写取字段：从 {"key":"value"} 里抠出字符串值
fn get_string(json: &str, key: &str) -> Option<String> {
    let pat = format!("\"{}\":\"", key);
    let rest = json.split_once(&pat)?.1;
    let end = rest.find('"')?;
    Some(rest[..end].to_string())
}

fn run() -> io::Result<()> {
    let path: PathBuf = std::env::temp_dir().join("brm-08-demo.json");
    let payload = to_json("爱丽丝", 30);
    println!("序列化结果: {}", payload);

    // 写：File::create 会截断或新建
    let mut file = fs::File::create(&path)?;
    file.write_all(payload.as_bytes())?;
    println!("已写入: {}", path.display());

    // 读回：read_to_string 一次性读入
    let content = fs::read_to_string(&path)?;
    println!("读回内容: {}", content);

    // 极简解析
    match get_string(&content, "name") {
        Some(name) => println!("解析 name = {}", name),
        None => println!("解析 name 失败"),
    }

    // 逐行读：BufReader + lines()
    let file = fs::File::open(&path)?;
    let reader = BufReader::new(file);
    let mut lines = 0;
    for line in reader.lines() {
        let line = line?;
        lines += 1;
        println!("第 {} 行: {}", lines, line);
    }

    // 元数据 + 清理
    println!("文件字节数: {}", fs::metadata(&path)?.len());
    fs::remove_file(&path)?;
    println!("临时文件已清理");
    Ok(())
}

fn main() {
    if let Err(e) = run() {
        eprintln!("IO 错误: {}", e);
    }
}
