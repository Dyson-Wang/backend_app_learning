// —— 闭包：内部函数「记住」并持续访问外层作用域的变量 ——
function makeCounter(start = 0) {
  let count = start;                 // 被闭包捕获，随返回的函数一起存活
  return () => ++count;
}
const c1 = makeCounter(10);
console.log('闭包计数器：', c1(), c1(), c1());
const c2 = makeCounter();
console.log('每次调用生成独立闭包：', c2(), c2(), '；c1 不受影响：', c1());

// —— 坑：循环里用 var，所有闭包共享同一个变量 ——
const bad = [];
for (var i = 0; i < 3; i++) bad.push(() => i);
console.log('var 捕获（全指向同一个 i）：', bad.map((f) => f()).join(','));
const good = [];
for (let j = 0; j < 3; j++) good.push(() => j); // let 每轮是新绑定
console.log('let 捕获（每轮独立）：', good.map((f) => f()).join(','));

// —— 捕获变量的生命周期：外层返回后变量仍存活，直到闭包被释放 ——
function makeAdder(n) {
  const base = n;
  return (x) => base + x;
}
const add5 = makeAdder(5);
console.log('闭包保活 base：', add5(3), add5(7), '；base 随 add5 一直存在');
console.log('函数也是值：typeof =', typeof add5, '；可作数组元素：', [add5].map((f) => f(1)).join(','));

// —— 高阶函数：接收函数 / 返回函数 ——
function compose(f, g) { return (x) => f(g(x)); }
const doubleThenSquare = compose((x) => x * x, (x) => x * 2);
console.log('compose 先翻倍再平方：', doubleThenSquare(3));
function once(fn) {
  let called = false, cached;
  return (...args) => {
    if (!called) { called = true; cached = fn(...args); }
    return cached;
  };
}
let hits = 0;
const init = once(() => ++hits);
console.log('once 只执行一次：', init(), init(), '；内部计数 hits =', hits);
function memoize(fn) {
  const cache = new Map();
  return (x) => (cache.has(x) ? cache.get(x) : (cache.set(x, fn(x)), cache.get(x)));
}
const slowSquare = memoize((x) => { console.log('  [memoize] 真正计算', x); return x * x; });
console.log('memoize：', slowSquare(4), slowSquare(4), '（第二次命中缓存）');

// —— 常见高阶函数：数组自带 map / filter / reduce ——
console.log('数组高阶函数：', [1, 2, 3, 4].filter((x) => x % 2).map((x) => x * 10).reduce((a, b) => a + b, 0));

// —— 回调与 this：普通函数作回调会丢 this，箭头/ bind 可救 ——
const obj = {
  factor: 3,
  scale(x) { return this.factor * x; },
  makeMultiplier() { return (x) => this.factor * x; }, // 箭头从外层方法捕获 this
};
function callIt(cb) { return cb(4); }
let lost;
try { lost = callIt(obj.scale); } catch (err) { lost = err.constructor.name; }
console.log('直接传方法，this 丢失：', lost);
console.log('bind 修正：', callIt(obj.scale.bind(obj)));
console.log('call / apply 指定 this：', obj.scale.call(obj, 4), obj.scale.apply(obj, [4]));
console.log('箭头从外层方法捕获 this：', callIt(obj.makeMultiplier()));
