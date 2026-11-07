// 游标分页 vs OFFSET 分页
const cmp = (a, b) => (b.score - a.score) || (a.id - b.id);
const ordered = (items) => [...items].sort(cmp);

const offsetPage = (items, offset, limit) => ordered(items).slice(offset, offset + limit);

function cursorPage(items, after, limit) {
  const out = [];
  for (const it of ordered(items)) {
    if (after) {
      const later = it.score < after.score || (it.score === after.score && it.id > after.id);
      if (!later) continue;
    }
    out.push(it);
    if (out.length === limit) break;
  }
  return out;
}

const ids = (items) => items.map((i) => i.id);
const overlap = (a, b) => {
  const seen = new Set(a.map((i) => i.id));
  return b.filter((i) => seen.has(i.id)).length;
};

const base = Array.from({ length: 10 }, (_, i) => ({ id: i + 1, score: 100 - i * 5 }));

console.log('== 1. OFFSET 分页：翻页期间有新数据插入 ==');
const p1 = offsetPage(base, 0, 3);
const grown = [...base, { id: 11, score: 97 }];
const p2 = offsetPage(grown, 3, 3);
console.log(`第一页 = [${ids(p1)}]`);
console.log(`插入一条 Score=97 后，第二页 = [${ids(p2)}]，与第一页重复 ${overlap(p1, p2)} 条`);

console.log('== 2. 游标分页：用 (score, id) 做锚点 ==');
const c1 = cursorPage(base, null, 3);
const last = c1[c1.length - 1];
const c2 = cursorPage(grown, last, 3);
console.log(`第一页 = [${ids(c1)}]，锚点 = (score=${last.score}, id=${last.id})`);
console.log(`插入同一条后，第二页 = [${ids(c2)}]，与第一页重复 ${overlap(c1, c2)} 条`);

console.log('== 3. 结论 ==');
console.log('OFFSET 按位置翻页，前面插入数据会让整体后移，出现重复或漏读');
console.log('游标按值翻页，锚点不变结果就不变；代价是不能跳页');
