// 令牌桶：burst 个突发额度，按 rate/s 匀速补充
class TokenBucket {
  constructor({ rate, burst, now = () => 0 }) {
    this.rate = rate;
    this.burst = burst;
    this.tokens = burst;
    this.last = now();
    this.now = now;
  }

  allow() {
    const now = this.now();
    const elapsedSec = (now - this.last) / 1000;
    this.tokens = Math.min(this.burst, this.tokens + elapsedSec * this.rate);
    this.last = now;
    if (this.tokens >= 1) {
      this.tokens -= 1;
      return true;
    }
    return false;
  }
}

let clock = 0;
const bucket = new TokenBucket({ rate: 2, burst: 5, now: () => clock });

function hit(label) {
  const ok = bucket.allow();
  console.log(`  ${label} -> allow=${ok}, tokens=${bucket.tokens.toFixed(2)}`);
  return ok;
}

console.log('== t=0ms：burst=5，连续 6 次 ==');
for (let i = 1; i <= 6; i += 1) hit(`#${i}`);

console.log('== t=500ms：补 1 个令牌 ==');
clock = 500;
hit('#7');
hit('#8');

console.log('== t=1500ms：补 2 个令牌 ==');
clock = 1500;
hit('#9');
hit('#10');
hit('#11');

console.log('结论：令牌桶允许突发（burst），但长期速率被 rate 卡死');
