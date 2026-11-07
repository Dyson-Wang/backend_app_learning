// 用逻辑时钟代替真实 sleep，保证输出可复现
const CHUNKS = 5;
const PER_CHUNK = 20;

console.log('== 1. 缓冲模式：全部生成完再发 ==');
{
  let t = 0;
  for (let i = 0; i < CHUNKS; i++) t += PER_CHUNK;
  console.log(`首字节 = ${t}ms，总耗时 = ${t}ms`);
}

console.log('== 2. 流式模式：生成一块发一块 ==');
{
  let t = 0;
  let firstByte = 0;
  for (let i = 0; i < CHUNKS; i++) {
    t += PER_CHUNK;
    if (i === 0) firstByte = t;
  }
  console.log(`首字节 = ${firstByte}ms，总耗时 = ${t}ms`);
}

console.log('== 3. 对比 ==');
console.log('总耗时一样，首字节差了 4 倍：用户感知的是首字节，不是总耗时');
console.log('SSE / chunked / WebSocket 的价值就在这里；代价是不能压缩、不能算 Content-Length');

console.log('== 4. 结论 ==');
console.log('大响应先想清楚：客户端要等全量，还是能边收边渲染');
console.log('流式要处理背压：res.write() 返回 false 就该 await drain');
