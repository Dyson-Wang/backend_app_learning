import http from 'node:http';
import { performance } from 'node:perf_hooks';

// —— 枚举（13 课）：Object.freeze 定义错误码常量 ——
const ErrorCode = Object.freeze({
  VALIDATION: 'VALIDATION',
  NOT_FOUND: 'NOT_FOUND',
  METHOD_NOT_ALLOWED: 'METHOD_NOT_ALLOWED',
  BAD_REQUEST: 'BAD_REQUEST',
});

// —— 自定义错误（04 课）：继承 Error，带上 status / code ——
class HttpError extends Error {
  constructor(status, code, message) {
    super(message);
    this.name = 'HttpError';
    this.status = status;
    this.code = code;
  }
}

// —— 泛型（14 课）：JSDoc @template 表达 T，运行时类型被擦除 ——
/**
 * 统一 JSON 响应
 * @template T
 * @param {object} res ServerResponse
 * @param {number} status 状态码
 * @param {T} payload 任意可序列化载荷
 */
function send(res, status, payload) {
  const body = JSON.stringify(payload);
  res.writeHead(status, {
    'Content-Type': 'application/json; charset=utf-8',
    'Content-Length': Buffer.byteLength(body),
  });
  res.end(body);
}
function fail(res, status, code, message) {
  send(res, status, { error: { code, message } });
}

// —— 08/09 课：异步迭代消费请求体并解析 JSON ——
async function readJSON(req) {
  const chunks = [];
  for await (const chunk of req) chunks.push(chunk);
  const text = Buffer.concat(chunks).toString('utf8');
  return text ? JSON.parse(text) : {};
}

const durations = [];   // 收集每次耗时「差值」，保证打印内容可复现

// —— 闭包 + 集合（16/05 课）：工厂函数捕获内存 store —
function createRouter() {
  const store = new Map();   // 内存存储，单线程无需锁（19 课）
  let nextId = 1;
  let seq = 0;               // 确定性 request id，不用随机值

  return async function route(req, res) {
    const requestId = `req-${++seq}`;
    const start = performance.now();
    const url = new URL(req.url, 'http://127.0.0.1');
    try {
      if (url.pathname === '/healthz') {
        if (req.method !== 'GET') throw new HttpError(405, ErrorCode.METHOD_NOT_ALLOWED, 'use GET');
        return send(res, 200, { status: 'ok' });
      }
      if (url.pathname === '/users') {
        if (req.method === 'GET') return send(res, 200, { data: [...store.values()] });
        if (req.method === 'POST') {
          const input = await readJSON(req);
          if (!input.name || !String(input.name).trim()) {           // 校验
            throw new HttpError(400, ErrorCode.VALIDATION, 'name is required');
          }
          const user = { id: nextId++, name: String(input.name) };
          store.set(user.id, user);
          return send(res, 201, { data: user });
        }
        res.setHeader('Allow', 'GET, POST');
        throw new HttpError(405, ErrorCode.METHOD_NOT_ALLOWED, 'use GET or POST');
      }
      const match = url.pathname.match(/^\/users\/(\d+)$/);          // 18 课正则
      if (match) {
        if (req.method !== 'GET') throw new HttpError(405, ErrorCode.METHOD_NOT_ALLOWED, 'use GET');
        const user = store.get(Number(match[1]));
        if (!user) throw new HttpError(404, ErrorCode.NOT_FOUND, 'user not found');
        return send(res, 200, { data: user });
      }
      throw new HttpError(404, ErrorCode.NOT_FOUND, 'route not found');
    } catch (err) {
      if (err instanceof HttpError) return fail(res, err.status, err.code, err.message);
      return fail(res, 400, ErrorCode.BAD_REQUEST, err.message);     // 统一错误码
    } finally {
      durations.push(performance.now() - start);                     // 18/23 课：单调时钟差值
      console.log(`[${requestId}] ${req.method} ${url.pathname} ${res.statusCode}`); // 访问日志
    }
  };
}

// —— 客户端：用 node:http 自请求，Connection: close 避免 keep-alive 拖住关闭 ——
function request(port, method, path, body) {
  return new Promise((resolve, reject) => {
    const payload = body === undefined ? null : JSON.stringify(body);
    const req = http.request({
      host: '127.0.0.1', port, method, path,
      headers: {
        'Content-Type': 'application/json',
        Connection: 'close',
        ...(payload ? { 'Content-Length': Buffer.byteLength(payload) } : {}),
      },
    }, (res) => {
      const chunks = [];
      res.on('data', (c) => chunks.push(c));
      res.on('end', () => {
        const text = Buffer.concat(chunks).toString('utf8');
        resolve({ status: res.statusCode, body: text ? JSON.parse(text) : null });
      });
    });
    req.on('error', reject);
    if (payload) req.write(payload);
    req.end();
  });
}

const server = http.createServer(createRouter());
server.requestTimeout = 2000;   // 显式超时，保证不会拖过 3 秒

// 兜底看门狗：unref 不阻止正常退出，仅在真卡住时强制收尾
const watchdog = setTimeout(() => {
  console.error('兜底：超过 2.5 秒强制退出');
  server.closeAllConnections();
  process.exit(1);
}, 2500);
watchdog.unref();

server.listen(0, '127.0.0.1', async () => {
  const port = server.address().port;
  console.log('监听随机端口（listen(0) 由系统分配），仅本机自请求');

  // 顺序自请求覆盖：路由 / JSON / 校验 / 错误码（顺序保证建号与输出可复现）
  const checks = [
    ['GET', '/healthz', undefined],
    ['POST', '/users', { name: 'alice' }],
    ['POST', '/users', { name: 'bob' }],
    ['GET', '/users', undefined],
    ['GET', '/users/1', undefined],
    ['GET', '/users/999', undefined],   // 404 NOT_FOUND
    ['POST', '/users', {}],             // 400 VALIDATION
    ['DELETE', '/users', undefined],    // 405 METHOD_NOT_ALLOWED
    ['GET', '/unknown', undefined],     // 404 NOT_FOUND
  ];
  const summary = [];
  for (const [method, path, body] of checks) {
    const r = await request(port, method, path, body);
    summary.push(`${method} ${path} -> ${r.status} ${r.body?.error?.code ?? 'OK'}`);
  }
  console.log('自请求结果：');
  for (const line of summary) console.log('  ', line);

  // 20 课：Promise.allSettled 并发等待多个只读请求（扇入）
  const settled = await Promise.allSettled([
    request(port, 'GET', '/healthz'),
    request(port, 'GET', '/users'),
  ]);
  console.log('并发自请求 allSettled 状态：', settled.map((s) => s.status).join(','));
  console.log('时间特性：performance.now() 记录', durations.length, '条耗时差值，均非负：', durations.every((d) => d >= 0));

  // 自关闭：先断掉可能残留的连接，再 server.close()
  server.closeAllConnections();
  server.close(() => {
    clearTimeout(watchdog);
    console.log('server.close() 完成，进程自行退出');
  });
});
