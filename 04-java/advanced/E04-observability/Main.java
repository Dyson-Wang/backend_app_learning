import java.util.concurrent.atomic.AtomicLong;

public class Main {
    static final class Metrics {
        private long requests = 0;
        private long errors = 0;
        private long totalMs = 0;
        private long maxMs = 0;

        synchronized void observe(int status, long ms) {
            requests++;
            if (status >= 500) {
                errors++;
            }
            totalMs += ms;
            if (ms > maxMs) {
                maxMs = ms;
            }
        }

        synchronized long[] snapshot() {
            return new long[] {requests, errors, totalMs, maxMs};
        }
    }

    private static final AtomicLong TRACE_SEQ = new AtomicLong(0);

    static String json(String level, String traceId, String method, String path,
                       int status, long ms, String msg) {
        StringBuilder sb = new StringBuilder();
        sb.append("{\"level\":\"").append(level).append("\"")
          .append(",\"trace_id\":\"").append(traceId).append("\"")
          .append(",\"method\":\"").append(method).append("\"")
          .append(",\"path\":\"").append(path).append("\"")
          .append(",\"status\":").append(status)
          .append(",\"latency_ms\":").append(ms);
        if (msg != null) {
            sb.append(",\"msg\":\"").append(msg).append("\"");
        }
        return sb.append("}").toString();
    }

    static void handle(Metrics metrics, String method, String path, int status, long ms) {
        String traceId = String.format("trace-%04d", TRACE_SEQ.incrementAndGet());
        metrics.observe(status, ms);
        String level = status >= 500 ? "error" : "info";
        String msg = status >= 500 ? "upstream failed" : null;
        System.out.println(json(level, traceId, method, path, status, ms, msg));
    }

    public static void main(String[] args) {
        Metrics metrics = new Metrics();
        Object[][] reqs = {
            {"GET", "/healthz", 200, 1L},
            {"GET", "/users", 200, 12L},
            {"POST", "/users", 201, 18L},
            {"GET", "/users/999", 404, 3L},
            {"POST", "/users", 400, 2L},
            {"GET", "/orders", 500, 240L},
            {"GET", "/users", 200, 9L},
            {"GET", "/users/1", 200, 7L},
            {"GET", "/orders", 500, 260L},
            {"DELETE", "/users/1", 204, 5L},
        };

        for (Object[] r : reqs) {
            handle(metrics, (String) r[0], (String) r[1], (Integer) r[2], (Long) r[3]);
        }

        long[] snap = metrics.snapshot();
        long requests = snap[0];
        long errors = snap[1];
        double avg = requests == 0 ? 0.0 : (double) snap[2] / requests;

        System.out.println("---- RED 指标 ----");
        System.out.printf("requests=%d errors=%d error_rate=%.2f%% avg_ms=%.1f max_ms=%d%n",
                requests, errors, (double) errors / requests * 100, avg, snap[3]);
        System.out.println("排查顺序：先看 error_rate 与 max_ms，再用 trace_id 串起同一次请求的全部日志");
    }
}
