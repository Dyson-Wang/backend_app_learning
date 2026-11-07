import * as pathNS from 'node:path';

// 单文件里要真实演示「模块」，用 data: URL 动态 import 造出独立的 ES 模块实例。
// 每个 data: URL 都有自己的模块作用域、导出表与 Module Namespace 对象，
// 语义和磁盘上的 .mjs 完全一致；它是内联语言特性，不是第三方包。
const load = (code) => import('data:text/javascript,' + encodeURIComponent(code));

// —— 导出 / 私有规则：没 export 的绑定对外完全不可见 ——
const modA = await load(`
  export const pub = '公开值';
  const secret = '内部私有';            // 没导出，外部拿不到
  export function reveal() { return secret; }
`);
console.log('导出可见：', modA.pub, '；私有绑定不可见：', modA.secret);
console.log('但模块内部能通过函数间接访问：', modA.reveal());

// —— 命名空间对象：不可写、键有序、带 Module 标记 ——
console.log('命名空间标记：', modA[Symbol.toStringTag], '；键：', Object.keys(modA).join(','));
let writeErr;
try { modA.pub = '篡改'; } catch (err) { writeErr = err.constructor.name; }
console.log('命名空间只读，赋值报：', writeErr);
const builtinNS = await import('node:path');
console.log('动态 import 内置模块也是命名空间：', builtinNS[Symbol.toStringTag], '；含 join：', typeof builtinNS.join);

// —— 命名冲突与别名：import { x as y } 的运行时等价是解构重命名 ——
const modB = await load(`export const value = 'B 的 value'; export default 42;`);
const modC = await load(`export const value = 'C 的 value';`);
const { value: valueFromB } = modB;
const { value: valueFromC } = modC;
console.log('别名消歧：', valueFromB, '|', valueFromC, '；default：', modB.default);

// 模块内部也能用 export { a as b } 重命名导出
const modD = await load(`const a = 1; export { a as renamed };`);
console.log('模块内 export 别名：', modD.renamed, '；原名 a 不可见：', modD.a);

// —— 活绑定（live binding）：export let 的变化会同步到命名空间 ——
const modE = await load(`
  export let counter = 0;
  export function inc() { counter++; }
`);
console.log('活绑定初值：', modE.counter);
modE.inc(); modE.inc();
console.log('调用 inc 后命名空间同步更新：', modE.counter);

// —— import * as ns：内置模块的命名空间对象 ——
console.log('import * as pathNS 的键数：', Object.keys(pathNS).length, '；basename：', pathNS.basename('/a/b.txt'));

// —— 目录即包 / 内部包可见性 / 循环依赖（概念，单文件无法真造多文件）——
// 真实工程里：含 package.json 的目录就是包；"exports" 字段精确列出对外子路径，
// 未列出的内部文件被 import 会报 ERR_PACKAGE_PATH_NOT_EXPORTED，这就是可见性边界。
// 循环依赖：A import B、B import A 时，ESM 靠「活绑定 + 实例化与求值分离」避免死锁；
// 但若在模块求值阶段就读取对方尚未初始化的 let/const，会命中 TDZ 抛
// ReferenceError（Cannot access 'x' before initialization）。Node 能加载循环，但要避免顶层读取。
console.log('可见性边界由 exports 控制，import.meta.resolve 可探测：', import.meta.resolve('node:path'));
console.log('循环依赖：ESM 活绑定可加载，但求值期读取未初始化变量会触发 TDZ 报错');
