public class Main {

    static class Msg {
        final String id;
        final int failures;
        final boolean permanent;

        Msg(String id, int failures, boolean permanent) {
            this.id = id;
            this.failures = failures;
            this.permanent = permanent;
        }
    }

    static class Result {
        final boolean success;
        final int attempts;

        Result(boolean success, int attempts) {
            this.success = success;
            this.attempts = attempts;
        }
    }

    static Result deliver(Msg m, int maxAttempts) {
        for (int attempt = 1; attempt <= maxAttempts; attempt++) {
            if (m.permanent) {
                System.out.println(m.id + " -> 永久错误，直接进死信（不重试）");
                return new Result(false, attempt);
            }
            if (attempt > m.failures) {
                System.out.println(m.id + " -> 成功（第 " + attempt + " 次）");
                return new Result(true, attempt);
            }
            if (attempt < maxAttempts) {
                int backoff = 10 << (attempt - 1);
                System.out.println(m.id + " -> 可重试错误，退避 " + backoff + "ms 后重试");
                continue;
            }
            System.out.println(m.id + " -> 重试耗尽，进死信");
        }
        return new Result(false, maxAttempts);
    }

    public static void main(String[] args) {
        Msg[] msgs = {
                new Msg("m1", 0, false),
                new Msg("m2", 2, false),
                new Msg("m3", 0, true),
                new Msg("m4", 99, false),
                new Msg("m5", 0, false),
        };

        System.out.println("== 1. 逐条投递（最多 3 次） ==");
        int ok = 0;
        int dead = 0;
        int total = 0;
        for (Msg m : msgs) {
            Result r = deliver(m, 3);
            total += r.attempts;
            if (r.success) {
                ok++;
            } else {
                dead++;
            }
        }

        System.out.println("== 2. 统计 ==");
        System.out.println("成功 = " + ok + "，死信 = " + dead + "，总投递 = " + total);

        System.out.println("== 3. 结论 ==");
        System.out.println("永久错误不要重试，重试只会拖慢整个队列");
        System.out.println("死信队列必须有人看：它进 DLQ 说明要么代码有 bug，要么数据有问题");
    }
}
