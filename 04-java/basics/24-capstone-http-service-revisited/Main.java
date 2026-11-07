import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicLong;

public class Main {
    // 枚举：HTTP 状态与错误码的映射
    enum Code {
        OK(200), CREATED(201), BAD_REQUEST(400), NOT_FOUND(404), INTERNAL(500);
        private final int http;
        Code(int http) { this.http = http; }
        int http() { return http; }
    }

    // 泛型：统一响应包装（泛型 record + 泛型静态工厂方法）
    record ApiResult<T>(boolean ok, T data, String code, String message) {
        static <R> ApiResult<R> ok(R data) { return new ApiResult<>(true, data, "OK", null); }
        static <R> ApiResult<R> fail(String code, String message) {
            return new ApiResult<>(false, null, code, message);
        }
    }

    record User(long id, String name) {}

    private static final AtomicLong NEXT_ID = new AtomicLong(1);
    private static final Map<Long, User> STORE = new ConcurrentHashMap<>();
    private static final AtomicLong REQ_SEQ = new AtomicLong(0);   // 单调序号做日志
    // 时间特性：用「固定时刻」格式化，避免把当前绝对时间当结论打印
    private static final String LOG_TIME =
            DateTimeFormatter.ofPattern("HH:mm:ss").format(LocalTime.of(9, 0, 0));

    public static void main(String[] args) throws Exception {
        // 随机端口：bind 到端口 0，由系统分配，不打印具体值以保证输出可复现
        HttpServer server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        int port = server.getAddress().getPort();

        // 闭包：下面的 lambda 捕获了 effectively final 的局部变量 service
        final String service = "brm-users";
        server.createContext("/health", ex -> dispatch(ex, service));
        server.createContext("/users", ex -> dispatch(ex, service));
        server.createContext("/users/", ex -> dispatch(ex, service));
        server.createContext("/", ex -> dispatch(ex, service));   // 兜底路由

        ExecutorService pool = Executors.newFixedThreadPool(4);
        server.setExecutor(pool);
        server.start();
        System.out.println("已在随机端口启动(端口由系统分配, 不打印具体值)");

        // 自请求：HttpClient 带超时，任何一步都不会卡住
        HttpClient client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(2))
                .build();
        System.out.println("GET /health -> " + call(client, port, "GET", "/health", null));
        System.out.println("POST /users(alice) -> "
                + call(client, port, "POST", "/users", "{\"name\":\"alice\"}"));
        System.out.println("POST /users(空名,校验失败) -> "
                + call(client, port, "POST", "/users", "{\"name\":\"\"}"));
        System.out.println("GET /users -> " + call(client, port, "GET", "/users", null));
        System.out.println("GET /users/999(不存在) -> "
                + call(client, port, "GET", "/users/999", null));
        System.out.println("内存用户数=" + STORE.size());

