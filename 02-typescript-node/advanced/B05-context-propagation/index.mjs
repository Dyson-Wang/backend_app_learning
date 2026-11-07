const sleep = (ms, signal) =>
  new Promise((resolve, reject) => {
    if (signal?.aborted) return reject(signal.reason ?? new Error('aborted'));
    const timer = setTimeout(() => {
      signal?.removeEventListener('abort', onAbort);
      resolve();
    }, ms);
    const onAbort = () => {
      clearTimeout(timer);
      reject(signal.reason ?? new Error('aborted'));
    };
    signal?.addEventListener('abort', onAbort, { once: true });
  });

// repo 尊重 signal：200ms 的活，取消就立刻退出
async function repo(signal) {
  await sleep(200, signal);
  return 'data';
}

async function service(signal) {
  if (signal.aborted) throw signal.reason ?? new Error('aborted');
  return repo(signal);
}

async function handler(signal) {
  return service(signal);
}

async function main() {
  console.log('== 1. 上游超时 50ms，下游要 200ms ==');
  {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(new Error('deadline exceeded')), 50);
    const start = Date.now();
    try {
      await handler(controller.signal);
      console.log('未取消（不应发生）');
    } catch (e) {
      console.log(`err = ${e.message}`);
    } finally {
      clearTimeout(timer);
    }
    console.log(`耗时 = ${Date.now() - start}ms（远小于 200ms）`);
  }

  console.log('== 2. 客户端主动取消 ==');
  {
    const controller = new AbortController();
    setTimeout(() => controller.abort(new Error('client closed')), 20);
    const start = Date.now();
    try {
      await handler(controller.signal);
    } catch (e) {
      console.log(`err = ${e.message}，耗时 = ${Date.now() - start}ms`);
    }
  }

  console.log('== 3. 反例：中间层丢掉 signal ==');
  {
    const start = Date.now();
    await repo(undefined); // 不传 signal，取消失效
    console.log(`耗时 = ${Date.now() - start}ms（取消了上游也没用）`);
  }

  console.log('== 4. 结论 ==');
  console.log('signal 必须一路显式传；任何一层忘记传，取消就断在这里');
  console.log('fetch 直接支持 signal，自定义异步函数要自己监听 abort 事件');
}

await main();
