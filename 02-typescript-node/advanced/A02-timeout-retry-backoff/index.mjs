const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

// 可重试错误带 retryable 标记；永久错误不重试
class DownstreamError extends Error {
  constructor(message, retryable) {
    super(message);
    this.retryable = retryable;
  }
}

async function doWithRetry({ maxAttempts, baseMs, fn }) {
  let lastError = null;
  for (let attempt = 1; attempt <= maxAttempts; attempt += 1) {
    try {
      const value = await fn(attempt);
      return { ok: true, attempt, value };
    } catch (e) {
      lastError = e;
      if (!e.retryable) {
        return { ok: false, attempt, error: e.message, permanent: true };
      }
      if (attempt === maxAttempts) break;
      const delay = baseMs * 2 ** (attempt - 1);
      console.log(`  attempt ${attempt} 失败（${e.message}），${delay}ms 后重试`);
      await sleep(delay);
    }
  }
  return { ok: false, attempt: maxAttempts, error: lastError.message, permanent: false };
}

// case 1：前 3 次失败，第 4 次成功
let calls = 0;
const flaky = async () => {
  calls += 1;
  if (calls < 4) throw new DownstreamError('503 upstream', true);
  return `第 ${calls} 次成功`;
};

// case 2：永久错误，不重试
const permanent = async () => {
  throw new DownstreamError('400 bad request', false);
};

// case 3：一直失败，耗尽重试
const alwaysFail = async () => {
  throw new DownstreamError('timeout', true);
};

async function main() {
  console.log('== 1. 可重试错误：指数退避后成功 ==');
  const r1 = await doWithRetry({ maxAttempts: 4, baseMs: 10, fn: flaky });
  console.log('结果:', r1);

  console.log('== 2. 永久错误：立即停止 ==');
  const r2 = await doWithRetry({ maxAttempts: 4, baseMs: 10, fn: permanent });
  console.log('结果:', r2);

  console.log('== 3. 一直失败：耗尽 3 次 ==');
  const r3 = await doWithRetry({ maxAttempts: 3, baseMs: 5, fn: alwaysFail });
  console.log('结果:', r3);

  console.log('结论：只重试可重试错误，退避必须指数增长；生产还要加抖动，避免重试风暴');
}

await main();
