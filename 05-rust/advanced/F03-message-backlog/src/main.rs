fn main() {
    let backlog = 1000;
    let per_sec = 10;      // 单消费者每秒 10 条
    let producers = 20;    // 生产端每秒 20 条

    println!("== 1. 单消费者，逐条处理 ==");
    println!("{} 条 / {} 条每秒 = {} 秒", backlog, per_sec, backlog / per_sec);

    println!("== 2. 扩到 3 个消费者 ==");
    let three = (backlog + per_sec * 3 - 1) / (per_sec * 3);
    println!("{} 条 / {} 条每秒 = {} 秒", backlog, per_sec * 3, three);

    println!("== 3. 批量 10 条一次（单消费者） ==");
    let batch = (backlog + per_sec * 10 - 1) / (per_sec * 10);
    println!("{} 条 / {} 条每秒 = {} 秒", backlog, per_sec * 10, batch);

    println!("== 4. 生产仍在继续（每秒 20 条） ==");
    let net = per_sec * 10 - producers;
    let clear = (backlog + net - 1) / net;
    println!("净消化 {} 条每秒，清空 {} 条需要 {} 秒", net, backlog, clear);

    println!("== 5. 结论 ==");
    println!("积压先算「净消化速率」，不是「消费速率」——生产端还在往里灌");
    println!("扩消费者最快但受分区数限制；批量提升单机吞吐，代价是延迟和失败重放");
}
