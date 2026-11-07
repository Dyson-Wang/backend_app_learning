import v8 from 'node:v8';
import process from 'node:process';

// —— 值与引用：原始值拷贝，对象拷贝的是「指向堆内存的引用」——
let a = 10;
let b = a;                 // 原始值：复制一份，互不影响
b = 99;
console.log('原始值拷贝：a =', a, '，b =', b);

const obj1 = { n: 1 };
const obj2 = obj1;         // 对象：两个变量指向同一块堆内存
obj2.n = 42;
console.log('对象引用共享：obj1.n =', obj1.n, '，是同一对象吗：', obj1 === obj2);

const origin = { inner: { v: 1 } };
const shallow = { ...origin };
shallow.inner.v = 5;       // 浅拷贝：嵌套对象仍共享
console.log('浅拷贝嵌套仍共享：origin.inner.v =', origin.inner.v);
const source = { x: [1, { y: 2 }] };
const deep = structuredClone(source);
deep.x[1].y = 99;          // 深拷贝：完全独立
console.log('structuredClone 深拷贝：source =', JSON.stringify(source), '，copy =', JSON.stringify(deep));

// —— 栈与堆：V8 分代堆 + 堆外内存 ——
const stats = v8.getHeapStatistics();
console.log('V8 堆存在：', stats.total_heap_size > 0, '；heap_size_limit 有限：', Number.isFinite(stats.heap_size_limit));
console.log('memoryUsage 含 heapUsed：', Object.keys(process.memoryUsage()).includes('heapUsed'));
const buf = Buffer.alloc(1024);
console.log('Buffer 在堆外内存，byteLength =', buf.length, '；属于 Uint8Array：', buf instanceof Uint8Array);

// —— GC：追踪式分代回收（非引用计数），从 GC Roots 标记可达对象 ——
// V8 新生代 Scavenge、老生代 Mark-Sweep/Mark-Compact；开发者不能手动 free，
// 也不像 Rust 有所有权/借用/移动语义，对象全靠引用与可达性管理生命周期。
console.log('WeakRef / FinalizationRegistry 可用：', typeof WeakRef, typeof FinalizationRegistry);
let strong = { data: '大对象' };
const weak = new WeakRef(strong);
console.log('强引用还在时 deref 拿到：', weak.deref()?.data);
strong = null;             // 去掉最后一个强引用 → 对象不可达，等待 GC（时机不确定）
console.log('去掉强引用后对象变为不可达；GC 何时回收不确定，只演示 API 不依赖结果');

// WeakMap / WeakSet 的键是弱引用，不阻止回收，也不可遍历
const wm = new WeakMap();
let key = { id: 1 };
wm.set(key, '元数据');
console.log('WeakMap 键存活时 has/get：', wm.has(key), wm.get(key));
console.log('WeakMap 不可遍历（无 size / iterator）：', typeof wm.size, typeof wm[Symbol.iterator]);
key = null;                // 键不可达后，该条目可被 GC 回收

// —— 常见内存陷阱与对象生命周期 ——
const cache = new Map();   // 陷阱1：全局缓存无限增长，强引用让对象永远可达
for (let i = 0; i < 3; i++) cache.set(i, { payload: 'x'.repeat(10) });
console.log('全局 Map 缓存条目数：', cache.size);
cache.clear();             // 显式释放强引用后条目才可回收
console.log('clear 后条目数：', cache.size);
function makeClosure() {   // 陷阱2：闭包捕获大对象，闭包存活期间无法回收
  const big = { data: Array(1000).fill(0) };
  return () => big.data.length;
}
const closure = makeClosure();
console.log('闭包持有大对象，调用结果：', closure(), '（闭包活着对象就活着）');
const timer = setTimeout(() => {}, 1000); // 陷阱3：未清理的 timer/listener 会长期持有引用
console.log('定时器句柄是对象：', typeof timer === 'object');
clearTimeout(timer);       // 用完必须清理，否则阻止进程退出与回收
console.log('已 clearTimeout 清理定时器');
console.log('JS 没有所有权/移动/借用；V8 的逃逸分析会把短命对象优化到栈上，但对开发者透明');
