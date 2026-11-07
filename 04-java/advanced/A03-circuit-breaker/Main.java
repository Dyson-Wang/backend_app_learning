public class Main {
    enum State {
        CLOSED, OPEN, HALF_OPEN;

        String label() {
            return name().toLowerCase().replace('_', '-');
        }
    }

    static final class Breaker {
        private State state = State.CLOSED;
        private int failures = 0;
        private final int threshold;
        private final long cooldownMs;
        private long openedAt = 0L;
        private boolean probeInFlight = false;
        private long now = 0L;

        Breaker(int threshold, long cooldownMs) {
            this.threshold = threshold;
            this.cooldownMs = cooldownMs;
        }

        void advance(long ms) {
            now += ms;
        }

        synchronized String allow() {
            if (state == State.OPEN) {
                if (now - openedAt >= cooldownMs) {
                    state = State.HALF_OPEN;
                    probeInFlight = true;
                    return null;
                }
                return "circuit open: fast fail";
            }
            if (state == State.HALF_OPEN) {
                if (probeInFlight) {
                    return "circuit half-open: probe already in flight";
                }
                probeInFlight = true;
                return null;
            }
            return null;
        }

        synchronized void success() {
            failures = 0;
            probeInFlight = false;
            state = State.CLOSED;
        }

        synchronized void failure() {
            failures++;
            probeInFlight = false;
            if (state == State.HALF_OPEN || failures >= threshold) {
                state = State.OPEN;
                openedAt = now;
            }
        }

        synchronized State state() {
            return state;
        }
    }

    static String call(Breaker b, boolean downstreamOk) {
        String reject = b.allow();
        if (reject != null) {
            return "rejected: " + reject;
        }
        if (downstreamOk) {
            b.success();
            return "ok";
        }
        b.failure();
        return "downstream error";
    }

    public static void main(String[] args) {
        Breaker b = new Breaker(5, 200);

        System.out.println("== 1. 正常请求 ==");
        System.out.println("call -> " + call(b, true) + ", state=" + b.state().label());

        System.out.println("== 2. 连续 5 次失败触发熔断 ==");
        for (int i = 1; i <= 5; i++) {
            System.out.println("failure #" + i + " -> " + call(b, false) + ", state=" + b.state().label());
        }

        System.out.println("== 3. 熔断期间快速失败 ==");
        System.out.println("call -> " + call(b, true) + ", state=" + b.state().label());

        System.out.println("== 4. 冷却 200ms 后进入半开 ==");
        b.advance(200);
        System.out.println("allow -> " + (b.allow() == null ? "ok" : "rejected") + ", state=" + b.state().label());

        System.out.println("== 5. 探测成功恢复 closed ==");
        b.success();
        System.out.println("state=" + b.state().label() + ", 后续请求 -> " + call(b, true));

        System.out.println("== 6. 半开探测失败则重新 open ==");
        for (int i = 0; i < 5; i++) {
            b.failure();
        }
        b.advance(200);
        b.allow();
        b.failure();
        System.out.println("state=" + b.state().label() + "（探测失败后立即回到 open）");
    }
}
