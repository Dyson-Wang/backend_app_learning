// —— 自定义生成器：function* + yield，返回「既是迭代器又是可迭代对象」的对象 ——
function* countdown(from) {
  for (let i = from; i > 0; i--) yield i;
  return '发射';   // return 值只在 done:true 那次出现，for...of 拿不到
}
const gen = countdown(3);
console.log('生成器是迭代器：', typeof gen.next, '；也是可迭代对象：', typeof gen[Symbol.iterator]);
console.log('手动 next：', JSON.stringify([gen.next(), gen.next()]));
console.log('剩余用展开消费：', [...gen].join(','), '（已消费的不会回来）');

// —— 惰性管道：map / filter / take / chain 都是「返回生成器的高阶函数」——
function* map(iterable, fn) { for (const x of iterable) yield fn(x); }
function* filter(iterable, pred) { for (const x of iterable) if (pred(x)) yield x; }
function* chain(...iterables) { for (const it of iterables) yield* it; }
function* take(iterable, n) {
  if (n <= 0) return;
  let i = 0;
  for (const x of iterable) {
    yield x;
    if (++i >= n) return; // return 会触发 IteratorClose，关闭上游生成器
  }
}

// 定义管道时什么都不算，只有消费时才逐个计算（惰性求值）
let computed = 0;
const pipeline = take(
  filter(
    map([1, 2, 3, 4, 5, 6], (x) => { computed++; return x * x; }),
    (v) => v % 2 === 0,
  ),
  2,
);
console.log('构造管道后尚未计算，computed =', computed);
const result = [...pipeline];
console.log('消费后结果：', result.join(','), '；computed =', computed, '（只算到够用为止）');

// —— 无限序列 + take(N)：懒生成，靠 take 限定，绝不会真跑成死循环 ——
function* naturals() {
  let n = 0;
  try {
    while (true) yield n++;   // 惰性：只有 next() 才前进一步
  } finally {
    console.log('  [naturals] 被 take 关闭，finally 执行，没有死循环');
  }
}
console.log('无限序列前 5 项：', [...take(naturals(), 5)].join(','));
console.log('前 3 项的平方：', [...take(map(naturals(), (x) => x * x), 3)].join(','));

// —— chain / yield* 拼接多个可迭代对象 ——
console.log('chain 拼接数组/Set/字符串：', [...chain([1, 2], new Set([3, 3, 4]), 'ab')].join(','));

// —— 一次性消费陷阱：生成器消费完就空了，数组/Map 可反复迭代 ——
const once = countdown(2);
console.log('第一次消费：', [...once].join(','), '；第二次消费：', JSON.stringify([...once]));
const arr = [1, 2, 3];
console.log('数组可反复迭代：', [...arr].join(','), [...arr].join(','));
console.log('结论：生成器一次性；要复用就先 [...gen] 成数组，或每次重新调用生成器函数');
