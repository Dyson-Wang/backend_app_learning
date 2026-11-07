import assert from 'node:assert/strict';

// —— 函数形态：声明 / 表达式 / 箭头 / 参数特性 ——
function add(a, b = 1) {              // 默认参数
  return a + b;
}
const mul = (a, b) => a * b;          // 箭头函数，隐式返回
const sumAll = (...nums) => nums.reduce((s, x) => s + x, 0);  // rest 参数
function describe({ name, age = 0 }) { // 解构参数
  return `${name}:${age}`;
}
console.log('默认参数 add(5)：', add(5));
console.log('箭头函数 mul(3,4)：', mul(3, 4));
console.log('rest 参数求和：', sumAll(1, 2, 3, 4));
console.log('解构参数：', describe({ name: 'alice' }));

// 箭头函数没有自己的 arguments；普通函数有类数组 arguments
function legacy() { return arguments.length; }
const arrow = (...a) => a.length;
console.log('普通函数 arguments：', legacy(1, 2, 3), '，箭头用 rest：', arrow(1, 2, 3));

// —— 自定义错误类型：继承 Error，用 { cause } 串起错误链 ——
class ValidationError extends Error {
  constructor(field, cause) {
    super(`字段 ${field} 校验失败`, { cause });
    this.name = 'ValidationError';
    this.field = field;
  }
}

function parseAge(input) {
  const age = Number(input);
  if (!Number.isInteger(age) || age < 0) {
    throw new ValidationError('age', `收到非法值 ${input}`);
  }
  return age;
}

// —— try / catch / finally：finally 无论是否 return 都会执行 ——
function handle(input) {
  try {
    const age = parseAge(input);
    assert.ok(age >= 0);
    return `成功：age=${age}`;
  } catch (err) {
    if (err instanceof ValidationError) {
      return `捕获自定义错误：${err.name}/${err.field}，cause=${err.cause}`;
    }
    throw err; // 认不出的错误继续向上抛
  } finally {
    console.log('  [finally] 本次处理结束，input =', input);
  }
}

console.log('合法输入：', handle('18'));
console.log('非法输入：', handle('abc'));

// 可省略 catch 绑定名；未捕获异常会终结进程，生产必须显式兜底
try {
  JSON.parse('{ 坏 JSON');
} catch {
  console.log('JSON 解析失败被 catch{} 接住');
}

// Error 的属性：name / message / stack / cause
const wrapped = new Error('外层', { cause: new Error('根因') });
console.log('错误链：', wrapped.message, '->', wrapped.cause.message);
console.log('错误类型名：', wrapped.name, '，含 stack：', wrapped.stack.includes('外层'));
