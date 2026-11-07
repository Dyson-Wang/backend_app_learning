// JS 运行时没有泛型：TypeScript 的 <T> 编译后被完全擦除。
// 本文件用 JSDoc 的 @template 表达类型参数（配合 tsc --checkJs 可做静态检查），
// 但 node 直接执行时这些注释不产生任何运行时行为——这正说明「类型擦除」。
// TypeScript 等价写法：
//   function identity<T>(x: T): T { return x; }
//   function first<T>(xs: T[]): T | undefined { return xs[0]; }
//   function mapArray<T, U>(xs: T[], f: (x: T) => U): U[] { return xs.map(f); }
//   class Stack<T> { private items: T[] = []; }

/**
 * 恒等函数：泛型函数的最简形态
 * @template T
 * @param {T} value
 * @returns {T}
 */
function identity(value) {
  return value;
}

/**
 * 带约束的泛型：T 限定为 string 或 number（对应 TS 的 <T extends ...>）
 * @template {string | number} T
 * @param {T} a
 * @param {T} b
 * @returns {T[]}
 */
function sortedPair(a, b) {
  return [a, b].sort();
}

/**
 * 多类型参数 + 泛型集合：把 T[] 映射成 U[]
 * @template T, U
 * @param {T[]} items
 * @param {(x: T) => U} fn
 * @returns {U[]}
 */
function mapArray(items, fn) {
  return items.map(fn);
}

console.log('泛型函数 identity：', identity(42), identity('文本'), identity({ k: 1 }).k);
console.log('约束泛型 sortedPair：', sortedPair(3, 1).join(','), sortedPair('b', 'a').join(','));
console.log('多参数 + 泛型集合：', mapArray([1, 2, 3], (x) => x * 2).join(','));

// —— 泛型类型（泛型容器）：用 JSDoc @template 标注一个栈 ——
/** @template T */
class Stack {
  /** @type {T[]} */
  #items = [];
  /** @param {T} x */
  push(x) { this.#items.push(x); return this; }
  /** @returns {T | undefined} */
  pop() { return this.#items.pop(); }
  /** @returns {number} */
  get size() { return this.#items.length; }
}
const numStack = new Stack();
numStack.push(10).push(20);
console.log('泛型容器 Stack：size =', numStack.size, '，pop =', numStack.pop());
// 运行时根本不检查 T：塞字符串也照样跑，类型信息已擦除
numStack.push('字符串也进来了');
console.log('无运行时类型约束，任意类型都能 push：', numStack.pop());

// —— 类型推断 & 类型擦除的实证 ——
// TS 泛型只在编译期存在；JS 里 identity 只有一个函数体，对所有类型复用
console.log('identity 是同一个函数：', identity === identity, '；形参个数：', identity.length);
console.log('函数源码里没有 <T>：', /<T>/.test(identity.toString()));
console.log('擦除后按值取类型：', typeof identity(1), typeof identity('s'), Array.isArray(identity([])));

/** @type {Map<string, number>} */
const scores = new Map([['alice', 90], ['bob', 85]]);
console.log('Map<K,V> 运行时只是普通 Map：', scores.constructor.name, '，alice =', scores.get('alice'));

// 单态化 vs 类型擦除：Rust/C++ 会为每种 T 生成一份机器码（单态化）；
// Java/TS 擦除成同一份代码；JS 连注解都不保留，最彻底。
console.log('结论：JS 泛型 = 编译期注释，运行时零开销、零检查、零特化');
