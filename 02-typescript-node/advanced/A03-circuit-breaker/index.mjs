// 熔断器：closed -> open -> half-open -> closed
const State = Object.freeze({ CLOSED: 'closed', OPEN: 'open', HALF_OPEN: 'half-open' });

class Breaker {
  constructor({ threshold = 5, cooldownMs = 200, now = () => Date.now() } = {}) {
    this.threshold = threshold;
    this.cooldownMs = cooldownMs;
    this.now = now;
    this.state = State.CLOSED;
    this.failures = 0;
    this.openedAt = 0;
    this.probeInFlight = false;
  }

  allow() {
    if (this.state === State.OPEN) {
      if (this.now() - this.openedAt >= this.cooldownMs) {
        this.state = State.HALF_OPEN;
        this.probeInFlight = false;
        return { ok: true };
      }
      return { ok: false, reason: 'circuit open: fast fail' };
    }
    if (this.state === State.HALF_OPEN) {
      if (this.probeInFlight) {
        return { ok: false, reason: 'circuit half-open: probe already in flight' };
      }
      this.probeInFlight = true;
      return { ok: true };
    }
    return { ok: true };
  }

  success() {
    this.failures = 0;
    this.probeInFlight = false;
    this.state = State.CLOSED;
  }

  failure() {
    this.failures += 1;
    this.probeInFlight = false;
    if (this.state === State.HALF_OPEN || this.failures >= this.threshold) {
      this.state = State.OPEN;
      this.openedAt = this.now();
    }
  }
}

function call(breaker, downstreamOk) {
  const gate = breaker.allow();
  if (!gate.ok) return `rejected: ${gate.reason}`;
  if (downstreamOk) {
    breaker.success();
    return 'ok';
  }
  breaker.failure();
  return 'downstream error';
}

let clock = 0;
const breaker = new Breaker({ threshold: 5, cooldownMs: 200, now: () => clock });

console.log('== 1. 正常请求 ==');
console.log(`call -> ${call(breaker, true)}, state=${breaker.state}`);

console.log('== 2. 连续 5 次失败触发熔断 ==');
for (let i = 1; i <= 5; i += 1) {
  console.log(`failure #${i} -> ${call(breaker, false)}, state=${breaker.state}`);
}

console.log('== 3. 熔断期间快速失败 ==');
console.log(`call -> ${call(breaker, true)}, state=${breaker.state}`);

console.log('== 4. 冷却 200ms 后进入半开 ==');
clock += 200;
console.log(`allow -> ${JSON.stringify(breaker.allow())}, state=${breaker.state}`);

console.log('== 5. 探测成功恢复 closed ==');
breaker.success();
console.log(`state=${breaker.state}, 后续请求 -> ${call(breaker, true)}`);

console.log('== 6. 半开探测失败则重新 open ==');
for (let i = 0; i < 5; i += 1) breaker.failure();
clock += 200;
breaker.allow();
breaker.failure();
console.log(`state=${breaker.state}（探测失败后立即回到 open）`);
