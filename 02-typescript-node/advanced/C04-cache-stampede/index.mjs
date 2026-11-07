// 缓存击穿：50 个并发请求同时发现 key 失效
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

class FakeDB {
  constructor() { this.hits = 0; }
  async query(key) {
    this.hits += 1;
    await sleep(20);
    return `value-of-${key}`;
  }
}

// 错误写法：查不到就回源
async function unprotectedGet(db, cache, key) {
  if (cache.has(key)) return cache.get(key);
  const v = await db.query(key); // 并发请求全部走到这里
  cache.set(key, v);
  return v;
}

// 正确写法：in-flight Promise 去重（Node 版 singleflight）
class GuardedCache {
  constructor() { this.data = new Map(); this.inflight = new Map(); }
  async get(db, key) {
    if (this.data.has(key)) return this.data.get(key);
    if (this.inflight.has(key)) return this.inflight.get(key);

    const p = (async () => {
      const v = await db.query(key);
      this.data.set(key, v);
      return v;
    })().finally(() => this.inflight.delete(key));

    this.inflight.set(key, p);
    return p;
  }
}

async function run(concurrency, guarded) {
  const db = new FakeDB();
  const start = Date.now();
  if (guarded) {
    const c = new GuardedCache();
    await Promise.all(Array.from({ length: concurrency }, () => c.get(db, 'user:1')));
  } else {
    const cache = new Map();
    await Promise.all(Array.from({ length: concurrency }, () => unprotectedGet(db, cache, 'user:1')));
  }
  return { hits: db.hits, ms: Date.now() - start };
}

const n = 50;
const bad = await run(n, false);
console.log(`无保护：${n} 并发 -> 回源 ${bad.hits} 次，耗时 ${bad.ms}ms`);

const good = await run(n, true);
console.log(`in-flight 去重：${n} 并发 -> 回源 ${good.hits} 次，耗时 ${good.ms}ms`);

console.log('结论：Node 单线程也会击穿——await 让出执行权后，后来者仍会看到「缓存为空」');
