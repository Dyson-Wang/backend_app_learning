import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {
    public static void main(String[] args) throws Exception {
        // 1) 直接创建 Thread + Runnable（lambda）
        var t = new Thread(() -> System.out.println("独立线程: 我跑完了"));
        t.start();
        t.join();                         // 主线程等它结束，保证输出顺序稳定

        int tasks = 8, loops = 1000;

        // 2) AtomicInteger：CAS 保证多线程自增无数据竞争
        var atomic = new AtomicInteger(0);
        var latch = new CountDownLatch(tasks);
        ExecutorService pool = Executors.newFixedThreadPool(4);
        for (int i = 0; i < tasks; i++) {
            pool.submit(() -> {
                for (int j = 0; j < loops; j++) atomic.incrementAndGet();
                latch.countDown();
            });
        }
        latch.await();                    // 等全部任务完成，不用 sleep
        System.out.println("AtomicInteger 期望=" + (tasks * loops) + " 实际=" + atomic.get());
        pool.shutdown();

        // 3) synchronized：用对象监视器保护临界区
        var counter = new SyncCounter();
        var latch2 = new CountDownLatch(tasks);
        ExecutorService pool2 = Executors.newFixedThreadPool(4);
        for (int i = 0; i < tasks; i++) {
            pool2.submit(() -> {
                for (int j = 0; j < loops; j++) counter.inc();
                latch2.countDown();
            });
        }
        latch2.await();
        System.out.println("synchronized 期望=" + (tasks * loops) + " 实际=" + counter.get());
        pool2.shutdown();

        // 4) volatile：保证可见性（不保证 i++ 这类复合操作原子性），这里只演示标志位
        volatileFlag = true;
        System.out.println("volatile 标志位=" + volatileFlag);

        System.out.println("主线程结束，executor 已 shutdown，进程自行退出");
    }

    private static volatile boolean volatileFlag = false;

    // 用 synchronized 方法把 n++ 变成原子临界区
    static class SyncCounter {
        private int n = 0;
        synchronized void inc() { n++; }
        synchronized int get() { return n; }
    }
}