        // 3 秒内自行关闭并退出
        server.stop(0);
        pool.shutdownNow();
        System.out.println("已 stop(0) 并关闭 executor");
        System.exit(0);   // 兜底：确保任何残留非守护线程都不拖住 JVM
    }

    // 路由分发：按真实请求路径 + 方法分派；错误处理统一成错误码
    private static void dispatch(HttpExchange ex, String service) throws IOException {
        String path = ex.getRequestURI().getPath();
        String method = ex.getRequestMethod();
        ApiResult<?> result;
        try {
            if (path.equals("/health") && "GET".equals(method)) {
                result = ApiResult.ok(Map.of("status", "up"));
            } else if (path.equals("/users") && "GET".equals(method)) {
                result = ApiResult.ok(STORE.values().toString());
            } else if (path.equals("/users") && "POST".equals(method)) {
                result = createUser(ex);
            } else if (path.startsWith("/users/") && "GET".equals(method)) {
                result = findUser(ex);
            } else {
                result = ApiResult.fail("NOT_FOUND", "没有这个路由");
            }
        } catch (Exception e) {
            result = ApiResult.fail("INTERNAL", e.getMessage());
        }
        send(ex, service, result);
    }

    private static ApiResult<?> createUser(HttpExchange ex) throws IOException {
        String raw = new String(ex.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
        String name = strField(raw, "name");
        // 校验 + 错误码
        if (name == null || name.isBlank()) {
            return ApiResult.fail("VALIDATION", "name 不能为空");
        }
        long id = NEXT_ID.getAndIncrement();
        User u = new User(id, name);
        STORE.put(id, u);
        return new ApiResult<>(true, u, "CREATED", null);
    }

    private static ApiResult<?> findUser(HttpExchange ex) {
        String tail = ex.getRequestURI().getPath().substring("/users/".length());
        try {
            long id = Long.parseLong(tail);
            User u = STORE.get(id);
            if (u == null) return ApiResult.fail("NOT_FOUND", "用户不存在");
            return ApiResult.ok(u);
        } catch (NumberFormatException e) {
            return ApiResult.fail("BAD_ID", "id 必须是整数");
        }
    }

    private static void send(HttpExchange ex, String service, ApiResult<?> result) throws IOException {
        int status = switch (result.code()) {
            case "OK" -> Code.OK.http();
            case "CREATED" -> Code.CREATED.http();
            case "VALIDATION", "BAD_ID" -> Code.BAD_REQUEST.http();
            case "NOT_FOUND" -> Code.NOT_FOUND.http();
            default -> Code.INTERNAL.http();
        };
        byte[] bytes = toJson(result).getBytes(StandardCharsets.UTF_8);
        // 日志：先记录再写响应体，保证输出顺序稳定；时间用固定值格式化
        System.out.println("  [" + LOG_TIME + " #" + REQ_SEQ.incrementAndGet() + "] "
                + service + " " + ex.getRequestMethod() + " " + ex.getRequestURI().getPath()
                + " -> " + status);
        ex.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        ex.sendResponseHeaders(status, bytes.length);
        try (var out = ex.getResponseBody()) {
            out.write(bytes);
        }
    }

    // 教学用手写 JSON 序列化（泛型 ApiResult 统一走这里）
    private static String toJson(ApiResult<?> r) {
        StringBuilder sb = new StringBuilder("{\"ok\":").append(r.ok());
        if (r.ok()) {
            sb.append(",\"data\":").append(jsonValue(r.data()));
        } else {
            sb.append(",\"error\":{\"code\":\"").append(r.code())
              .append("\",\"message\":\"").append(escape(r.message())).append("\"}");
        }
        return sb.append("}").toString();
    }

    private static String jsonValue(Object v) {
        if (v instanceof User u) {   // instanceof 模式匹配做解构
            return "{\"id\":" + u.id() + ",\"name\":\"" + escape(u.name()) + "\"}";
        }
        return "\"" + escape(String.valueOf(v)) + "\"";
    }

    // 教学用极简解析：取字符串字段
    private static String strField(String json, String key) {
        int k = json.indexOf("\"" + key + "\"");
        if (k < 0) return null;
        int colon = json.indexOf(':', k);
        if (colon < 0) return null;
        int first = json.indexOf('"', colon);
        if (first < 0) return null;
        int last = json.indexOf('"', first + 1);
        if (last < 0) return null;
        return json.substring(first + 1, last);
    }

    private static String escape(String s) {
        return s == null ? "" : s.replace("\\", "\\\\").replace("\"", "\\\"");
    }

    // 自请求辅助：返回状态码
    private static int call(HttpClient client, int port, String method, String path, String body)
            throws Exception {
        HttpRequest.Builder b = HttpRequest.newBuilder()
                .uri(URI.create("http://127.0.0.1:" + port + path))
                .timeout(Duration.ofSeconds(2));
        if (body == null) {
            b = "GET".equals(method) ? b.GET()
                    : b.method(method, HttpRequest.BodyPublishers.noBody());
        } else {
            b = b.method(method, HttpRequest.BodyPublishers.ofString(body));
        }
        HttpResponse<String> resp = client.send(b.build(), HttpResponse.BodyHandlers.ofString());
        return resp.statusCode();
    }
}
