// 有界队列 + 满则拒绝。消费者先被 gate 挡住，保证提交阶段队列必满。
const capacity = 5;
const total = 20;

const queue = [];
let accepted = 0;
let rejected = 0;
let processed = 0;
let release;
const gate = new Promise((resolve) => { release = resolve; });

async function consumer() {
  await gate;
  while (queue.length > 0) {
    queue.shift();
    processed += 1;
  }
}

const consuming = consumer();

for (let i = 0; i < total; i += 1) {
  if (queue.length < capacity) {
    queue.push(i);
    accepted += 1;
  } else {
    rejected += 1; // 满则快速拒绝
  }
}

release();
await consuming;

console.log(`容量=${capacity}，提交=${total}`);
console.log(`accepted=${accepted} rejected=${rejected} processed=${processed}`);
console.log('结论：拒绝发生在入口，内存和延迟都是常数；无界数组只是把 OOM 推迟');

const unbounded = Array.from({ length: total }, (_, i) => i);
console.log(`无界队列长度=${unbounded.length}（全部接收，没有上限）`);
