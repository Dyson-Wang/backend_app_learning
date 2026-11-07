// 用 Map 模拟 Redis
class FakeRedis {
  constructor(now = () => 0) {
    this.data = new Map();
    this.exp = new Map();
    this.now = now;
  }

  alive(key) {
    const t = this.exp.get(key);
    return t !== undefined && this.now() < t;
  }

  setNX(key, value, ttlMs) {
    if (this.alive(key)) return false;
    this.data.set(key, value);
    this.exp.set(key, this.now() + ttlMs);
    return true;
  }

  get(key) {
    return this.alive(key) ? this.data.get(key) : '';
  }

  // 错误解锁：不校验 token
  del(key) {
    this.data.delete(key);
    this.exp.delete(key);
  }

  // 正确解锁：比较 token 再删（等价 Lua）
  delIfMatch(key, token) {
    if (!this.alive(key) || this.data.get(key) !== token) return false;
    this.data.delete(key);
    this.exp.delete(key);
    return true;
  }
}

let clock = 0;
const redis = new FakeRedis(() => clock);
const key = 'lock:order:1';

console.log('== 1. A 抢到锁 ==');
console.log(`A SetNX = ${redis.setNX(key, 'token-A', 100)}, holder=${redis.get(key)}`);

console.log('== 2. B 抢锁失败 ==');
console.log(`B SetNX = ${redis.setNX(key, 'token-B', 100)}, holder=${redis.get(key)}`);

console.log('== 3. A 的业务超时，锁过期 ==');
clock += 150;
console.log(`锁是否还有效 = ${redis.get(key) !== ''}`);

console.log('== 4. B 抢到锁 ==');
console.log(`B SetNX = ${redis.setNX(key, 'token-B', 100)}, holder=${redis.get(key)}`);

console.log('== 5. A 醒来直接 DEL：删掉了 B 的锁（错误） ==');
redis.del(key);
console.log(`holder="${redis.get(key)}"（B 的锁被误删）`);

console.log('== 6. 修正：delIfMatch 校验 token ==');
redis.setNX(key, 'token-B', 100);
console.log(`A delIfMatch(token-A) = ${redis.delIfMatch(key, 'token-A')}, holder=${redis.get(key)}`);
console.log(`B delIfMatch(token-B) = ${redis.delIfMatch(key, 'token-B')}, holder="${redis.get(key)}"`);

console.log('== 7. 结论 ==');
console.log('加锁 SET key token NX PX ttl；解锁必须比较 token 再删，且要原子');
console.log('业务可能超过 ttl 时要续期，否则锁提前释放');
