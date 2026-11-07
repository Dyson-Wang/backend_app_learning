import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {

    static class Order {
        final int id;
        final String key;
        final int amount;

        Order(int id, String key, int amount) {
            this.id = id;
            this.key = key;
            this.amount = amount;
        }
    }

    static class IdemStore {
        private final ConcurrentHashMap<String, Order> done = new ConcurrentHashMap<>();
        private final ConcurrentHashMap<String, CountDownLatch> inflight = new ConcurrentHashMap<>();
        private final AtomicInteger seq = new AtomicInteger();
        private final AtomicInteger created = new AtomicInteger();
        private final AtomicInteger replayed = new AtomicInteger();

        Order create(String key, int amount) throws InterruptedException {
            Order hit = done.get(key);
            if (hit != null) {
                replayed.incrementAndGet();
                return hit;
            }
            CountDownLatch mine = new CountDownLatch(1);
            CountDownLatch prev = inflight.putIfAbsent(key, mine);
            if (prev != null) {
                prev.await();                     // 等第一个请求跑完
                replayed.incrementAndGet();
                return done.get(key);
            }
            try {
                Thread.sleep(20);                 // 模拟写库
                Order o = new Order(seq.incrementAndGet(), key, amount);
                done.put(key, o);
                created.incrementAndGet();
                return o;
            } finally {
                inflight.remove(key);
                mine.countDown();
            }
        }
    }

    public static void main(String[] args) throws Exception {
        IdemStore store = new IdemStore();
        int n = 20;
        Order[] results = new Order[n];
        Thread[] ts = new Thread[n];

        System.out.println("== 1. 20 个并发请求，同一个 Idempotency-Key ==");
        for (int i = 0; i < n; i++) {
            final int idx = i;
            ts[i] = new Thread(() -> {
                try {
                    results[idx] = store.create("idem-7f3a", 9900);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
            ts[i].start();
        }
        for (Thread t : ts) {
            t.join();
        }

        java.util.Set<Integer> ids = new java.util.HashSet<>();
        for (Order o : results) {
            ids.add(o.id);
        }
        System.out.println("真实建单 = " + store.created.get()
                + "，回放 = " + store.replayed.get()
                + "，订单 id 去重后 = " + ids.size());

        System.out.println("== 2. 换一个 key ==");
        Order o2 = store.create("idem-other", 100);
        System.out.println("新 key 建单 id = " + o2.id);

        System.out.println("== 3. 结论 ==");
        System.out.println("幂等键必须落库并加唯一索引；ConcurrentHashMap 只是演示");
        System.out.println("putIfAbsent 做占位，比「先查再写」少一个竞态窗口");
    }
}
