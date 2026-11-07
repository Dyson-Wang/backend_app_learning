// —— 数组：有序、可变长；map/filter/slice 返回新数组，splice 原地改 ——
const nums = [5, 2, 9, 1];
nums.push(7);
console.log('push 后原数组：', nums.join(','));
console.log('map 平方：', nums.map((x) => x * x).join(','));
console.log('filter 偶数：', nums.filter((x) => x % 2 === 0).join(','));
console.log('reduce 求和：', nums.reduce((s, x) => s + x, 0));

// sort 默认按字符串字典序，数字排序必须给比较函数
console.log('错误排序（字典序）：', [10, 2, 1].sort().join(','));
console.log('正确数字排序：', [10, 2, 1].sort((a, b) => a - b).join(','));

// slice 浅拷贝不改动原数组；splice 会原地增删
const copy = nums.slice(1, 3);
console.log('slice(1,3)：', copy.join(','), '，原数组长度仍是：', nums.length);

// —— Set 去重 / Map 任意键类型的有序字典 ——
const set = new Set([1, 2, 2, 3]);
set.add(4).delete(1);
console.log('Set 去重且保持插入顺序：', [...set].join(','), '，has(3)：', set.has(3));

const map = new Map();
map.set('name', 'alice').set(1, '数字键').set({}, '对象键');
console.log('Map 键的类型：', [...map.keys()].map((k) => typeof k).join(','));
console.log('Map.get(name)：', map.get('name'), '，size：', map.size);

// WeakMap 的键是弱引用，不阻止 GC，也不可遍历（21 课细讲）
const wm = new WeakMap();
const key = {};
wm.set(key, '私有数据');
console.log('WeakMap.get：', wm.get(key), '，不可遍历：', typeof wm[Symbol.iterator]);

// —— 迭代器协议：实现 [Symbol.iterator] 就能被 for...of / 展开 ——
const range = {
  from: 1,
  to: 4,
  [Symbol.iterator]() {
    let cur = this.from;
    const last = this.to;
    return {
      next() {
        return cur <= last
          ? { value: cur++, done: false }
          : { value: undefined, done: true };
      },
      [Symbol.iterator]() { return this; }, // 迭代器本身也是可迭代对象
    };
  },
};
console.log('自定义可迭代对象展开：', [...range].join(','));
console.log('用 Array.from 消费：', Array.from(range).join('-'));

// 普通对象默认不可迭代，需 Object.keys / values / entries
const obj = { a: 1, b: 2 };
console.log('对象无迭代器：', typeof obj[Symbol.iterator]);
console.log('Object.entries：', JSON.stringify(Object.entries(obj)));
