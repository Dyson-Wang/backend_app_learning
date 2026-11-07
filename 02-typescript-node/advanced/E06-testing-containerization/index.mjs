// 极简测试框架：只看失败输出长什么样
let failures = 0;

function check(name, ok, detail) {
  if (ok) {
    console.log(`PASS ${name}`);
    return;
  }
  failures++;
  console.log(`FAIL ${name}: ${detail}`);
}

const add = (a, b) => a + b;

function divide(a, b) {
  if (b === 0) throw new Error('division by zero');
  return a / b;
}

function parseAge(s) {
  if (s === '') throw new Error('empty');
  return 18;                 // 故意不校验非数字，用来触发失败用例
}

console.log('== 1. 单元测试 ==');
check('add(1,2)=3', add(1, 2) === 3, '期望 3');

let threw = false;
try {
  divide(1, 0);
} catch {
  threw = true;
}
check('divide(1,0) 抛错', threw, '期望抛错');

let ageErr = null;
try {
  parseAge('abc');
} catch (e) {
  ageErr = e;
}
check('parseAge("abc") 抛错', ageErr !== null, '期望抛错，实际正常返回 18');

console.log('== 2. 汇总 ==');
console.log(`3 个用例：${3 - failures} 通过，${failures} 失败`);

console.log('== 3. 结论 ==');
console.log('最后一条是故意失败的用例，用来看失败输出长什么样');
console.log('单测不碰网络/DB；集成测试才起真实依赖，两者分开跑');
console.log('容器化：多阶段构建、非 root、HEALTHCHECK 打 /healthz');
