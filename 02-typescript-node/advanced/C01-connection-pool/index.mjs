class Semaphore {
  constructor(size) {
    this.size = size;
    this.free = size;
    this.waiters = [];
  }

  acquire(timeoutMs) {
    if (this.free > 0) {
      this.free -= 1;
      return Promise.resolve(true);
    }
    return new Promise((resolve) => {
      const waiter = { resolve, timer: null };
      if (timeoutMs > 0) {
        waiter.timer = setTimeout(() => {
          const i = this.waiters.indexOf(waiter);
          if (i >= 0) this.waiters.splice(i, 1);
          resolve(false);
        }, timeoutMs);
      }
      this.waiters.push(waiter);
    });
  }

  release() {
    const next = this.waiters.shift();
    if (next) {
      clearTimeout(next.timer);
      next.resolve(true);
    } else {
      this.free += 1;
    }
  }
}

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

async function run(size, requests, workMs, acquireTimeoutMs) {
  const sem = new Semaphore(size);
  let ok = 0;
  let timeout = 0;
  const start = Date.now();

  await Promise.all(
    Array.from({ length: requests }, async () => {
      const got = await sem.acquire(acquireTimeoutMs);
      if (!got) {
        timeout += 1;
        return;
      }
      try {
        await sleep(workMs);
        ok += 1;
      } finally {
        sem.release();
      }
    }),
  );

  return { ok, timeout, ms: Date.now() - start };
}

const SIZE = 3;
const REQUESTS = 10;
const WORK = 100;

console.log('== 1. 无超时：全部排队，慢但都成功 ==');
const a = await run(SIZE, REQUESTS, WORK, 0);
console.log(`ok=${a.ok} timeout=${a.timeout} 耗时=${a.ms}ms`);

console.log('== 2. 等待超过 50ms 就快速失败 ==');
const b = await run(SIZE, REQUESTS, WORK, 50);
console.log(`ok=${b.ok} timeout=${b.timeout} 耗时=${b.ms}ms`);

console.log('== 3. 结论 ==');
console.log('连接池/并发上限必须配等待超时：没有超时，故障时请求无限排队');
console.log('Node 里没有内建连接池，pg.Pool / undici 的 Pool 都要显式配 connectionTimeout');
