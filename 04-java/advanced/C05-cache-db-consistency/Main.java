import java.util.HashMap;
import java.util.Map;

public class Main {

    static class Cache {
        private final Map<String, String> data = new HashMap<>();

        String get(String k) {
            return data.get(k);
        }

        void set(String k, String v) {
            data.put(k, v);
        }

        void delete(String k) {
            data.remove(k);
        }
    }

    static String readThrough(Cache cache, Map<String, String> db, String k) {
        String hit = cache.get(k);
        if (hit != null) {
            return hit;
        }
        String v = db.get(k);          // 缓存未命中，回源
        cache.set(k, v);
        return v;
    }

    public static void main(String[] args) {
        String key = "user:1";

        System.out.println("== 1. 先删缓存，再写库 ==");
        Map<String, String> db = new HashMap<>();
        db.put(key, "v1");
        Cache cache = new Cache();
        cache.set(key, "v1");
        cache.delete(key);                  // 写请求：先删缓存
        readThrough(cache, db, key);        // 读请求：回源读到旧值
        db.put(key, "v2");                  // 写请求：再写库
        System.out.println("缓存 = " + cache.get(key) + "，库 = " + db.get(key)
                + "，是否脏 = " + (!cache.get(key).equals(db.get(key))));

        System.out.println("== 2. 先写库，再删缓存 ==");
        db = new HashMap<>();
        db.put(key, "v1");
        cache = new Cache();
        cache.set(key, "v1");
        String old = cache.get(key);        // 读请求先拿到旧值
        db.put(key, "v2");                  // 写请求：先写库
        cache.delete(key);                  // 再删缓存
        cache.set(key, old);                // 读请求迟到，把旧值写回
        System.out.println("缓存 = " + cache.get(key) + "，库 = " + db.get(key)
                + "，是否脏 = " + (!cache.get(key).equals(db.get(key))));

        System.out.println("== 3. 延迟双删 ==");
        db = new HashMap<>();
        db.put(key, "v1");
        cache = new Cache();
        cache.set(key, "v1");
        old = cache.get(key);
        cache.delete(key);
        db.put(key, "v2");
        cache.set(key, old);                // 迟到的读请求仍回填旧值
        cache.delete(key);                  // 延迟再删一次
        System.out.println("缓存命中 = " + (cache.get(key) != null)
                + "，下一次读会回源拿到 " + readThrough(cache, db, key));

        System.out.println("== 4. 结论 ==");
        System.out.println("先删缓存有窗口，先写库再删缓存也有窗口，双删只能把窗口压小");
        System.out.println("要强一致就别用缓存读；要最终一致就配 binlog 订阅 + 重试删除");
    }
}
