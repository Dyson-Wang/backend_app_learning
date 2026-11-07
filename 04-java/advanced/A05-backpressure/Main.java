import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

public class Main {
    public static void main(String[] args) throws Exception {
        final int capacity = 5;
        final int total = 20;

        ArrayBlockingQueue<Integer> queue = new ArrayBlockingQueue<>(capacity);
        CountDownLatch gate = new CountDownLatch(1);
        CountDownLatch done = new CountDownLatch(1);
        final int[] processed = {0};

        Thread consumer = new Thread(() -> {
            try {
                gate.await();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                return;
            }
            while (true) {
                Integer v = queue.poll();
                if (v == null) {
                    break;
                }
                processed[0]++;
            }
            done.countDown();
        });
        consumer.start();

        int accepted = 0;
        int rejected = 0;
        for (int i = 0; i < total; i++) {
            if (queue.offer(i)) {
                accepted++;
            } else {
                rejected++; // 满则快速拒绝
            }
        }

        gate.countDown();
        done.await(2, TimeUnit.SECONDS);
        consumer.join();

        System.out.println("容量=" + capacity + "，提交=" + total);
        System.out.println("accepted=" + accepted + " rejected=" + rejected
                + " processed=" + processed[0]);
        System.out.println("结论：拒绝发生在入口，内存和延迟都是常数；无界队列只是把 OOM 推迟");
        System.out.println("无界队列长度=" + total + "（全部接收，没有上限）");
    }
}
