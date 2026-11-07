import http from 'node:http';
import { randomUUID } from 'node:crypto';

// 与 Go 版同样的最小 CRUD 服务：内存存储 + 路由 + 统一错误 + trace id + 优雅关闭
const store = new Map();
let nextId = 1;

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

async function readJSON(req) {
  const chunks = [];
  for await (const chunk of req) {
    chunks.push(chunk);
    if (Buffer.concat(chunks).length > 1_000_000) {
      throw new Error('body too large');
    }
  }
  return JSON.parse(Buffer.concat(chunks).toString('utf8') || '{}');
}

const server = http.createServer(async (req, res) => {
  const requestId = req.headers['x-request-id'] || randomUUID();
  res.setHeader('X-Request-ID', requestId);
  const start = process.hrtime.bigint();
  const url = new URL(req.url, `http://${req.headers.host || 'localhost'}`);

  try {
    if (url.pathname === '/healthz') return send(res, 200, { status: 'ok' });
    if (url.pathname === '/readyz') return send(res, 200, { status: 'ready' });

    if (url.pathname === '/users') {
      if (req.method === 'GET') {
        return send(res, 200, { data: [...store.values()] });
      }
      if (req.method === 'POST') {
        const input = await readJSON(req);
        if (!input.name || !String(input.name).trim()) {
          return fail(res, 400, 'VALIDATION', 'name is required');
        }
        const user = { id: nextId++, name: String(input.name) };
        store.set(user.id, user);
        return send(res, 201, { data: user });
      }
      res.setHeader('Allow', 'GET, POST');
      return fail(res, 405, 'METHOD_NOT_ALLOWED', 'use GET or POST');
    }

    const match = url.pathname.match(/^\/users\/(\d+)$/);
    if (match) {
      if (req.method !== 'GET') return fail(res, 405, 'METHOD_NOT_ALLOWED', 'use GET');
      const user = store.get(Number(match[1]));
      return user
        ? send(res, 200, { data: user })
        : fail(res, 404, 'NOT_FOUND', 'user not found');
    }

    return fail(res, 404, 'NOT_FOUND', 'route not found');
  } catch (err) {
    console.error(`[${requestId}]`, err);
    return fail(res, 400, 'BAD_REQUEST', err.message);
  } finally {
    const ms = Number(process.hrtime.bigint() - start) / 1e6;
    console.log(`${req.method} ${url.pathname} ${res.statusCode} ${ms.toFixed(1)}ms request_id=${requestId}`);
  }
});

// server 的超时也要显式设置，默认值对生产来说太宽松
server.requestTimeout = 10_000;
server.headersTimeout = 2_000;
server.keepAliveTimeout = 5_000;

let shuttingDown = false;
function shutdown(signal) {
  if (shuttingDown) return;
  shuttingDown = true;
  console.log(`${signal} received, draining`);
  server.close(() => {
    console.log('all connections closed');
    process.exit(0);
  });
  setTimeout(() => {
    console.error('forced shutdown after timeout');
    process.exit(1);
  }, 5000).unref();
}
process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));

const port = Number(process.env.PORT || 8081);
server.listen(port, () => console.log(`listening on http://localhost:${port}`));
