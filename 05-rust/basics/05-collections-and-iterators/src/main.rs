// 05 · 集合与迭代器
// Vec / HashMap / BTreeMap / HashSet 与惰性迭代器 map / filter / collect。

use std::collections::{BTreeMap, HashMap, HashSet};

fn main() {
    // 1) Vec：动态数组，push / sort / 索引
    let mut nums = vec![3, 1, 2];
    nums.push(5);
    nums.sort();
    println!("Vec 排序后: {:?}, 首元素 {}", nums, nums[0]);

    // 2) HashMap：无序键值映射，entry 处理「存在则更新」
    let mut scores: HashMap<&str, i32> = HashMap::new();
    scores.insert("alice", 90);
    scores.insert("bob", 85);
    *scores.entry("alice").or_insert(0) += 1;
    println!("HashMap alice = {:?}", scores.get("alice"));

    // 3) BTreeMap：键有序，输出稳定
    let mut ordered: BTreeMap<i32, &str> = BTreeMap::new();
    ordered.insert(2, "二");
    ordered.insert(1, "一");
    println!("BTreeMap 按键序: {:?}", ordered);

    // 4) HashSet 去重；为让输出稳定，转成有序 Vec 再打印
    let set: HashSet<i32> = [1, 2, 2, 3].into_iter().collect();
    let mut unique: Vec<i32> = set.into_iter().collect();
    unique.sort();
    println!("HashSet 去重: {:?}, 个数 {}", unique, unique.len());

    // 5) 迭代器惰性：map / filter 直到 collect 才真正执行
    let doubled: Vec<i32> = nums.iter().map(|x| *x * 2).filter(|x| *x > 4).collect();
    println!("惰性 map+filter 结果: {:?}", doubled);

    // 6) enumerate / sum 消费迭代器
    let total: i32 = nums.iter().sum();
    for (idx, val) in nums.iter().enumerate() {
        println!("  nums[{}] = {}", idx, val);
    }
    println!("求和 sum = {}", total);
}
