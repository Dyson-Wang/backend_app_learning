public class Main {
    static final class DownstreamException extends Exception {
        final boolean retryable;

        DownstreamException(String message, boolean retryable) {
            super(message);
            this.retryable = retryable;
        }
    }

    interface Attempt {
        String run(int attempt) throws DownstreamException;
    }

    static String doWithRetry(int maxAttempts, long baseMs, Attempt fn) throws InterruptedException {
        DownstreamException last = null;
        for (int attempt = 1; attempt <= maxAttempts; attempt++) {
            try {
                String value = fn.run(attempt);
                return "ok=true attempt=" + attempt + " value=" + value;
            } catch (DownstreamException e) {
                last = e;
                if (!e.retryable) {
                    return "ok=false attempt=" + attempt + " error=" + e.getMessage() + " permanent=true";
                }
                if (attempt == maxAttempts) {
                    break;
                }
                long delay = baseMs * (1L << (attempt - 1));
                System.out.println("  attempt " + attempt + " 失败（" + e.getMessage() + "），" + delay + "ms 后重试");
                Thread.sleep(delay);
            }
        }
        return "ok=false attempt=" + maxAttempts + " error=" + (last == null ? "" : last.getMessage()) + " permanent=false";
    }

    public static void main(String[] args) throws Exception {
        final int[] calls = {0};

        System.out.println("== 1. 可重试错误：指数退避后成功 ==");
        System.out.println("结果: " + doWithRetry(4, 10, attempt -> {
            calls[0]++;
            if (calls[0] < 4) {
                throw new DownstreamException("503 upstream", true);
            }
            return "第 " + calls[0] + " 次成功";
        }));

        System.out.println("== 2. 永久错误：立即停止 ==");
        System.out.println("结果: " + doWithRetry(4, 10, attempt -> {
            throw new DownstreamException("400 bad request", false);
        }));

        System.out.println("== 3. 一直失败：耗尽 3 次 ==");
        System.out.println("结果: " + doWithRetry(3, 5, attempt -> {
            throw new DownstreamException("timeout", true);
        }));

        System.out.println("结论：只重试可重试错误，退避必须指数增长；生产还要加抖动，避免重试风暴");
    }
}
