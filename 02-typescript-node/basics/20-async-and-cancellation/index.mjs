import { AsyncLocalStorage } from 'node:async_hooks';

// ============ 1. 异步任务模型 + 上下文传值（AsyncLocalStorage）============
// 把 requestId 放进「异步上下文」，无需一层层手动传参，跨 await / 微任务 / 宏任务都在
const als = new AsyncLocalStorage();
console.log('【上下文】run 之前：', als.getStore());
await als.run({ requestId: 'req-42' }, async () => {
  console.log('【上下文】run 内：', als.getStore().requestId);
  await Promise.resolve();
  console.log('【上下文】await 微任务后：', als.getStore().requestId);
  await new Promise((resolve) => setImmediate(resolve));
  console.log('【上下文】setImmediate 宏任务后：', als.getStore().requestId);
});
console.log('【上下文】run 结束后恢复：', als.getStore());

// ============ 2. 可取消任务：共享 signal，取消时清理资源 ============
function cancellable(name, signal, value) {
  return new Promise((resolve, reject) => {
    if (signal.aborted) { reject(signal.reason); return; }
    const timer = setTimeout(() => resolve(`${name}=${value}`), 500);
    signal.addEventListener('abort', () => {
      clearTimeout(timer);                       // 取消后清理：释放定时器等资源
      console.log(`  [清理] ${name} 的定时器已清除`);
      reject(signal.reason);
    }, { once: true });
  });
}

// —— 取消传播：一个 controller 的 signal 传给多个子任务 ——
const controller = new AbortController();
const parent = controller.signal;
const extra = new AbortController();
const combined = AbortSignal.any([parent, extra.signal]); // 多个取消源任一触发即取消
parent.addEventListener('abort', () => {
  console.log('  [abort 事件] 父信号已取消，reason =', parent.reason.message);
}, { once: true });

const tasks = [
  cancellable('任务1', parent, 1),
  cancellable('任务2', parent, 2),
  cancellable('任务3', combined, 3),
];
controller.abort(new Error('主动取消'));   // 同步传播给 任务1 / 任务2 / 任务3

// —— 并发等待：allSettled 收齐每个任务的结局（不因某个 reject 而短路）——
const settled = await Promise.allSettled(tasks);
console.log('【并发等待】allSettled 状态：', settled.map((r) => r.status).join(','));
console.log('【取消原因】', settled.map((r) => r.reason.message).join(' / '));

// —— 超时：AbortSignal.timeout 到点自动 abort（内部定时器不阻止进程退出）——
const timeoutSignal = AbortSignal.timeout(30);
console.log('【超时】timeout(30) 初始 aborted：', timeoutSignal.aborted);
try {
  await cancellable('超时任务', timeoutSignal, 'ok');
} catch (err) {
  console.log('【超时】任务被取消：', err.name, '；signal.aborted：', timeoutSignal.aborted);
}

// 未被取消的 signal 保持可用；正常任务靠自身逻辑完成，不依赖 signal
const idle = new AbortController();
console.log('【未取消】idle.signal.aborted：', idle.signal.aborted, '；可继续用于后续任务');
console.log('结论：AbortController 管取消、AbortSignal.timeout/any 管超时与组合、AsyncLocalStorage 管上下文');
