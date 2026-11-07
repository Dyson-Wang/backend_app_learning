import { createHmac } from 'node:crypto';

const SECRET = 's3cr3t';
const sign = (payload) =>
  createHmac('sha256', SECRET).update(`HS256.${payload}`).digest('hex').slice(0, 32);

const issue = (sub, jti, exp) => {
  const payload = `sub=${sub};jti=${jti};exp=${exp}`;
  return `HS256.${payload}.${sign(payload)}`;
};

const field = (payload, key) =>
  payload.split(';').find((kv) => kv.startsWith(`${key}=`))?.slice(key.length + 1) ?? '';

function verify(token, now, revoked) {
  const parts = token.split('.');
  if (parts.length !== 3 || parts[0] !== 'HS256') return ['', 'bad format or alg'];
  if (parts[2] !== sign(parts[1])) return ['', 'signature mismatch'];
  if (now >= Number(field(parts[1], 'exp'))) return ['', 'token expired'];
  if (revoked.has(field(parts[1], 'jti'))) return ['', 'token revoked'];
  return [field(parts[1], 'sub'), ''];
}

const now = 1000;
const revoked = new Set();

console.log('== 1. 签发 ==');
const access = issue('u1', 'jti-a1', now + 60);
const refresh = issue('u1', 'jti-r1', now + 86400);
console.log('access.jti = jti-a1，exp = +60s');
console.log('refresh.jti = jti-r1，exp = +86400s');

console.log('== 2. 校验正常 token ==');
let [sub, err] = verify(access, now, revoked);
console.log(`ok = ${err === ''}，sub = ${sub}`);

console.log('== 3. 篡改 payload ==');
[sub, err] = verify(access.replace('sub=u1', 'sub=u2'), now, revoked);
console.log(`ok = ${err === ''}，reason = ${err}`);

console.log('== 4. 刷新：轮换 refresh token ==');
revoked.add('jti-r1');
issue('u1', 'jti-r2', now + 86400);
console.log('新 refresh.jti = jti-r2，旧 jti-r1 进黑名单');
[sub, err] = verify(refresh, now, revoked);
console.log(`用旧 refresh 再刷：ok = ${err === ''}，reason = ${err}`);

console.log('== 5. 吊销 access token ==');
revoked.add('jti-a1');
[sub, err] = verify(access, now, revoked);
console.log(`ok = ${err === ''}，reason = ${err}`);

console.log('== 6. 结论 ==');
console.log('签名只保证没被篡改，吊销必须靠 jti 黑名单（或短 TTL + 版本号）');
console.log('refresh 必须轮换：旧的用过一次就作废，否则泄露后可长期续命');
