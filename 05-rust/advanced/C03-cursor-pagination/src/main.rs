#[derive(Clone, Debug, PartialEq)]
struct Item {
    id: i32,
    score: i32,
}

fn ordered(items: &[Item]) -> Vec<Item> {
    let mut out = items.to_vec();
    out.sort_by(|a, b| b.score.cmp(&a.score).then(a.id.cmp(&b.id)));
    out
}

fn offset_page(items: &[Item], offset: usize, limit: usize) -> Vec<Item> {
    let s = ordered(items);
    if offset >= s.len() {
        return Vec::new();
    }
    s[offset..(offset + limit).min(s.len())].to_vec()
}

fn cursor_page(items: &[Item], after: Option<&Item>, limit: usize) -> Vec<Item> {
    let mut out = Vec::new();
    for it in ordered(items) {
        if let Some(a) = after {
            let later = it.score < a.score || (it.score == a.score && it.id > a.id);
            if !later {
                continue;
            }
        }
        out.push(it);
        if out.len() == limit {
            break;
        }
    }
    out
}

fn ids(items: &[Item]) -> Vec<i32> {
    items.iter().map(|i| i.id).collect()
}

fn overlap(a: &[Item], b: &[Item]) -> usize {
    b.iter().filter(|i| a.iter().any(|x| x.id == i.id)).count()
}

fn main() {
    let base: Vec<Item> = (0..10)
        .map(|i| Item {
            id: i + 1,
            score: 100 - i * 5,
        })
        .collect();

    println!("== 1. OFFSET 分页：翻页期间有新数据插入 ==");
    let p1 = offset_page(&base, 0, 3);
    let mut grown = base.clone();
    grown.push(Item { id: 11, score: 97 });
    let p2 = offset_page(&grown, 3, 3);
    println!("第一页 = {:?}", ids(&p1));
    println!(
        "插入一条 Score=97 后，第二页 = {:?}，与第一页重复 {} 条",
        ids(&p2),
        overlap(&p1, &p2)
    );

    println!("== 2. 游标分页：用 (score, id) 做锚点 ==");
    let c1 = cursor_page(&base, None, 3);
    let last = c1.last().unwrap().clone();
    let c2 = cursor_page(&grown, Some(&last), 3);
    println!(
        "第一页 = {:?}，锚点 = (score={}, id={})",
        ids(&c1),
        last.score,
        last.id
    );
    println!(
        "插入同一条后，第二页 = {:?}，与第一页重复 {} 条",
        ids(&c2),
        overlap(&c1, &c2)
    );

    println!("== 3. 结论 ==");
    println!("OFFSET 按位置翻页，前面插入数据会让整体后移，出现重复或漏读");
    println!("游标按值翻页，锚点不变结果就不变；代价是不能跳页");
}
