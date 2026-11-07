public class Main {

    static class Deadline {
        private final long atNanos;

        Deadline(long budgetMs) {
            this.atNanos = System.nanoTime() + budgetMs * 1_000_000L;
        }

        long remainingMs() {
            return (atNanos - System.nanoTime()) / 1_000_000L;
        }
    }

    // deadline 为 null 表示中间层没传，下游就完全不看上游预算
    static String downstream(Deadline d, long workMs) throws InterruptedException {
        if (d == null) {
            Thread.sleep(workMs);
            return "ok";
        }
        long left = d.remainingMs();
        if (left < workMs) {
            if (left > 0) {
                Thread.sleep(left);      // 有预算意识：跑到截止就放弃
            }
            return "deadline exceeded";
        }
        Thread.sleep(workMs);
        return "ok";
    }

    public static void main(String[] args) throws Exception {
        System.out.println("== 1. 上游预算 50ms，下游要 200ms ==");
        long t0 = System.nanoTime();
        String r1 = downstream(new Deadline(50), 200);
        long c1 = (System.nanoTime() - t0) / 1_000_000L;
        System.out.println("res = " + r1 + "，耗时 = " + c1 + "ms（远小于 200ms）");

        System.out.println("== 2. 中间层丢掉 deadline ==");
        long t1 = System.nanoTime();
        String r2 = downstream(null, 200);
        long c2 = (System.nanoTime() - t1) / 1_000_000L;
        System.out.println("res = " + r2 + "，耗时 = " + c2 + "ms（上游再急也没用）");

        System.out.println("== 3. 正确写法：orTimeout 把预算交给 CompletableFuture ==");
        long t2 = System.nanoTime();
        try {
            java.util.concurrent.CompletableFuture
                    .supplyAsync(() -> {
                        try {
                            Thread.sleep(200);
                        } catch (InterruptedException ignored) {
                            Thread.currentThread().interrupt();
                        }
                        return "ok";
                    })
                    .orTimeout(50, java.util.concurrent.TimeUnit.MILLISECONDS)
                    .join();
        } catch (java.util.concurrent.CompletionException e) {
            System.out.println(e.getCause().getClass().getSimpleName()
                    + "，耗时 = " + ((System.nanoTime() - t2) / 1_000_000L) + "ms");
        }

        System.out.println("== 4. 结论 ==");
        System.out.println("deadline 必须一路显式传；任何一层忘记传，取消就断在这里");
        System.out.println("跨服务用 X-Request-Deadline 头传绝对时间戳，避免时钟漂移累积");
    }
}
