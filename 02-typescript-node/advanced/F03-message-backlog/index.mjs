const BACKLOG = 1000;
const PER_SEC = 10;      // 单消费者每秒处理 10 条
const PRODUCERS = 20;    // 生产端每秒 20 条

console.log('== 1. 单消费者，逐条处理 ==');
console.log(`${BACKLOG} 条 / ${PER_SEC} 条每秒 = ${BACKLOG / PER_SEC} 秒`);

console.log('== 2. 扩到 3 个消费者 ==');
const three = Math.ceil(BACKLOG / (PER_SEC * 3));
console.log(`${BACKLOG} 条 / ${PER_SEC * 3} 条每秒 = ${three} 秒`);

console.log('== 3. 批量 10 条一次（单消费者） ==');
const batch = Math.ceil(BACKLOG / (PER_SEC * 10));
console.log(`${BACKLOG} 条 / ${PER_SEC * 10} 条每秒 = ${batch} 秒`);

console.log('== 4. 生产仍在继续（每秒 20 条） ==');
const net = PER_SEC * 10 - PRODUCERS;
const clear = Math.ceil(BACKLOG / net);
console.log(`净消化 ${net} 条每秒，清空 ${BACKLOG} 条需要 ${clear} 秒`);

console.log('== 5. 结论 ==');
console.log('积压先算「净消化速率」，不是「消费速率」——生产端还在往里灌');
console.log('扩消费者最快但受分区数限制；批量提升单机吞吐，代价是延迟和失败重放');
