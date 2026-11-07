// Node 的并发模型：单线程事件循环 + 非阻塞 I/O。
// 这里先讲「任务调度顺序」与 Promise 组合子，真正的线程/消息传递放到 19 课。

async function main() {
  // —— 阶段一：ES 微任务 先于 宏任务（check 阶段）——
  // 用数组收集先后，再一次性打印，输出稳定、不依赖真实耗时
  const order = ['同步'];
  queueMicrotask(() => order.push('queueMicrotask'));
  Promise.resolve().then(() => order.push('Promise.then'));
  setImmediate(() => {
    order.push('setImmediate');
    console.log('【顺序·微任务先于宏任务】', order.join(' -> '));
  });
  await new Promise((resolve) => setImmediate(resolve)); // 屏障：等上面排空

  // —— 阶段二：process.nextTick 是 Node 专属高优先级队列 ——
  // 它会连同回调里嵌套排入的 nextTick 一起，在事件循环前进前被完整排空
  const tickLog = [];
  process.nextTick(() => {
    tickLog.push('第 1 个');
    process.nextTick(() => tickLog.push('嵌套第 2 个'));
  });
  await new Promise((resolve) => setImmediate(resolve));
  console.log('【nextTick】嵌套回调在宏任务前排空：', tickLog.join(' -> '));

  // —— 阶段三：Promise 组合子，一次发起、统一等待 ——
  const task = (id, ok = true) =>
    ok ? Promise.resolve(`任务${id}完成`) : Promise.reject(new Error(`任务${id}失败`));
  const all = await Promise.all([task(1), task(2), task(3)]);
  console.log('【Promise.all】', all.join(' / '));
  const race = await Promise.race([task('快'), task('慢')]);
  console.log('【Promise.race】取最先：', race);
  const settled = await Promise.allSettled([task(4), task(5, false)]);
  console.log('【allSettled】状态：', settled.map((r) => r.status).join(','));

  // —— 阶段四：async/await 只是 Promise 语法糖，错误仍靠 try/catch ——
  async function risky(flag) {
    if (!flag) throw new Error('boom');
    return 42;
  }
  try {
    await risky(false);
  } catch (err) {
    console.log('【async 错误】catch 接住：', err.message);
  }
  console.log('【async 返回值】一定是 Promise：', risky(true) instanceof Promise);

  // —— 阶段五：Promise 执行器同步跑，只有 .then 回调进微任务 ——
  const execOrder = [];
  new Promise((resolve) => { execOrder.push('executor'); resolve(); })
    .then(() => execOrder.push('then'));
  execOrder.push('after');
  await Promise.resolve();
  console.log('【执行器同步】', execOrder.join(' -> '));
  console.log('【并发原语】单线程内靠微任务 + 事件循环交错，CPU 密集同步循环会卡住一切');
}

await main();
