import { Worker, isMainThread } from 'node:worker_threads';
import { EventEmitter } from 'node:events';

// ============================================================
// 第一部分：worker_threads —— 真线程、消息传递、共享状态保护
// ============================================================
// 用 data: URL 内联 worker 源码，避免多文件；worker 内用 ESM import node: 内置模块。
const workerCode = `
  import { parentPort, workerData } from 'node:worker_threads';
  const { id, numbers, shared } = workerData;
  const view = new Int32Array(shared);          // 多线程共享的一块内存
  let partial = 0;
  for (const n of numbers) {
    partial += n;
    Atomics.add(view, 0, n);                   // 原子累加，无需互斥锁
  }
  parentPort.postMessage({ id, partial });     // 把结果通过消息传回主线程
`;
const workerUrl = 'data:text/javascript,' + encodeURIComponent(workerCode);
const shared = new SharedArrayBuffer(4);       // 1 个 int32
const total = new Int32Array(shared);

// 一个 worker 跑一片数据（扇出），主线程用 Promise.all 汇总（扇入）
function runWorker(id, numbers) {
  return new Promise((resolve, reject) => {
    const w = new Worker(new URL(workerUrl), { workerData: { id, numbers, shared } });
    let message;
    w.on('message', (m) => { message = m; });
    w.on('error', reject);
    w.on('exit', (code) => resolve({ ...message, exitCode: code })); // 等 worker 真正退出
  });
}

const parts = await Promise.all([
  runWorker('A', [1, 2, 3]),
  runWorker('B', [4, 5, 6]),
]);
console.log('主线程 isMainThread：', isMainThread);
console.log('扇入各 worker 部分和：', parts.map((p) => `${p.id}=${p.partial}`).join(' / '));
console.log('worker 退出码：', parts.map((p) => p.exitCode).join(','));
console.log('Atomics.load 共享总和：', Atomics.load(total, 0), '（跨线程累加无数据竞争）');

// ============================================================
// 第二部分：EventEmitter + 有界异步队列 —— 生产消费 / 背压 / 关闭信号
// ============================================================
class BoundedQueue {
  #items = [];
  #takers = [];       // 队列空时等待数据的消费者 resolve
  #putters = [];      // 队列满时等待空位的生产者 { item, resolve }
  #closed = false;
  #events = new EventEmitter();
  constructor(capacity) { this.capacity = capacity; }

  push(item) {
    if (this.#closed) return Promise.reject(new Error('队列已关闭'));
    if (this.#takers.length) {                 // 有消费者在等，直接交付
      this.#takers.shift()(item);
      return Promise.resolve();
    }
    if (this.#items.length < this.capacity) {
      this.#items.push(item);
      this.#events.emit('push', this.#items.length);
      return Promise.resolve();
    }
    // 满：返回 pending Promise，生产者 await 形成背压（等价于「阻塞」）
    return new Promise((resolve) => this.#putters.push({ item, resolve }));
  }

  take() {
    if (this.#items.length) {
      const item = this.#items.shift();
      if (this.#putters.length) {             // 空出位置，放行一个生产者
        const { item: pending, resolve } = this.#putters.shift();
        this.#items.push(pending);
        resolve();
      }
      this.#events.emit('take', this.#items.length);
      return Promise.resolve(item);
    }
    if (this.#closed) return Promise.resolve(undefined); // 关闭信号
    return new Promise((resolve) => this.#takers.push(resolve));
  }

  close() {
    this.#closed = true;
    for (const resolve of this.#takers.splice(0)) resolve(undefined); // 唤醒全部等待者
    this.#events.emit('close');
  }

  on(evt, fn) { return this.#events.on(evt, fn); }
}

// —— 有界队列与背压：容量 1，塞满后第二个 push 变成 pending Promise ——
const q = new BoundedQueue(1);
await q.push('任务A');
const blocked = q.push('任务B');
console.log('队列满时 push 返回 pending：', blocked instanceof Promise, '；有界容量：', q.capacity);
console.log('取走一个释放空位：', await q.take());
await blocked; // 生产者被放行、背压解除
console.log('背压解除后取到第二个：', await q.take());
q.close();

// —— 生产消费 + 关闭信号传播 + 等待多个任务 ——
const queue = new BoundedQueue(2);
queue.on('close', () => console.log('  [队列] 收到关闭信号，等待中的消费者被唤醒'));

async function produce() {
  for (let i = 1; i <= 6; i++) {
    await queue.push({ id: i });
  }
  queue.close();   // 关闭信号沿队列传播给所有消费者
}

async function consume(name, sink) {
  for (;;) {
    const item = await queue.take();
    if (item === undefined) return `${name} 退出`;  // 收到关闭信号
    sink.push(item.id);
    await Promise.resolve();  // 模拟处理：纯微任务，不 sleep
  }
}

const consumed = [];
const [, ...exits] = await Promise.all([
  produce(),                                     // 生产者
  consume('消费者1', consumed), consume('消费者2', consumed), // 扇出到 2 个消费者
]);
console.log('扇入汇总：共消费', consumed.length, '项，去重 id：', [...new Set(consumed)].sort((a, b) => a - b).join(','));
console.log('消费者退出：', exits.join(' / '));
console.log('单线程里 Map/数组 push 无需锁；跨线程共享状态才需要 Atomics（见第一部分）');
