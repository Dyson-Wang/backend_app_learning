// 简化模型：所有请求同时到达，按到达顺序 FIFO 抢池子里的连接
function simulate(durations, pool, waitTimeout) {
  const freeAt = new Array(pool).fill(0);
  let ok = 0;
  let timeout = 0;
  let total = 0;

  for (const d of durations) {
    let best = 0;
    for (let i = 1; i < pool; i++) if (freeAt[i] < freeAt[best]) best = i;
    const wait = freeAt[best];
    if (wait > waitTimeout) {
      timeout++;
      continue;
    }
    freeAt[best] = wait + d;
    total = Math.max(total, freeAt[best]);
    ok++;
  }
  return { ok, timeout, total };
}

const normal = new Array(10).fill(30);

console.log('== 1. 正常：全部 30ms ==');
let r = simulate(normal, 2, 100);
console.log(`成功 = ${r.ok}，超时 = ${r.timeout}，总耗时 = ${r.total}ms`);

console.log('== 2. 混入 1 个 300ms 慢查询（等待上限 100ms） ==');
const slow = [...normal];
slow[0] = 300;
r = simulate(slow, 2, 100);
console.log(`成功 = ${r.ok}，超时 = ${r.timeout}，总耗时 = ${r.total}ms`);

console.log('== 3. 结论 ==');
console.log('池子只有 2 个连接，1 个慢查询就占掉一半；等待超时把雪崩挡在入口');
console.log('慢查询必须配 statement_timeout + 索引，光加连接数只是把问题推后');
