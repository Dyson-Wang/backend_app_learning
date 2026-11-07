// 至少一次投递：重复消息必须由消费端去重
const deliveries = [
  { id: 'e1', amount: 10 },
  { id: 'e2', amount: 10 },
  { id: 'e2', amount: 10 },   // ACK 丢了，Broker 重投
  { id: 'e3', amount: 10 },
];

function run(seen) {
  let balance = 100;
  let applied = 0;
  for (const e of deliveries) {
    if (seen) {
      if (seen.has(e.id)) continue;   // 已处理，直接 ACK
      seen.add(e.id);
    }
    balance += e.amount;
    applied++;
  }
  return { balance, applied };
}

console.log('== 1. 消费端没有幂等表 ==');
const a = run(null);
console.log(`投递 ${deliveries.length} 次，余额 100 -> ${a.balance}（多加了 10）`);

console.log('== 2. 消费端按消息 ID 去重 ==');
const b = run(new Set());
console.log(`投递 ${deliveries.length} 次，实际处理 ${b.applied} 次，余额 100 -> ${b.balance}`);

console.log('== 3. 结论 ==');
console.log('至少一次投递下，重复是常态；正确性必须由消费端幂等保证');
console.log('去重表要落库并和业务写在同一事务里，否则还是会重复');
