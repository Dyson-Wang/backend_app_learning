import assert from 'node:assert/strict';
import { performance } from 'node:perf_hooks';

// —— 文档注释：JSDoc 描述用途/参数/返回，供 IDE 与文档生成器读取 ——
/**
 * 计算字符串的码元长度
 * @param {string} s 输入字符串
 * @returns {number} 长度
 */
function charCount(s) { return s.length; }

/**
 * 安全除法，除数为 0 时抛错
 * @param {number} a 被除数
 * @param {number} b 除数
 * @returns {number} 商
 */
function divide(a, b) {
  if (b === 0) throw new RangeError('除数不能为 0');
  return a / b;
}

// —— 表驱动测试：一组用例批量验证 ——
const cases = [
  { name: 'charCount 空串', input: [''], run: charCount, expected: 0 },
  { name: 'charCount 中文', input: ['后端'], run: charCount, expected: 2 },
  { name: 'divide 正常', input: [10, 2], run: divide, expected: 5 },
  { name: 'divide 小数', input: [7, 2], run: divide, expected: 3.5 },
];

let passed = 0;
const failures = [];
for (const c of cases) {
  try {
    assert.strictEqual(c.run(...c.input), c.expected);
    passed++;
    console.log('  [通过]', c.name);
  } catch (err) {
    failures.push(`${c.name} -> ${err.name}`);
  }
}
console.log('表驱动结果：通过', passed, '/', cases.length, '，失败', failures.length);

// —— 失败信息：故意埋一个错误用例，看 assert 的 actual/expected/operator ——
let failureInfo = '';
try {
  assert.strictEqual(charCount('abc'), 99);  // 故意写错，演示失败输出
} catch (err) {
  failureInfo = `actual=${err.actual}, expected=${err.expected}, operator=${err.operator}`;
}
console.log('断言失败信息示例：', failureInfo);

// —— 更多断言形式：deepStrictEqual / throws / ok ——
assert.deepStrictEqual({ a: 1, b: [2] }, { a: 1, b: [2] });
assert.throws(() => divide(1, 0), RangeError);
assert.ok(Array.isArray([1, 2]));
console.log('其余断言全部通过：deepStrictEqual / throws / ok');

// —— 基准测试：固定迭代次数，用 performance.now() 量「差值」——
const ITERATIONS = 20000;
let acc = 0;
const t0 = performance.now();
for (let i = 0; i < ITERATIONS; i++) acc += charCount('abc');
const elapsed = performance.now() - t0;
console.log('基准完成：循环', ITERATIONS, '次，累加结果', acc, '（只演示 API，不依赖具体耗时）');
console.log('performance.now 返回数值：', typeof elapsed === 'number', '，且非负：', elapsed >= 0);

// —— 覆盖率思路（工具层，本脚本不依赖第三方）——
// node --experimental-test-coverage --test 会打印 V8 覆盖率；
// NODE_V8_COVERAGE 环境变量可导出原始覆盖率数据，再用 c8 / nyc 生成报告。
// 标准测试写法是 node:test 的 test()/describe() + `node --test`；
// 这里用自制迷你 runner，保证 node index.mjs 直接可跑且退出码 0。
console.log('文档：JSDoc 由 IDE/jsdoc 读取；TS 下 tsc --declaration 生成 .d.ts');
console.log('覆盖率：node --experimental-test-coverage / NODE_V8_COVERAGE + c8');
console.log('入口脚本仍可直接运行，全部断言通过，退出码 0');
