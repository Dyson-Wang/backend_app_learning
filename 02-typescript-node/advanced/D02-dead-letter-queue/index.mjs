// 死信队列与重试策略
function deliver(msg, maxAttempts) {
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    if (msg.permanent) {
      console.log(`${msg.id} -> 永久错误，直接进死信（不重试）`);
      return { success: false, attempts: attempt };
    }
    if (attempt > msg.failures) {
      console.log(`${msg.id} -> 成功（第 ${attempt} 次）`);
      return { success: true, attempts: attempt };
    }
    if (attempt < maxAttempts) {
      const backoff = 10 * 2 ** (attempt - 1);
      console.log(`${msg.id} -> 可重试错误，退避 ${backoff}ms 后重试`);
      continue;
    }
    console.log(`${msg.id} -> 重试耗尽，进死信`);
  }
  return { success: false, attempts: maxAttempts };
}

const msgs = [
  { id: 'm1', failures: 0, permanent: false },
  { id: 'm2', failures: 2, permanent: false },
  { id: 'm3', failures: 0, permanent: true },
  { id: 'm4', failures: 99, permanent: false },
  { id: 'm5', failures: 0, permanent: false },
];

console.log('== 1. 逐条投递（最多 3 次） ==');
let ok = 0;
let dead = 0;
let total = 0;
for (const m of msgs) {
  const r = deliver(m, 3);
  total += r.attempts;
  if (r.success) ok++;
  else dead++;
}

console.log('== 2. 统计 ==');
console.log(`成功 = ${ok}，死信 = ${dead}，总投递 = ${total}`);

console.log('== 3. 结论 ==');
console.log('永久错误不要重试，重试只会拖慢整个队列');
console.log('死信队列必须有人看：它进 DLQ 说明要么代码有 bug，要么数据有问题');
