// 顺序性：状态机 + 分区
const apply = (state, e) => {
  if (e.kind === 'create' && state === '') return ['CREATED', ''];
  if (e.kind === 'pay' && state === 'CREATED') return ['PAID', ''];
  if (e.kind === 'ship' && state === 'PAID') return ['SHIPPED', ''];
  return [state, `${e.kind} 在 ${state || 'INIT'} 状态下不合法`];
};

const consume = (events) => {
  let state = '';
  for (const e of events) {
    const [next, err] = apply(state, e);
    state = next;
    if (err) return [state, err];
  }
  return [state, ''];
};

const partition = (key, n) => {
  let h = 0;
  for (const ch of key) h = (h * 31 + ch.codePointAt(0)) % n;
  return h;
};

const create = { id: 'c1', kind: 'create' };
const pay = { id: 'c2', kind: 'pay' };
const ship = { id: 'c3', kind: 'ship' };

console.log('== 1. 同一 key 被并发消费（顺序打乱） ==');
let [state, err] = consume([ship, create, pay]);
console.log(`最终状态 = ${state || 'INIT'}，错误 = ${err}`);

console.log('== 2. 按 key 路由到同一分区、单消费者串行处理 ==');
[state, err] = consume([create, pay, ship]);
console.log(`最终状态 = ${state}，错误 = ${err === '' ? '无' : err}`);

console.log('== 3. 分区函数：同一 order_id 必然落到同一分区 ==');
for (const k of ['order-42', 'order-43', 'order-42']) {
  console.log(`${k} -> 分区 ${partition(k, 4)}`);
}

console.log('== 4. 结论 ==');
console.log('顺序性不是全局的，是「同一分区内」的；跨 key 的顺序本来就没有意义');
console.log('要顺序就牺牲并发：单分区单消费者，或者消费端按 key 做串行队列');
