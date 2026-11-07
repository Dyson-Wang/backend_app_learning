// 幂等键：同一个 Idempotency-Key 只创建一次订单，其余回放结果
const orders = new Map();     // key -> order
const inflight = new Map();   // key -> Promise

async function createOrder(key, amount) {
  if (orders.has(key)) {
    return { order: orders.get(key), replayed: true };
  }
  if (inflight.has(key)) {
    // 同 key 的并发请求：等第一个跑完，直接拿结果，不重复建单
    const order = await inflight.get(key);
    return { order, replayed: true };
  }

  const p = (async () => {
    await new Promise((r) => setTimeout(r, 20));   // 模拟写库耗时
    const order = { id: orders.size + 1, key, amount };
    orders.set(key, order);
    return order;
  })();

  inflight.set(key, p);
  try {
    return { order: await p, replayed: false };
  } finally {
    inflight.delete(key);
  }
}

const key = 'idem-7f3a';
const results = await Promise.all(
  Array.from({ length: 20 }, () => createOrder(key, 9900)),
);

const created = results.filter((r) => !r.replayed).length;
const replayed = results.filter((r) => r.replayed).length;
const ids = new Set(results.map((r) => r.order.id));

console.log('== 1. 20 个并发请求，同一个 Idempotency-Key ==');
console.log(`真实建单 = ${created}，回放 = ${replayed}，订单 id 去重后 = ${ids.size}`);

console.log('== 2. 换一个 key ==');
const r2 = await createOrder('idem-other', 100);
console.log(`新 key 建单 id = ${r2.order.id}，replayed = ${r2.replayed}`);

console.log('== 3. 结论 ==');
console.log('幂等键必须落库并加唯一索引；内存 Map 只是演示，进程重启就没了');
console.log('并发同 key 要用「占位 + 等待」而不是「先查再写」，否则照样重复');
