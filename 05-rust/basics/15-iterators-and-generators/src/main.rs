// 自定义迭代器：只要实现 Iterator，就自动获得 map/filter/take/sum 等几十个方法
struct Counter {
    current: u32,
    max: u32,
}

impl Counter {
    fn new(max: u32) -> Self {
        Counter { current: 0, max }
    }
}

impl Iterator for Counter {
    type Item = u32;

    fn next(&mut self) -> Option<u32> {
        if self.current < self.max {
            self.current += 1;
            Some(self.current)
        } else {
            None
        }
    }
}

fn main() {
    let sum: u32 = Counter::new(5).sum();
    println!("自定义迭代器求和 1..=5 = {}", sum);

    let doubled: Vec<u32> = Counter::new(5).map(|x| x * 2).collect();
    println!("map 后 = {:?}", doubled);

    let evens: Vec<u32> = Counter::new(10).filter(|x| x % 2 == 0).collect();
    println!("filter 偶数 = {:?}", evens);

    let chained: Vec<u32> = Counter::new(3).chain(Counter::new(3)).collect();
    println!("chain 两个迭代器 = {:?}", chained);

    println!("take(2) = {:?}", Counter::new(100).take(2).collect::<Vec<_>>());

    // 惰性：不消费就不计算
    let lazy = Counter::new(1_000_000).map(|x| {
        println!("  计算 {}", x);
        x
    });
    println!("还没消费，上面不会打印任何「计算」");
    let first = lazy.take(1).collect::<Vec<_>>();
    println!("消费 1 个 = {:?}", first);

    // 无限序列 + take
    let powers = std::iter::successors(Some(1u64), |&n| Some(n * 2));
    println!("2 的幂前 6 项 = {:?}", powers.take(6).collect::<Vec<_>>());

    // 迭代器是一次性的：next 会推进内部状态
    let mut it = Counter::new(2);
    println!("next = {:?}", it.next());
    println!("next = {:?}", it.next());
    println!("next = {:?}", it.next());
}
