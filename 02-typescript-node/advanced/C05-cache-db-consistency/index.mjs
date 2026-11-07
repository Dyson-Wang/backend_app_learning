// 缓存与 DB 双写一致性
class Cache {
  constructor() {
    this.data = new Map();
  }
  get(k) {
    return this.data.get(k);
  }
  set(k, v) {
    this.data.set(k, v);
  }
  del(k) {
    this.data.delete(k);
  }
}

function readThrough(cache, db, k) {
  const hit = cache.get(k);
  if (hit !== undefined) return hit;
  const v = db.get(k);      // 缓存未命中，回源
  cache.set(k, v);
  return v;
}

const key = 'user:1';
const fresh = (v) => ({ db: new Map([[key, v]]), cache: new Cache() });

console.log('== 1. 先删缓存，再写库 ==');
{
  const { db, cache } = fresh('v1');
  cache.set(key, 'v1');
  cache.del(key);                       // 写请求：先删缓存
  readThrough(cache, db, key);          // 读请求：回源读到旧值
  db.set(key, 'v2');                    // 写请求：再写库
  console.log(`缓存 = ${cache.get(key)}，库 = ${db.get(key)}，是否脏 = ${cache.get(key) !== db.get(key)}`);
}

console.log('== 2. 先写库，再删缓存 ==');
{
  const { db, cache } = fresh('v1');
  cache.set(key, 'v1');
  const old = cache.get(key);           // 读请求先拿到旧值
  db.set(key, 'v2');                    // 写请求：先写库
  cache.del(key);                       // 再删缓存
  cache.set(key, old);                  // 读请求迟到，把旧值写回
  console.log(`缓存 = ${cache.get(key)}，库 = ${db.get(key)}，是否脏 = ${cache.get(key) !== db.get(key)}`);
}

console.log('== 3. 延迟双删 ==');
{
  const { db, cache } = fresh('v1');
  cache.set(key, 'v1');
  const old = cache.get(key);
  cache.del(key);
  db.set(key, 'v2');
  cache.set(key, old);                  // 迟到的读请求仍回填旧值
  cache.del(key);                       // 延迟再删一次
  const hit = cache.get(key);
  console.log(`缓存命中 = ${hit !== undefined}，下一次读会回源拿到 ${readThrough(cache, db, key)}`);
}

console.log('== 4. 结论 ==');
console.log('先删缓存有窗口，先写库再删缓存也有窗口，双删只能把窗口压小');
console.log('要强一致就别用缓存读；要最终一致就配 binlog 订阅 + 重试删除');
