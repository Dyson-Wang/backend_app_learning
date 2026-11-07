// 配置优先级：环境变量 > 配置文件 > 默认值
const resolve = (key, def, file, env) => {
  if (key in env) return { value: env[key], origin: 'env' };
  if (key in file) return { value: file[key], origin: 'file' };
  return { value: def, origin: 'default' };
};

const migrate = (applied, all) => {
  const ran = [];
  for (const m of all) {
    if (applied.has(m.version)) continue;
    applied.add(m.version);
    ran.push(m.version);
  }
  return ran;
};

const file = { 'db.name': 'brm', port: '9090' };
const env = { port: '8080' };

console.log('== 1. 配置优先级 ==');
for (const k of ['port', 'db.host', 'db.name']) {
  const s = resolve(k, 'localhost', file, env);
  console.log(`${k} = ${s.value}（来源：${s.origin}）`);
}

console.log('== 2. 密钥缺失拒绝启动 ==');
const secrets = {};
if (!('DATABASE_URL' in secrets)) {
  console.log('启动失败：DATABASE_URL 未设置（不要给密钥默认值）');
}

console.log('== 3. migration ==');
const all = [
  { version: 1, sql: 'create accounts' },
  { version: 2, sql: 'add email' },
  { version: 3, sql: 'add index' },
];
const applied = new Set();
for (const v of migrate(applied, all)) {
  console.log(`应用 v${v} ${all[v - 1].sql}`);
}
console.log(`当前版本 = ${applied.size}`);
const again = migrate(applied, all);
console.log(`再次执行：${again.length} 个待应用，当前版本 = ${applied.size}`);

console.log('== 4. 结论 ==');
console.log('配置分层要能说出「这个值从哪来」，否则排查时全靠猜');
console.log('migration 必须可重复执行：靠版本表跳过，不靠人记得跑过没有');
