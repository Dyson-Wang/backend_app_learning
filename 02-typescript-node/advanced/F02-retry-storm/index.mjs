const CLIENTS = 100;
const DOWNSTREAM_CAP = 100;

console.log('== 1. 无退避：失败立即重试 3 次 ==');
const noBackoff = CLIENTS * (1 + 3);
console.log(`下游收到 ${noBackoff} QPS（容量 ${DOWNSTREAM_CAP}，超载 ${noBackoff / DOWNSTREAM_CAP} 倍）`);

console.log('== 2. 指数退避 + 全抖动 ==');
// 3 次重试分散到 1s / 2s / 4s，每个窗口只承接 1/3 的客户端
const peak = CLIENTS + Math.floor(CLIENTS / 3);
console.log(`下游峰值 ≈ ${peak} QPS（接近容量，不再放大）`);

console.log('== 3. 再加熔断 ==');
console.log('下游请求 = 0 QPS（open 状态本地快速失败，直接返回降级结果）');

console.log('== 4. 结论 ==');
console.log('重试必须有上限、退避和抖动；三个都没有，就是给故障加杠杆');
console.log('熔断才是止损的关键：下游已经挂了，再发请求只会让它更晚恢复');
