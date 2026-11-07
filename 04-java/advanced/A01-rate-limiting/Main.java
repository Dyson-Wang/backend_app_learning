public class Main {
    static final class TokenBucket {
        private final double rate;
        private final double burst;
        private double tokens;
        private long lastMs = 0L;
        private long nowMs = 0L;

        TokenBucket(double rate, double burst) {
            this.rate = rate;
            this.burst = burst;
            this.tokens = burst;
        }

        void advance(long ms) {
            nowMs += ms;
        }

        boolean allow() {
            double elapsedSec = (nowMs - lastMs) / 1000.0;
            tokens = Math.min(burst, tokens + elapsedSec * rate);
            lastMs = nowMs;
            if (tokens >= 1.0) {
                tokens -= 1.0;
                return true;
            }
            return false;
        }

        double tokens() {
            return tokens;
        }
    }

    static void hit(TokenBucket bucket, String label) {
        boolean ok = bucket.allow();
        System.out.printf("  %s -> allow=%s, tokens=%.2f%n", label, ok, bucket.tokens());
    }

    public static void main(String[] args) {
        TokenBucket bucket = new TokenBucket(2, 5);

        System.out.println("== t=0ms：burst=5，连续 6 次 ==");
        for (int i = 1; i <= 6; i++) {
            hit(bucket, "#" + i);
        }

        System.out.println("== t=500ms：补 1 个令牌 ==");
        bucket.advance(500);
        hit(bucket, "#7");
        hit(bucket, "#8");

        System.out.println("== t=1500ms：补 2 个令牌 ==");
        bucket.advance(1000);
        hit(bucket, "#9");
        hit(bucket, "#10");
        hit(bucket, "#11");

        System.out.println("结论：令牌桶允许突发（burst），但长期速率被 rate 卡死");
    }
}
