import process from 'node:process';

console.log('hello, backend roadmap');
console.log('node:    ', process.version);
console.log('platform:', process.platform, process.arch);
console.log('pid:     ', process.pid);
console.log('module:  ', import.meta.url);

// 事件循环的最直观证据：同步代码先跑完，timer 回调最后跑
console.log('1: sync');
setTimeout(() => console.log('3: timer callback (event loop)'), 0);
console.log('2: sync');

// 未捕获异常与 Promise 拒绝会直接终结进程，生产服务必须显式兜底
process.on('unhandledRejection', (reason) => {
  console.error('unhandledRejection:', reason);
});
