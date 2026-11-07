import process from 'node:process';

// —— 变量声明：let / const / var 的作用域与可变性 ——
let mutable = 1;
mutable = 2;
const fixed = 10;
console.log('let 可重新赋值：', mutable, '，const 只读绑定：', fixed);

// const 只锁「绑定」，不锁「对象内容」——JS 最常见的误解之一
const user = { name: 'alice', tags: ['dev'] };
user.name = 'bob';
user.tags.push('ops');
console.log('const 对象内部仍可变：', user.name, user.tags.join('/'));

// var 是函数/模块作用域且会提升；let/const 是块作用域
if (true) {
  var leaked = 'var 泄漏到块外';
  let scoped = 'let 只在块内';
}
console.log('块外读 var：', typeof leaked, '，块外读 let：', typeof scoped);

// —— 七种原始类型与 typeof ——
const samples = [42, 3.14, 'hi', true, undefined, null, 10n, Symbol('id')];
for (const value of samples) {
  console.log('值', String(value), '的 typeof =', typeof value);
}

// JS 没有 Go 那种统一零值：未赋值是 undefined，不是 0 / "" / false
let uninitialized;
console.log('未赋值的「零值」：', uninitialized, '，typeof：', typeof uninitialized);
console.log('显式转换默认零值：', Number(), JSON.stringify(String()), Boolean());

// 数字边界：超出安全整数后精度丢失，用 bigint 解决
console.log('最大安全整数：', Number.MAX_SAFE_INTEGER);
console.log('精度丢失 2^53 === 2^53+1：', 9007199254740992 === 9007199254740993);
console.log('BigInt 精确结果：', (9007199254740992n + 1n).toString());

// NaN 与判等：全等 === 不做隐式转换
console.log('NaN === NaN：', NaN === NaN, '，Number.isNaN：', Number.isNaN(NaN));
console.log('1 === "1"：', 1 === '1', '，1 == "1"（隐式转换）：', 1 == '1');
console.log('运行时：', process.version);
