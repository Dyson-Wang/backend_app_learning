import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class Main {

    static void reply(com.sun.net.httpserver.HttpExchange ex, String body) throws IOException {
        byte[] b = body.getBytes("UTF-8");
        ex.sendResponseHeaders(200, b.length);
        try (OutputStream os = ex.getResponseBody()) {
            os.write(b);
        }
    }

    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        server.createContext("/slow", ex -> {
            try {
                Thread.sleep(120);          // 模拟在途慢请求
            } catch (InterruptedException ignored) {
                Thread.currentThread().interrupt();
            }
            reply(ex, "slow done");
        });
        server.createContext("/", ex -> reply(ex, "ok"));
        server.setExecutor(Executors.newFixedThreadPool(4));
        server.start();

        int port = server.getAddress().getPort();
        String base = "http://127.0.0.1:" + port;
        HttpClient client = HttpClient.newHttpClient();

        System.out.println("== 1. 发起一个 120ms 的在途请求 ==");
        CompletableFuture<HttpResponse<String>> future = client.sendAsync(
                HttpRequest.newBuilder(URI.create(base + "/slow")).build(),
                HttpResponse.BodyHandlers.ofString());
        Thread.sleep(30);

        System.out.println("== 2. 触发优雅关闭 stop(2)：最多等 2 秒 ==");
        long t0 = System.nanoTime();
        server.stop(2);                      // 停止 accept，等待在途 handler
        long cost = (System.nanoTime() - t0) / 1_000_000;

        System.out.println("== 3. 在途请求的结果 ==");
        System.out.println("body = " + future.get(3, TimeUnit.SECONDS).body()
                + "，关闭耗时 = " + cost + "ms");

        System.out.println("== 4. 关闭后新请求 ==");
        try (Socket s = new Socket()) {
            s.connect(new InetSocketAddress("127.0.0.1", port), 300);
            System.out.println("新请求 -> 连接成功（不应该发生）");
        } catch (IOException e) {
            System.out.println("新请求 -> 被拒绝 " + e.getClass().getSimpleName());
        }

        System.out.println("== 5. 结论 ==");
        System.out.println("stop(grace) 停止接收新连接，并给在途请求留 grace 秒");
        System.out.println("生产上 grace 要小于编排系统的 SIGKILL 宽限期，否则会被硬杀");
    }
}
