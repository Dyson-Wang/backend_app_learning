import java.util.concurrent.Semaphore;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {

    static class Pool {
        private final Semaphore sem;
        private final AtomicInteger peak = new AtomicInteger();
        private final AtomicInteger inUse = new AtomicInteger();

        Pool(int size) {
            this.sem = new Semaphore(size);
        }

        boolean acquire(Long waitMs) throws InterruptedException {
            boolean ok = waitMs == null
                    ? acquireBlocking()
                    : sem.tryAcquire(waitMs, TimeUnit.MILLISECONDS);
            if (ok) {
                int now = inUse.incrementAndGet();
                peak.accumulateAndGet(now, Math::max);
            }
            return ok;
        }

        private boolean acquireBlocking() throws InterruptedException {
            sem.acquire();
            return true;
        }

        void release() {
            inUse.decrementAndGet();
            sem.release();
        }
    }

    static int[] run(Pool pool, int n, long holdMs, Long waitMs) throws Exception {
        AtomicInteger ok = new AtomicInteger();
        AtomicInteger timeout = new AtomicInteger();
        Thread[] ts = new Thread[n];
        long t0 = System.nanoTime();
        for (int i = 0; i < n; i++) {
            ts[i] = new Thread(() -> {
                try {
                    if (pool.acquire(waitMs)) {
                        try {
                            Thread.sleep(holdMs);
                            ok.incrementAndGet();
                        } finally {
                            pool.release();
                        }
                    } else {
                        timeout.incrementAndGet();
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
            ts[i].start();
        }
        for (Thread t : ts) {
            t.join();
        }
        int cost = (int) ((System.nanoTime() - t0) / 1_000_000L);
        return new int[]{ok.get(), timeout.get(), cost};
    }

    public static void main(String[] args) throws Exception {
        System.out.println("== 1. 池大小 3，10 个请求，无等待超时 ==");
        Pool p1 = new Pool(3);
        int[] r1 = run(p1, 10, 120, null);
        System.out.println("ok=" + r1[0] + " timeout=" + r1[1] + " 耗时=" + r1[2] + "ms peak=" + p1.peak.get());

        System.out.println("== 2. 同样的池，等待超过 40ms 就快速失败 ==");
        Pool p2 = new Pool(3);
        int[] r2 = run(p2, 10, 120, 40L);
        System.out.println("ok=" + r2[0] + " timeout=" + r2[1] + " 耗时=" + r2[2] + "ms peak=" + p2.peak.get());

        System.out.println("== 3. 结论 ==");
        System.out.println("连接池/并发上限必须配等待超时：没有超时，故障时请求无限排队");
        System.out.println("HikariCP 的 connectionTimeout、Semaphore.tryAcquire 都是这个作用");
    }
}
