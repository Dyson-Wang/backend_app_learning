import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {
    static final class FakeDb {
        private final AtomicInteger hits = new AtomicInteger();

        String query(String key) {
            hits.incrementAndGet();
            try {
                Thread.sleep(20);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
            return "value-of-" + key;
        }

        int hits() {
            return hits.get();
        }
    }

    static String unprotectedGet(FakeDb db, Map<String, String> cache, String key) {
        String value = cache.get(key);
        if (value != null) {
            return value;
        }
        value = db.query(key);
        cache.put(key, value);
        return value;
    }

    static final class GuardedCache {
        private final Map<String, String> data = new ConcurrentHashMap<>();
        private final Map<String, Object> keyLocks = new ConcurrentHashMap<>();

        String get(FakeDb db, String key) {
            String value = data.get(key);
            if (value != null) {
                return value;
            }
            Object lock = keyLocks.computeIfAbsent(key, k -> new Object());
            synchronized (lock) {
                value = data.get(key);
                if (value != null) {
                    return value;
                }
                value = db.query(key);
                data.put(key, value);
                return value;
            }
        }
    }

    static int run(int n, boolean guarded) throws Exception {
        FakeDb db = new FakeDb();
        ExecutorService pool = Executors.newFixedThreadPool(n);
        CountDownLatch start = new CountDownLatch(1);
        CountDownLatch done = new CountDownLatch(n);

        if (guarded) {
            GuardedCache cache = new GuardedCache();
            for (int i = 0; i < n; i++) {
                pool.submit(() -> {
                    try {
                        start.await();
                        cache.get(db, "user:1");
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    } finally {
                        done.countDown();
                    }
                });
            }
        } else {
            Map<String, String> cache = new ConcurrentHashMap<>();
            for (int i = 0; i < n; i++) {
                pool.submit(() -> {
                    try {
                        start.await();
                        unprotectedGet(db, cache, "user:1");
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    } finally {
                        done.countDown();
                    }
                });
            }
        }

        long t0 = System.nanoTime();
        start.countDown();
        done.await();
        long ms = (System.nanoTime() - t0) / 1_000_000;
        pool.shutdownNow();

        String label = guarded ? "按 key 加锁 + 双重检查" : "无保护";
        System.out.println(label + "：" + n + " 并发 -> 回源 " + db.hits() + " 次，耗时 " + ms + "ms");
        return db.hits();
    }

    public static void main(String[] args) throws Exception {
        int n = 50;
        run(n, false);
        run(n, true);
        System.out.println("结论：缓存击穿是「同一热 key 失效瞬间的并发回源」，线程池越大越明显");
    }
}
