fn main() {
    let clients = 100;
    let downstream_cap = 100;

    println!("== 1. 无退避：失败立即重试 3 次 ==");
    let no_backoff = clients * (1 + 3);
    println!(
        "下游收到 {} QPS（容量 {}，超载 {} 倍）",
        no_backoff,
        downstream_cap,
        no_backoff / downstream_cap
    );

    println!("== 2. 指数退避 + 全抖动 ==");
    // 3 次重试分散到 1s / 2s / 4s，每个窗口只承接 1/3 的客户端
    let peak = clients + clients / 3;
    println!("下游峰值 ≈ {} QPS（接近容量，不再放大）", peak);

    println!("== 3. 再加熔断 ==");
    println!("下游请求 = 0 QPS（open 状态本地快速失败，直接返回降级结果）");

    println!("== 4. 结论 ==");
    println!("重试必须有上限、退避和抖动；三个都没有，就是给故障加杠杆");
    println!("熔断才是止损的关键：下游已经挂了，再发请求只会让它更晚恢复");
}
