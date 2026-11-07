use std::time::{Duration, Instant, SystemTime, UNIX_EPOCH};

// 标准库没有随机数：用固定种子的 LCG 保证可复现
fn lcg(seed: &mut u64) -> u64 {
    *seed = seed
        .wrapping_mul(6364136223846793005)
        .wrapping_add(1442695040888963407);
    *seed >> 33
}

fn main() {
    // Instant：单调时钟，用来测耗时
    let start = Instant::now();
    let mut acc = 0u64;
    for i in 0..100_000u64 {
        acc = acc.wrapping_add(i);
    }
    let elapsed = start.elapsed();
    println!("循环求和 = {}，耗时 {:?}（仅演示计时 API）", acc, elapsed);

    // Duration：时长
    let d = Duration::from_millis(1500);
    println!("1500ms = {}s，as_millis = {}", d.as_secs(), d.as_millis());

    // SystemTime：墙上时间，可早于或晚于 UNIX_EPOCH
    let since_epoch = SystemTime::now().duration_since(UNIX_EPOCH).unwrap();
    println!(
        "距 UNIX_EPOCH 的天数 = {}（时间点演示，不作结论）",
        since_epoch.as_secs() / 86400
    );

    // 标准库没有正则：手写前缀判断 + 字符截取，生产用 regex crate
    let text = "order-2026-11-08.log";
    let has_order = text.starts_with("order-");
    let date_part: String = text.chars().skip(6).take(10).collect();
    println!("starts_with = {}，截取日期 = {}", has_order, date_part);

    // 排序
    let mut nums = vec![5, 2, 9, 1, 7];
    nums.sort();
    println!("升序 = {:?}", nums);
    nums.sort_by(|a, b| b.cmp(a));
    println!("降序 = {:?}", nums);

    let mut words = vec!["pear", "apple", "fig"];
    words.sort_by_key(|w| w.len());
    println!("按长度排序 = {:?}", words);

    // 固定种子伪随机，保证每次输出一致
    let mut seed = 42u64;
    let rolls: Vec<u64> = (0..5).map(|_| lcg(&mut seed) % 6 + 1).collect();
    println!("固定种子伪随机骰子 = {:?}", rolls);

    println!(
        "max/min = {}/{}",
        nums.iter().max().unwrap(),
        nums.iter().min().unwrap()
    );
}
