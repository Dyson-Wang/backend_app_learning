import http from 'node:http';

// 一个会慢 120ms 的接口，用来模拟「在途请求」
const server = http.createServer(async (req, res) => {
  res.setHeader('Connection', 'close');   // 不让 keep-alive 干扰后续验证
  if (req.url === '/slow') {
    await new Promise((r) => setTimeout(r, 120));
    res.end('slow done');
    return;
  }
  res.end('ok');
});

server.listen(0, '127.0.0.1', async () => {
  const port = server.address().port;
  const base = `http://127.0.0.1:${port}`;

  console.log('== 1. 发起一个 120ms 的在途请求 ==');
  const inflight = fetch(`${base}/slow`).then((r) => r.text());
  await new Promise((r) => setTimeout(r, 30));   // 确保请求已经进到 handler

  console.log('== 2. 触发优雅关闭 ==');
  const t0 = Date.now();
  const closed = new Promise((resolve) => server.close(resolve));

  console.log('== 3. 在途请求的结果 ==');
  console.log(`body = ${await inflight}，关闭耗时 = ${Date.now() - t0}ms`);

  console.log('== 4. 关闭后新请求 ==');
  try {
    await fetch(`${base}/`);
    console.log('新请求：成功（不应该发生）');
  } catch (err) {
    console.log(`新请求：被拒绝 -> ${err.cause?.code ?? err.message}`);
  }

  await closed;
  console.log('== 5. 结论 ==');
  console.log('server.close() 停止接收新连接，但会等在途请求跑完');
  console.log('生产上还要配超时兜底：server.closeAllConnections() 或进程级 kill 定时器');
});
