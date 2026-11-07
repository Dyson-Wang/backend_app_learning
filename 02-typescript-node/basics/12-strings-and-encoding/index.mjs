import { TextEncoder, TextDecoder } from 'node:util';

const encoder = new TextEncoder();
const s = 'Hi 后端 😀';

// —— 不可变性：所有「修改」方法都返回新串，原串不变 ——
const upper = s.toUpperCase();
console.log('原串：', s, '；toUpperCase 新串：', upper, '；两者相等吗：', s === upper);
let assignErr = '无';
try { s[0] = 'X'; } catch (err) { assignErr = err.constructor.name; }
console.log('尝试改 s[0] 的结果：', assignErr, '（ESM 严格模式直接抛错，原串首字符仍是', s[0], '）');

// —— 码元 vs 码点 vs 字节：三个完全不同的「长度」——
const emoji = '😀';
console.log('😀 的 .length（UTF-16 码元）：', emoji.length);
console.log('[...emoji] 长度（Unicode 码点）：', [...emoji].length);
console.log('Buffer.byteLength（UTF-8 字节）：', Buffer.byteLength(emoji));
console.log('codePointAt(0) =', emoji.codePointAt(0).toString(16), '；String.fromCodePoint 还原：', String.fromCodePoint(0x1f600));
console.log('TextEncoder 编码 UTF-8 字节：', [...encoder.encode(emoji)].join(','));
console.log('TextDecoder 解码回：', new TextDecoder().decode(encoder.encode('中文')));

// —— 遍历差异：下标按码元，for...of 按码点 ——
const units = [];
for (let i = 0; i < s.length; i++) units.push(s[i]);
console.log('按下标取到', units.length, '个码元，含孤立高代理项：', units.includes('\ud83d'));
console.log('for...of 按码点取到', [...s].length, '个码点；末个码点：', [...s].at(-1));

// —— 拼接与格式化：+ / 模板字符串 / padStart / repeat ——
console.log('模板字符串：', `用户=${'alice'}，序号=${String(7).padStart(3, '0')}`);
console.log('repeat + concat：', 'ab'.repeat(3), '-'.concat('x', 'y'));
console.log('padEnd：', '标题'.padEnd(6, '·'));

// —— 转义与原始字符串：String.raw 保留反斜杠 ——
console.log('普通转义：', JSON.stringify('第一行\n第二行\t制表'));
console.log('Unicode 转义：', '\u4e2d', '\u{1f600}');
console.log('String.raw 原始：', String.raw`路径 C:\new\tab 里的 \n 不换行`);
console.log('JSON.stringify 转义：', JSON.stringify('引号" 换行\n'));

// —— 大小写 / 查找 / 切分 ——
console.log('大小写：', 'Äpfel'.toLocaleLowerCase('de'), 'ß'.toUpperCase());
console.log('查找：', s.includes('后端'), s.startsWith('Hi'), s.indexOf('端'));
console.log('切分：', JSON.stringify('a,b,,c'.split(',')));
console.log('替换：', 'a-b-a'.replace('a', 'X'), '，replaceAll：', 'a-b-a'.replaceAll('a', 'X'));
console.log('排序（码元序）：', ['b', 'A', 'a'].sort().join(','));

// —— 编码解码：Buffer 的 base64 / hex 与百分号编码 ——
const buf = Buffer.from('后端', 'utf8');
console.log('Buffer hex：', buf.toString('hex'), '；base64：', buf.toString('base64'));
console.log('base64 还原：', Buffer.from(buf.toString('base64'), 'base64').toString('utf8'));
console.log('encodeURIComponent：', encodeURIComponent('a b/中'));
console.log('decodeURIComponent：', decodeURIComponent('a%20b%2F%E4%B8%AD'));
console.log('trim：', JSON.stringify('  x  '.trim()));
