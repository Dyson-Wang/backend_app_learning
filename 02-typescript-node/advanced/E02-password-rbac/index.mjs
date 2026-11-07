import { pbkdf2Sync, timingSafeEqual } from 'node:crypto';

const hashPassword = (password, salt) =>
  pbkdf2Sync(password, salt, 1000, 16, 'sha256').toString('hex');

// 固定时间比较，避免按字节提前返回泄露信息
const safeEqual = (a, b) => {
  const x = Buffer.from(a, 'hex');
  const y = Buffer.from(b, 'hex');
  return x.length === y.length && timingSafeEqual(x, y);
};

const rolePerms = {
  admin: ['order:read', 'order:write', 'user:delete'],
  viewer: ['order:read'],
};
const allowed = (role, perm) => (rolePerms[role] ?? []).includes(perm);

const salt = 'salt-8f2a';
const hash = hashPassword('correct-horse', salt);

console.log('== 1. 加盐哈希 ==');
console.log(`salt = ${salt}，hash = ${hash.slice(0, 16)}...`);
console.log(`同一密码 + 同一盐 结果一致 = ${hash === hashPassword('correct-horse', salt)}`);
console.log(`同一密码 + 不同盐 结果不同 = ${hash !== hashPassword('correct-horse', 'salt-other')}`);

console.log('== 2. 校验 ==');
console.log(`正确密码 = ${safeEqual(hash, hashPassword('correct-horse', salt))}`);
console.log(`错误密码 = ${safeEqual(hash, hashPassword('wrong', salt))}`);

console.log('== 3. RBAC ==');
console.log(`admin 能 order:write = ${allowed('admin', 'order:write')}`);
console.log(`viewer 能 order:write = ${allowed('viewer', 'order:write')}`);
console.log(`viewer 能 user:delete = ${allowed('viewer', 'user:delete')}`);

console.log('== 4. 结论 ==');
console.log('密码永远不存明文，也不要可逆加密；盐必须每用户独立随机');
console.log('1000 轮只是演示，生产用 bcrypt/argon2id，成本参数按机器调');
