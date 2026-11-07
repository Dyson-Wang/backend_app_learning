import java.util.HashMap;
import java.util.Map;

public class Main {

    static class Entry {
        final String token;
        final long expireAtNanos;

        Entry(String token, long expireAtNanos) {
            this.token = token;
            this.expireAtNanos = expireAtNanos;
        }
    }

    // 用 HashMap 模拟 Redis 的 SET NX PX / GET / DEL
    static class FakeRedis {
        private final Map<String, Entry> data = new HashMap<>();

        synchronized boolean setNxPx(String key, String token, long ttlMs) {
            Entry cur = data.get(key);
            long now = System.nanoTime();
            if (cur != null && cur.expireAtNanos > now) {
                return false;
            }
            data.put(key, new Entry(token, now + ttlMs * 1_000_000L));
            return true;
        }

        synchronized String get(String key) {
            Entry cur = data.get(key);
            if (cur == null || cur.expireAtNanos <= System.nanoTime()) {
                return "";
            }
            return cur.token;
        }

        synchronized void delete(String key) {
            data.remove(key);
        }

        // 真实环境用 Lua 保证「比较 + 删除」原子
        synchronized boolean delIfMatch(String key, String token) {
            if (get(key).equals(token)) {
                data.remove(key);
                return true;
            }
            return false;
        }
    }

    public static void main(String[] args) throws Exception {
        FakeRedis r = new FakeRedis();
        String key = "lock:order:42";

        System.out.println("== 1. A 加锁 ==");
        System.out.println("A SetNX = " + r.setNxPx(key, "token-A", 50) + ", holder=" + r.get(key));

        System.out.println("== 2. B 尝试加锁（A 还持有时） ==");
        System.out.println("B SetNX = " + r.setNxPx(key, "token-B", 50));

        System.out.println("== 3. A 的业务超时，锁过期 ==");
        Thread.sleep(60);
        System.out.println("锁是否还有效 = " + (!r.get(key).isEmpty()));

        System.out.println("== 4. B 抢到锁 ==");
        System.out.println("B SetNX = " + r.setNxPx(key, "token-B", 50) + ", holder=" + r.get(key));

        System.out.println("== 5. A 醒来直接 DEL：删掉了 B 的锁（错误） ==");
        r.delete(key);
        System.out.println("holder=\"" + r.get(key) + "\"（B 的锁被误删）");

        System.out.println("== 6. 修正：delIfMatch 校验 token ==");
        r.setNxPx(key, "token-B", 50);
        System.out.println("A delIfMatch(token-A) = " + r.delIfMatch(key, "token-A")
                + ", holder=" + r.get(key));
        System.out.println("B delIfMatch(token-B) = " + r.delIfMatch(key, "token-B")
                + ", holder=" + r.get(key));

        System.out.println("== 7. 结论 ==");
        System.out.println("加锁 SET key token NX PX ttl；解锁必须比较 token 再删，且要原子");
        System.out.println("业务可能超过 ttl 时要续期，Redisson 的 watchdog 就是干这个");
    }
}
