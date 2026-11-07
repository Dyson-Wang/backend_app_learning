// —— 算术运算符：/ 是浮点除，** 是幂（右结合）——
console.log('算术：7 / 2 =', 7 / 2, '，取余 7 % 3 =', 7 % 3, '，幂 2 ** 10 =', 2 ** 10);
console.log('幂右结合 2 ** 3 ** 2 =', 2 ** 3 ** 2, '（等价于 2 ** (3 ** 2)）');
console.log('BigInt 算术：', (2n ** 10n).toString(), '（任意精度，不溢出）');

// —— 比较 / 相等：=== 严格，== 会先隐式转换 ——
console.log('比较：', 3 > 2, 3 >= 3, 'a' < 'b');
console.log('== vs ===：', 0 == false, 0 === false, '' == 0, '' === 0);
console.log('Object.is 特例：', Object.is(NaN, NaN), Object.is(0, -0), 0 === -0);

// —— 逻辑与短路：用计数器证明右侧「没有被执行」——
let calls = 0;
const bump = () => { calls++; return true; };
const a = false && bump();   // 短路，bump 不执行
const b = true || bump();    // 短路，bump 不执行
console.log('短路后 calls =', calls, '；a =', a, '，b =', b);
console.log('?? 只在 null/undefined 时取右：', null ?? '默认', 0 ?? '默认', '' ?? '默认');
console.log('&& / || 返回操作数本身：', 0 || 5, 1 && 2, 'x' && 0 || 7);

// —— 位运算：操作数被转成 32 位有符号整数，会「回绕」——
console.log('位与 5 & 3 =', 5 & 3, '；位或 5 | 3 =', 5 | 3, '；异或 5 ^ 3 =', 5 ^ 3);
console.log('取反 ~5 =', ~5, '（即 -6）；左移 1 << 31 =', 1 << 31, '（溢出成负数）');
console.log('无符号右移 -1 >>> 0 =', -1 >>> 0, '；有符号右移 -8 >> 2 =', -8 >> 2);
console.log('32 位回绕实证：2147483648 | 0 =', 2147483648 | 0, '（普通加法仍是', 2147483648, '）');

// —— 整数「溢出」行为：number 是 float64，无溢出但有精度上限 ——
console.log('超安全整数静默丢精度：', Number.MAX_SAFE_INTEGER + 1 === Number.MAX_SAFE_INTEGER + 2);
console.log('超 MAX_VALUE 变 Infinity：', Number.MAX_VALUE * 2);
console.log('BigInt 无上限：', (2n ** 64n).toString().length, '位十进制数');

// —— 赋值与复合赋值 / 逻辑赋值 / 前后缀自增 ——
let x = 10;
x += 3; x *= 2; x **= 2;
console.log('复合赋值后 x =', x);
let lo = 0, hi = 5, unset;
lo ||= 3; hi &&= 9; unset ??= 'init';   // 逻辑赋值同样短路
console.log('逻辑赋值：', lo, hi, unset);
let i = 1;
console.log('前后缀自增：', i++, ++i, i);

// —— 优先级、三元、隐式类型转换 ——
console.log('优先级 2 + 3 * 4 =', 2 + 3 * 4, '；(2 + 3) * 4 =', (2 + 3) * 4);
console.log('优先级 && 高于 ||：', true || false && false, '（等价于 true || (false && false)）');
console.log('三元：', 7 % 2 === 0 ? '偶数' : '奇数', '；链式：', 85 >= 90 ? 'A' : 85 >= 60 ? 'B' : 'C');
console.log('?? 不能与 &&/|| 裸混用，加括号：', (null ?? 0) && 1);
console.log('+ 遇字符串变拼接：', '3' + 4, '；- 强制转数字：', '8' - 3, '，3 * 2：', '3' * '2');
console.log('一元 + 转数字：', +'42', '；Number("") =', Number(''), '，Number("abc") =', Number('abc'));
console.log('falsy 列表里为真的个数：', [0, '', null, undefined, NaN, false].filter(Boolean).length);
