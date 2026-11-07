fn main() {
    let s: String = String::from("héllo 世界");
    println!("字符串 = {}", s);

    // len() 是字节数，chars().count() 才是字符数
    println!("len() 字节数 = {}", s.len());
    println!("chars().count() 字符数 = {}", s.chars().count());

    // UTF-8 变长编码：每个字符占 1~4 字节
    for (i, ch) in s.chars().enumerate() {
        println!("  第 {} 个字符 {:?} 占 {} 字节", i, ch, ch.len_utf8());
    }

    // 不能 s[0]：索引可能切在多字节字符中间
    println!("get(0..1)  = {:?}（'h' 合法）", s.get(0..1));
    println!("get(0..2)  = {:?}（切在 'é' 中间，返回 None）", s.get(0..2));
    println!("get(1..3)  = {:?}（'é' 的完整 2 字节）", s.get(1..3));

    // 拼接：String 可变，push / push_str / format!
    let mut owned = String::from("hello");
    owned.push_str(", ");
    owned.push('R');
    owned.push_str("ust");
    println!("push 拼接 = {}", owned);
    println!("format! = {}", format!("{}-{}", owned, 42));

    // 原始字符串：不处理转义
    println!("原始字符串 = {}", r"C:\temp\new");
    println!("r# 包裹可含双引号 = {}", r#"里面有 "双引号""#);

    // 查找 / 切分 / 大小写
    let csv = "go,node,python,java,rust";
    println!("split 第 2 项 = {:?}", csv.split(',').nth(1));
    println!("contains(\"java\") = {}", csv.contains("java"));
    println!("to_uppercase = {}", "rust".to_uppercase());
    println!("trim = {:?}", "  padded  ".trim());

    // 解析
    let n: i32 = "42".parse().unwrap();
    println!("parse 42 + 1 = {}", n + 1);
}
