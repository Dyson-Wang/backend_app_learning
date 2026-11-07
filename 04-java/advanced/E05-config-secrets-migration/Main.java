import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

public class Main {

    static class Source {
        final String value;
        final String origin;

        Source(String value, String origin) {
            this.value = value;
            this.origin = origin;
        }
    }

    // 优先级：环境变量 > 配置文件 > 默认值
    static Source resolve(String key, String def, Map<String, String> file, Map<String, String> env) {
        if (env.containsKey(key)) {
            return new Source(env.get(key), "env");
        }
        if (file.containsKey(key)) {
            return new Source(file.get(key), "file");
        }
        return new Source(def, "default");
    }

    static List<Integer> migrate(Set<Integer> applied, int[][] all, String[] sqls) {
        List<Integer> ran = new ArrayList<>();
        for (int i = 0; i < all.length; i++) {
            int version = all[i][0];
            if (applied.contains(version)) {
                continue;
            }
            applied.add(version);
            ran.add(version);
        }
        return ran;
    }

    public static void main(String[] args) {
        Map<String, String> file = new HashMap<>();
        file.put("db.name", "brm");
        file.put("port", "9090");
        Map<String, String> env = new HashMap<>();
        env.put("port", "8080");

        System.out.println("== 1. 配置优先级 ==");
        for (String k : new String[]{"port", "db.host", "db.name"}) {
            Source s = resolve(k, "localhost", file, env);
            System.out.println(k + " = " + s.value + "（来源：" + s.origin + "）");
        }

        System.out.println("== 2. 密钥缺失拒绝启动 ==");
        Map<String, String> secrets = new HashMap<>();
        if (!secrets.containsKey("DATABASE_URL")) {
            System.out.println("启动失败：DATABASE_URL 未设置（不要给密钥默认值）");
        }

        System.out.println("== 3. migration ==");
        int[][] all = {{1}, {2}, {3}};
        String[] sqls = {"create accounts", "add email", "add index"};
        Set<Integer> applied = new HashSet<>();
        for (int v : migrate(applied, all, sqls)) {
            System.out.println("应用 v" + v + " " + sqls[v - 1]);
        }
        System.out.println("当前版本 = " + applied.size());
        List<Integer> again = migrate(applied, all, sqls);
        System.out.println("再次执行：" + again.size() + " 个待应用，当前版本 = " + applied.size());

        System.out.println("== 4. 结论 ==");
        System.out.println("配置分层要能说出「这个值从哪来」，否则排查时全靠猜");
        System.out.println("migration 必须可重复执行：靠版本表跳过，不靠人记得跑过没有");
    }
}
