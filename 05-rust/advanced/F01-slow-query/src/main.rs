// 简化模型：所有请求同时到达，按到达顺序 FIFO 抢池子里的连接
fn simulate(durations: &[i32], pool: usize, wait_timeout: i32) -> (i32, i32, i32) {
    let mut free_at = vec![0i32; pool];
    let (mut ok, mut timeout, mut total) = (0, 0, 0);

    for &d in durations {
        let mut best = 0usize;
        for i in 1..pool {
            if free_at[i] < free_at[best] {
                best = i;
            }
        }
        let wait = free_at[best];
        if wait > wait_timeout {
            timeout += 1;
            continue;
        }
        free_at[best] = wait + d;
        if free_at[best] > total {
            total = free_at[best];
        }
        ok += 1;
    }
    (ok, timeout, total)
}

fn main() {
    let normal = vec![30; 10];

    println!("== 1. 正常：全部 30ms ==");
    let (ok, to, total) = simulate(&normal, 2, 100);
    println!("成功 = {}，超时 = {}，总耗时 = {}ms", ok, to, total);

    println!("== 2. 混入 1 个 300ms 慢查询（等待上限 100ms） ==");
    let mut slow = normal.clone();
    slow[0] = 300;
    let (ok, to, total) = simulate(&slow, 2, 100);
    println!("成功 = {}，超时 = {}，总耗时 = {}ms", ok, to, total);

    println!("== 3. 结论 ==");
    println!("池子只有 2 个连接，1 个慢查询就占掉一半；等待超时把雪崩挡在入口");
    println!("慢查询必须配 statement_timeout + 索引，光加连接数只是把问题推后");
}
