import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.Callable;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {
    // 用特殊值作为「关闭信号」在生产消费之间传播
    static final String POISON = "__STOP__";

    public static void main(String[] args) throws Exception {
        // 1) 有界阻塞队列 + 生产消费 + 毒丸关闭信号
        BlockingQueue<String> queue = new ArrayBlockingQueue<>(4);
        ExecutorService pool = Executors.newFixedThreadPool(4);
        List<String> consumed = Collections.synchronizedList(new ArrayList<>());

        Future<?> producer = pool.submit(() -> {
            for (int i = 0; i < 6; i++) {
                try { queue.put("任务" + i); }
                catch (InterruptedException e) { Thread.currentThread().interrupt(); return; }
            }
            try { queue.put(POISON); }       // 关闭信号放进队列传播给消费者
            catch (InterruptedException e) { Thread.currentThread().interrupt(); }
        });

        Future<?> consumer = pool.submit(() -> {
            try {
                for (int i = 0; i < 32; i++) {          // 显式上限，杜绝死循环
                    String item = queue.poll(1, TimeUnit.SECONDS);
                    if (item == null) break;            // 超时兜底
                    if (POISON.equals(item)) break;     // 关闭信号
                    consumed.add(item);
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        });

        producer.get();   // 等待多个任务：先等生产者
        consumer.get();   // 再等消费者
        System.out.println("有界队列(容量4)生产消费: 收集任务数=" + consumed.size()
                + " 队列剩余=" + queue.size());
        System.out.println("关闭信号(毒丸)传播后消费者退出: 首项=" + consumed.get(0));

        // 2) 扇出扇入：提交多个任务，用 Future 汇总结果
        List<Callable<Integer>> jobs = new ArrayList<>();
        for (int i = 1; i <= 4; i++) {
            final int n = i;
            jobs.add(() -> n * n);
        }
        int fanSum = 0;
        for (Future<Integer> f : pool.invokeAll(jobs)) fanSum += f.get();
        System.out.println("扇出扇入(invokeAll) 平方和=" + fanSum);

        // 3) 等待多个任务 + 共享状态保护
        CountDownLatch latch = new CountDownLatch(4);
        AtomicInteger shared = new AtomicInteger(0);
        for (int i = 0; i < 4; i++) {
            pool.submit(() -> {
                shared.incrementAndGet();
                latch.countDown();
            });
        }
        latch.await();
        System.out.println("等待 CountDownLatch 归零: 共享计数=" + shared.get());

        pool.shutdown();
        System.out.println("关闭信号已传播，executor 已 shutdown，进程退出");
    }
}
