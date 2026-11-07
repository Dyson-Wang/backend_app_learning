import java.util.concurrent.CountDownLatch;
import java.util.concurrent.CyclicBarrier;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.ReentrantLock;

public class Main {
    private static int racy = 0;

    public static void main(String[] args) throws Exception {
        System.out.println("== 1. 无锁读改写：两个线程各加一次 ==");
        CyclicBarrier barrier = new CyclicBarrier(2);
        Thread t1 = new Thread(() -> {
            int value = racy;
            try { barrier.await(); } catch (Exception e) { Thread.currentThread().interrupt(); return; }
            racy = value + 1;
        });
        Thread t2 = new Thread(() -> {
            int value = racy;
            try { barrier.await(); } catch (Exception e) { Thread.currentThread().interrupt(); return; }
            racy = value + 1;
        });
        t1.start();
        t2.start();
        t1.join();
        t2.join();
        System.out.println("期望 2，实际 " + racy + "（丢更新）");

        System.out.println("== 2. 加锁：4 线程 × 1000 次 ==");
        final int[] counter = {0};
        final Object lock = new Object();
        Thread[] workers = new Thread[4];
        for (int i = 0; i < 4; i++) {
            workers[i] = new Thread(() -> {
                for (int j = 0; j < 1000; j++) {
                    synchronized (lock) {
                        counter[0]++;
                    }
                }
            });
            workers[i].start();
        }
        for (Thread w : workers) {
            w.join();
        }
        System.out.println("期望 4000，实际 " + counter[0]);

        System.out.println("== 3. 反序加锁导致死锁 ==");
        ReentrantLock a = new ReentrantLock();
        ReentrantLock b = new ReentrantLock();
        CountDownLatch done = new CountDownLatch(1);

        Thread ab = new Thread(() -> {
            a.lock();
            try {
                Thread.sleep(10);
                b.lock();
                try { done.countDown(); } finally { b.unlock(); }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } finally {
                a.unlock();
            }
        });
        Thread ba = new Thread(() -> {
            b.lock();
            try {
                Thread.sleep(10);
                a.lock();
                try { done.countDown(); } finally { a.unlock(); }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } finally {
                b.unlock();
            }
        });
        ab.setDaemon(true);
        ba.setDaemon(true);
        ab.start();
        ba.start();

        if (done.await(100, TimeUnit.MILLISECONDS)) {
            System.out.println("未死锁");
        } else {
            System.out.println("检测到死锁：100ms 内无人完成（各持一把锁等对方）");
        }

        System.out.println("== 4. tryLock 带超时，避免永久阻塞 ==");
        ReentrantLock c = new ReentrantLock();
        ReentrantLock d = new ReentrantLock();
        Thread ab2 = new Thread(() -> {
            c.lock();
            try {
                Thread.sleep(10);
                if (d.tryLock(50, TimeUnit.MILLISECONDS)) {
                    try { System.out.println("  ab2 拿到两把锁"); } finally { d.unlock(); }
                } else {
                    System.out.println("  ab2 放弃 d，避免死锁");
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } finally {
                c.unlock();
            }
        });
        Thread ba2 = new Thread(() -> {
            d.lock();
            try {
                Thread.sleep(10);
                if (c.tryLock(50, TimeUnit.MILLISECONDS)) {
                    try { System.out.println("  ba2 拿到两把锁"); } finally { c.unlock(); }
                } else {
                    System.out.println("  ba2 放弃 c，避免死锁");
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } finally {
                d.unlock();
            }
        });
        ab2.start();
        ba2.start();
        ab2.join(1000);
        ba2.join(1000);

        System.out.println("结论：死锁靠统一加锁顺序或 tryLock(timeout) 避免；daemon 线程不阻止 JVM 退出");
    }
}
