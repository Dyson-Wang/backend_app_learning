import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicLong;

public class Main {
    private static final AtomicLong NEXT_ID = new AtomicLong(1);
    private static final Map<Long, String> USERS = new ConcurrentHashMap<>();

    public static void main(String[] args) throws Exception {
        int port = Integer.parseInt(System.getenv().getOrDefault("PORT", "8083"));
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);

        server.createContext("/healthz", ex -> send(ex, 200, "{\"status\":\"ok\"}"));
        server.createContext("/readyz", ex -> send(ex, 200, "{\"status\":\"ready\"}"));
        server.createContext("/users", Main::handleUsers);
        server.createContext("/users/", Main::handleUserById);

        // 不用默认的单线程 executor，生产会换成带界队列的 ThreadPoolExecutor 或虚拟线程
        server.setExecutor(Executors.newFixedThreadPool(8));

        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("SIGTERM received, draining");
            server.stop(2); // 最多等 2 秒，超时强退
            System.out.println("stopped");
        }));

        server.start();
        System.out.println("listening on http://localhost:" + port);
    }

    private static void handleUsers(HttpExchange ex) throws IOException {
        switch (ex.getRequestMethod()) {
            case "GET" -> send(ex, 200, listJson());
            case "POST" -> {
                String raw = new String(ex.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
                String name = extractName(raw);
                if (name == null || name.isBlank()) {
                    send(ex, 400, "{\"error\":{\"code\":\"VALIDATION\",\"message\":\"name is required\"}}");
                    return;
                }
                long id = NEXT_ID.getAndIncrement();
                USERS.put(id, name);
                send(ex, 201, "{\"data\":{\"id\":" + id + ",\"name\":\"" + escape(name) + "\"}}");
            }
            default -> {
                ex.getResponseHeaders().set("Allow", "GET, POST");
                send(ex, 405, "{\"error\":{\"code\":\"METHOD_NOT_ALLOWED\",\"message\":\"use GET or POST\"}}");
            }
        }
    }

    private static void handleUserById(HttpExchange ex) throws IOException {
        if (!"GET".equals(ex.getRequestMethod())) {
            send(ex, 405, "{\"error\":{\"code\":\"METHOD_NOT_ALLOWED\",\"message\":\"use GET\"}}");
            return;
        }
        String raw = ex.getRequestURI().getPath().substring("/users/".length());
        try {
            long id = Long.parseLong(raw);
            String name = USERS.get(id);
            if (name == null) {
                send(ex, 404, "{\"error\":{\"code\":\"NOT_FOUND\",\"message\":\"user not found\"}}");
            } else {
                send(ex, 200, "{\"data\":{\"id\":" + id + ",\"name\":\"" + escape(name) + "\"}}");
            }
        } catch (NumberFormatException e) {
            send(ex, 400, "{\"error\":{\"code\":\"BAD_ID\",\"message\":\"id must be an integer\"}}");
        }
    }

    private static String listJson() {
        StringBuilder sb = new StringBuilder("{\"data\":[");
        boolean first = true;
        for (var entry : USERS.entrySet()) {
            if (!first) sb.append(',');
            first = false;
            sb.append("{\"id\":").append(entry.getKey())
              .append(",\"name\":\"").append(escape(entry.getValue())).append("\"}");
        }
        return sb.append("]}").toString();
    }

    // 教学用简化解析：生产请用 Jackson / Gson，别手写 JSON
    private static String extractName(String json) {
        int key = json.indexOf("\"name\"");
        if (key < 0) return null;
        int colon = json.indexOf(':', key);
        if (colon < 0) return null;
        int first = json.indexOf('"', colon);
        if (first < 0) return null;
        int last = json.indexOf('"', first + 1);
        if (last < 0) return null;
        return json.substring(first + 1, last);
    }

    private static String escape(String s) {
        return s.replace("\\", "\\\\").replace("\"", "\\\"");
    }

    private static void send(HttpExchange ex, int status, String body) throws IOException {
        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);
        ex.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        String requestId = ex.getRequestHeaders().getFirst("X-Request-ID");
        ex.getResponseHeaders().set("X-Request-ID", requestId == null ? "-" : requestId);
        ex.sendResponseHeaders(status, bytes.length);
        try (var out = ex.getResponseBody()) {
            out.write(bytes);
        }
        System.out.printf("%s %s %d request_id=%s%n",
                ex.getRequestMethod(), ex.getRequestURI().getPath(), status,
                requestId == null ? "-" : requestId);
    }
}
