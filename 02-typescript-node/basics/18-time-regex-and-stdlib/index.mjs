import { randomInt, randomUUID, randomBytes } from 'node:crypto';
import { format, inspect, types } from 'node:util';
import { performance } from 'node:perf_hooks';

// —— 时间点 / 时长 / 时区：全部用「固定时刻」，不打印当前绝对时间 ——
const instant = new Date(Date.UTC(2024, 0, 2, 3, 4, 5)); // 2024-01-02T03:04:05Z
console.log('固定时刻 ISO：', instant.toISOString());
console.log('UTC 完整格式：', new Intl.DateTimeFormat('zh-CN', { timeZone: 'UTC', dateStyle: 'full', timeStyle: 'long' }).format(instant));
console.log('上海时区小时：', new Intl.DateTimeFormat('zh-CN', { timeZone: 'Asia/Shanghai', hour: '2-digit', hour12: false }).format(instant));
console.log('纽约时区小时：', new Intl.DateTimeFormat('en-US', { timeZone: 'America/New_York', hour: '2-digit', hour12: false }).format(instant));
const later = new Date(Date.UTC(2024, 0, 2, 4, 4, 5));
console.log('时长（两固定时刻之差，毫秒）：', later - instant, '；折合分钟：', (later - instant) / 60000);
console.log('解析 ISO 字符串：', Date.parse('2024-01-02T03:04:05Z') === instant.getTime());

// —— 正则：exec / 命名捕获 / matchAll / replace / 标志 ——
const re = /(?<user>\w+)@(?<host>[\w.]+)/g;   // g 正则带状态：lastIndex
const text = '联系 alice@ex.com 或 bob@site.org';
const m = re.exec(text);
console.log('exec 命名捕获：', m.groups.user, m.groups.host, '；索引：', m.index);
console.log('exec 推进了 lastIndex：', re.lastIndex);
re.lastIndex = 0; // 复用同一个 g 正则前必须归零，否则 matchAll 会漏掉开头
console.log('matchAll 全部用户名：', [...text.matchAll(re)].map((x) => x.groups.user).join(','));
console.log('replace + 回调：', text.replace(re, (all, u, h) => `${u} [at] ${h}`));
console.log('test 校验：', /^\d{3}-\d{4}$/.test('123-4567'), '；不合法：', /^\d{3}-\d{4}$/.test('12-4567'));
console.log('点号匹配换行需 s 标志：', /a.b/.test('a\nb'), /a.b/s.test('a\nb'));

// —— 随机数：固定种子可复现序列 + crypto API 只演示不依赖具体值 ——
function seeded(seed) {                 // 线性同余，固定种子 → 结果可复现
  let s = seed >>> 0;
  return () => ((s = (1664525 * s + 1013904223) >>> 0) / 2 ** 32);
}
const rnd = seeded(42);
console.log('固定种子(42) 伪随机序列：', [rnd(), rnd(), rnd()].map((x) => x.toFixed(4)).join(','));
console.log('crypto.randomInt 类型：', typeof randomInt, '；randomUUID 长度：', randomUUID().length);
console.log('crypto.randomBytes(8) 字节数：', randomBytes(8).length, '（值为随机，不打印内容）');

// —— 排序：数值 / 字段 / 本地化 / 不改原数组的 toSorted ——
console.log('数值排序：', [10, 2, 1, 21].sort((a, b) => a - b).join(','));
const people = [{ n: 'bob', a: 30 }, { n: 'alice', a: 25 }];
console.log('按字段排序：', people.map((p) => p.n).sort((x, y) => x.localeCompare(y)).join(','));
console.log('Intl.Collator 中文排序：', ['张三', '李四', '王五'].sort(new Intl.Collator('zh-CN').compare).join(','));
const origin = [3, 1, 2];
console.log('toSorted 不改原数组：', origin.toSorted().join(','), '；原数组：', origin.join(','));

// —— 常用标准库工具函数 ——
console.log('util.format：', format('%s=%d', '端口', 8080));
console.log('util.inspect：', inspect({ a: [1, 2] }, { compact: true }));
console.log('util.types 判断：', types.isDate(instant), types.isPromise(Promise.resolve()));
console.log('structuredClone 深拷贝：', JSON.stringify(structuredClone({ x: [1, { y: 2 }] })));
console.log('Object.groupBy 分组：', JSON.stringify(Object.groupBy([1, 2, 3, 4], (x) => (x % 2 ? '奇' : '偶'))));
console.log('Intl.NumberFormat 千分位：', new Intl.NumberFormat('en-US').format(1234567));
console.log('performance.now 可用：', typeof performance.now === 'function');
