// —— 分支：if / else if / else ——
function grade(score) {
  if (score >= 90) return 'A';
  else if (score >= 60) return 'B';
  else return 'C';
}
console.log('分支评级：', [95, 70, 40].map(grade).join(' '));

// 三元表达式：适合简单二选一，嵌套多了就该换 if
const n = 7;
console.log('三元判奇偶：', n % 2 === 0 ? '偶数' : '奇数');

// switch 用严格匹配 ===，每个 case 记得 break / return，否则 fall-through 贯穿
function httpLabel(code) {
  switch (code) {
    case 200: return 'OK';
    case 404: return 'Not Found';
    case 500: return 'Server Error';
    default: return '未知';
  }
}
console.log('switch 状态码：', [200, 404, 500, 418].map(httpLabel).join(' / '));

// —— 经典 for / while / do...while ——
let sum = 0;
for (let i = 1; i <= 5; i++) sum += i;
console.log('for 求和 1..5：', sum);

let countdown = 3;
while (countdown > 0) countdown--;
console.log('while 结束后的值：', countdown);

let once = 0;
do { once++; } while (once < 1);
console.log('do...while 至少执行一次：', once);

// —— break / continue / 标签 break ——
const kept = [];
for (let i = 0; i < 10; i++) {
  if (i % 2 === 0) continue;   // 跳过偶数
  if (i > 7) break;            // 提前结束
  kept.push(i);
}
console.log('continue + break 结果：', kept.join(','));

let found = '未命中';
outer: for (let i = 0; i < 3; i++) {
  for (let j = 0; j < 3; j++) {
    if (i === 1 && j === 1) { found = `命中 (${i},${j})`; break outer; }
  }
}
console.log('标签 break 跳出双层循环：', found);

// —— 迭代：for...of（值，靠迭代器协议）vs for...in（键）——
const arr = ['a', 'b', 'c'];
const keys = [];
for (const idx in arr) keys.push(idx);   // 遍历下标/键
const vals = [];
for (const v of arr) vals.push(v);       // 遍历值
console.log('for...in 拿下标：', keys.join(','), '，for...of 拿值：', vals.join(','));

// for...of 可直接遍历字符串、Map、Set，并配合解构
const m = new Map([['x', 1], ['y', 2]]);
const pairs = [];
for (const [k, v] of m) pairs.push(`${k}=${v}`);
console.log('for...of 遍历 Map 并解构：', pairs.join(' '));
console.log('for...of 遍历字符串：', [...'后端'].join('-'));
