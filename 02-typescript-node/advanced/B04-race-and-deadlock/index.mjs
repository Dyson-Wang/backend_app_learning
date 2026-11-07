// 1. 异步交错导致丢更新：两个 async 函数都读到旧值
let counter = 0;

async function racyInc() {
  const v = counter;
  await Promise.resolve(); // 让出执行权，制造交错
  counter = v + 1;
}

// 2. 串行化：用一条 Promise 链保证「读-改-写」不被打断
let safeCounter = 0;
let chain = Promise.resolve();

function safeInc() {
  chain = chain.then(() => {
    safeCounter += 1;
  });
  return chain;
}

// 3. 死锁：两个 Promise 互相等待
async function deadlockDemo() {
  let resolveA;
  let resolveB;
  const a = new Promise((r) => { resolveA = r; });
  const b = new Promise((r) => { resolveB = r; });

  const t1 = (async () => { await b; resolveA(); })();
  const t2 = (async () => { await a; resolveB(); })();

  const timeout = new Promise((r) => setTimeout(() => r('timeout'), 100));
  const result = await Promise.race([Promise.all([t1, t2]).then(() => 'done'), timeout]);
  return result;
}

async function main() {
  console.log('== 1. 异步交错：两个并发各加一次 ==');
  await Promise.all([racyInc(), racyInc()]);
  console.log(`期望 2，实际 ${counter}（丢更新）`);

  console.log('== 2. Promise 链串行化：1000 次 ==');
  for (let i = 0; i < 1000; i += 1) safeInc();
  await chain;
  console.log(`期望 1000，实际 ${safeCounter}`);

  console.log('== 3. 互相 await 造成死锁 ==');
  const r = await deadlockDemo();
  console.log(`结果: ${r}（100ms 内两个 Promise 都没完成）`);

  console.log('结论：Node 没有数据竞争，但有异步交错；死锁表现为 Promise 永不 settle');
}

await main();
