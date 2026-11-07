import java.util.concurrent.CancellationException;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class Main {
    // 上下文传值：Java 无内置 context，显式把上下文对象沿异步链传递
    record Ctx(String traceId) {}

    public static void main(String[] args) throws Exception {
        ExecutorService pool = Executors.newFixedThreadPool(3);

        // 1) 异步任务模型：supplyAsync 把任务交给线程池
        CompletableFuture<Integer> async = CompletableFuture.supplyAsync(() -> 21 * 2, pool);
        System.out.println("异步任务结果=" + async.get());

        // 2) 超时：一个永不完成的 future，orTimeout 触发 TimeoutException（不用 sleep）
        CompletableFuture<String> never = new CompletableFuture<>();
        CompletableFuture<String> timed = never.orTimeout(40, TimeUnit.MILLISECONDS);
        try {
            timed.join();
        } catch (CompletionException e) {
            System.out.println("orTimeout 超时被捕获: " + e.getCause().getClass().getSimpleName());
        }

        // 3) 取消：手动 future 的 cancel 一定成功；取消后 join 抛 CancellationException
        CompletableFuture<String> cancellable = new CompletableFuture<>();
        boolean cancelled = cancellable.cancel(true);
        System.out.println("cancel(true) 返回=" + cancelled);
        try {
            cancellable.join();
        } catch (CancellationException e) {
            System.out.println("取消后 join -> CancellationException");
        }

        // 4) 取消/异常沿链传播：上游异常，下游 thenApply 也异常
        CompletableFuture<Integer> upstream = new CompletableFuture<>();
        CompletableFuture<String> downstream = upstream.thenApply(n -> "值=" + n);
        upstream.completeExceptionally(new IllegalStateException("上游失败"));
        try {
            downstream.join();
        } catch (CompletionException e) {
            System.out.println("异常沿链传播: " + e.getCause().getMessage());
        }

        // 5) 并发等待多个异步任务
        CompletableFuture<Integer> f1 = CompletableFuture.supplyAsync(() -> 10, pool);
        CompletableFuture<Integer> f2 = CompletableFuture.supplyAsync(() -> 20, pool);
        CompletableFuture<Integer> f3 = CompletableFuture.supplyAsync(() -> 30, pool);
        CompletableFuture.allOf(f1, f2, f3).join();
        System.out.println("allOf 并发等待后求和=" + (f1.join() + f2.join() + f3.join()));

        // 6) 取消后清理 + 上下文传值：whenComplete 无论成败都执行，Ctx 沿链传递
        Ctx ctx = new Ctx("trace-1");
        CompletableFuture<String> withCtx = CompletableFuture
                .supplyAsync(() -> "处理数据", pool)
                .thenApply(data -> data + " | ctx=" + ctx.traceId())
                .whenComplete((r, ex) -> System.out.println("whenComplete 清理(无论成败), 结果=" + r));
        System.out.println("上下文传值: " + withCtx.join());

        pool.shutdownNow();   // 收尾：终止线程池，保证进程 2 秒内退出
        System.out.println("executor 已 shutdownNow，所有异步任务结束");
    }
}
