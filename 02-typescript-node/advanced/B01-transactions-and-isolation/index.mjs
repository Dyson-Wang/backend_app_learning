// 内存模拟一行账户
class Account {
  constructor(balance, version) {
    this.balance = balance;
    this.version = version;
  }

  read() {
    return { balance: this.balance, version: this.version };
  }

  // 乐观锁：version 没变才写入
  cas(expectedVersion, newBalance) {
    if (this.version !== expectedVersion) return false;
    this.balance = newBalance;
    this.version += 1;
    return true;
  }
}

function debit(acc, amount, maxAttempts = 10) {
  for (let attempt = 0; attempt <= maxAttempts; attempt += 1) {
    const { balance, version } = acc.read();
    if (acc.cas(version, balance - amount)) return { ok: true, attempt };
  }
  return { ok: false, attempt: maxAttempts + 1 };
}

console.log('== 1. 无锁：两个事务都基于旧值写回 ==');
const acc = new Account(100, 1);
const t1 = acc.read();
const t2 = acc.read();
acc.balance = t1.balance - 10;
acc.version = t1.version + 1;
acc.balance = t2.balance - 10; // 覆盖 T1
acc.version = t2.version + 1;
console.log(`期望 80，实际 ${acc.balance}（丢失更新）`);

console.log('== 2. 乐观锁：version 校验 + 重试 ==');
const acc2 = new Account(100, 1);
const stale = acc2.read(); // T2 先读旧值
const r1 = debit(acc2, 10);
console.log(`T1 ok=${r1.ok} attempts=${r1.attempt}，balance=${acc2.balance} version=${acc2.version}`);

if (acc2.cas(stale.version, stale.balance - 10)) {
  console.log('T2 用旧 version 写入成功（不应发生）');
} else {
  console.log('T2 用旧 version 写入失败：version 已被 T1 改掉');
}

const r2 = debit(acc2, 10);
console.log(`T2 重试 ok=${r2.ok} attempts=${r2.attempt}，balance=${acc2.balance} version=${acc2.version}`);
console.log(`期望 80，实际 ${acc2.balance}`);

console.log('== 3. 结论 ==');
console.log('悲观锁锁整行，吞吐低；乐观锁读不加锁、写校验 version，冲突才重试');
console.log('读已提交 + 乐观锁是多数业务的折中，范围查询要小心幻读');
