class FakeDB {
  constructor() {
    this.queries = 0;
  }

  orders() {
    this.queries += 1;
    return Array.from({ length: 10 }, (_, i) => ({ id: i + 1, userId: (i % 3) + 1 }));
  }

  userById(id) {
    this.queries += 1;
    return { id, name: `user-${id}` };
  }

  userMap(ids) {
    this.queries += 1; // 模拟 SELECT * FROM users WHERE id IN (...)
    const map = new Map();
    for (const id of ids) map.set(id, { id, name: `user-${id}` });
    return map;
  }
}

function ordersNPlusOne(db) {
  const orders = db.orders();
  return orders.map((o) => {
    const u = db.userById(o.userId);
    return `order-${o.id} by ${u.name}`;
  });
}

function ordersBatch(db) {
  const orders = db.orders();
  const ids = [...new Set(orders.map((o) => o.userId))];
  const users = db.userMap(ids);
  return orders.map((o) => `order-${o.id} by ${users.get(o.userId).name}`);
}

const db1 = new FakeDB();
const t0 = Date.now();
const res1 = ordersNPlusOne(db1);
console.log(`N+1：${res1.length} 条订单 -> 查询 ${db1.queries} 次，耗时 ${Date.now() - t0}ms`);

const db2 = new FakeDB();
const t1 = Date.now();
const res2 = ordersBatch(db2);
console.log(`批量：${res2.length} 条订单 -> 查询 ${db2.queries} 次，耗时 ${Date.now() - t1}ms`);

console.log(`第一条结果一致：${res1[0] === res2[0]}`);
console.log('结论：N+1 的代价是「1 + N 次往返」，批量把 N 压成 1');
