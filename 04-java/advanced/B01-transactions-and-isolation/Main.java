public class Main {

    static class Account {
        private int balance;
        private long version;

        Account(int balance) {
            this.balance = balance;
        }

        synchronized long version() {
            return version;
        }

        synchronized int balance() {
            return balance;
        }

        // 乐观锁：只有 version 没被别人改过才写入
        synchronized boolean cas(long expectedVersion, int newBalance) {
            if (version != expectedVersion) {
                return false;
            }
            balance = newBalance;
            version++;
            return true;
        }
    }

    static int debit(Account acc, int amount) {
        for (int attempt = 0; attempt <= 10; attempt++) {
            long v = acc.version();
            int b = acc.balance();
            if (acc.cas(v, b - amount)) {
                return attempt;
            }
        }
        return -1;
    }

    public static void main(String[] args) {
        System.out.println("== 1. 无锁：两个事务都基于旧值写回 ==");
        Account naive = new Account(100);
        int read1 = naive.balance();
        int read2 = naive.balance();
        // 两个事务都读到 100，各自算出 90 写回
        naive.balance = read1 - 10;
        naive.balance = read2 - 10;
        System.out.println("期望 80，实际 " + naive.balance + "（丢失更新）");

        System.out.println("== 2. 乐观锁：version 校验 + 重试 ==");
        Account acc = new Account(100);
        int a1 = debit(acc, 10);
        System.out.println("T1 attempts=" + a1 + "，balance=" + acc.balance() + " version=" + acc.version());

        long stale = acc.version() - 1;
        System.out.println("T2 用旧 version=" + stale + " 写入：" + acc.cas(stale, 70));

        int a2 = debit(acc, 10);
        System.out.println("T2 重试 attempts=" + a2 + "，balance=" + acc.balance() + " version=" + acc.version());
        System.out.println("期望 80，实际 " + acc.balance());

        System.out.println("== 3. 结论 ==");
        System.out.println("悲观锁锁整行，吞吐低；乐观锁读不加锁、写校验 version，冲突才重试");
        System.out.println("JPA 里对应 @Version，冲突抛 OptimisticLockException");
    }
}
