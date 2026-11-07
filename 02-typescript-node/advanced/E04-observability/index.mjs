import { AsyncLocalStorage } from 'node:async_hooks';

// 结构化日志：一行一个 JSON
function log(fields) {
  console.log(JSON.stringify({ level: 'info', ...fields }));
}

// trace id：入口生成，通过 AsyncLocalStorage 透传
const als = new AsyncLocalStorage();
let seq = 0;
const nextTraceId = () => `trace-${String((seq += 1)).padStart(4, '0')}`;

// RED 指标
class Metrics {
  constructor() { this.requests = 0; this.errors = 0; this.totalMs = 0; this.maxMs = 0; }
  observe(status, ms) {
    this.requests += 1;
    if (status >= 500) this.errors += 1;
    this.totalMs += ms;
    if (ms > this.maxMs) this.maxMs = ms;
  }
  snapshot() {
    const avg = this.requests ? this.totalMs / this.requests : 0;
    return { requests: this.requests, errors: this.errors, avg, max: this.maxMs };
  }
}

const metrics = new Metrics();

function handle({ method, path, status, ms }) {
  const traceId = nextTraceId();
  return als.run({ traceId }, async () => {
    metrics.observe(status, ms);
    const entry = { trace_id: als.getStore().traceId, method, path, status, latency_ms: ms };
    if (status >= 500) {
      log({ ...entry, level: 'error', msg: 'upstream failed' });
    } else {
      log(entry);
    }
  });
}

const reqs = [
  { method: 'GET', path: '/healthz', status: 200, ms: 1 },
  { method: 'GET', path: '/users', status: 200, ms: 12 },
  { method: 'POST', path: '/users', status: 201, ms: 18 },
  { method: 'GET', path: '/users/999', status: 404, ms: 3 },
  { method: 'POST', path: '/users', status: 400, ms: 2 },
  { method: 'GET', path: '/orders', status: 500, ms: 240 },
  { method: 'GET', path: '/users', status: 200, ms: 9 },
  { method: 'GET', path: '/users/1', status: 200, ms: 7 },
  { method: 'GET', path: '/orders', status: 500, ms: 260 },
  { method: 'DELETE', path: '/users/1', status: 204, ms: 5 },
];

for (const r of reqs) {
  await handle(r);
}

const { requests, errors, avg, max } = metrics.snapshot();
console.log('---- RED 指标 ----');
console.log(`requests=${requests} errors=${errors} error_rate=${((errors / requests) * 100).toFixed(2)}% avg_ms=${avg.toFixed(1)} max_ms=${max}`);
console.log('排查顺序：先看 error_rate 与 max_ms，再用 trace_id 串起同一次请求的全部日志');
