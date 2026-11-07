public class Main {

    public static void main(String[] args) {
        final int clients = 100;
        final int downstreamCap = 100;

        System.out.println("== 1. 无退避：失败立即重试 3 次 ==");
        int noBackoff = clients * (1 + 3);
        System.out.println("下游收到 " + noBackoff + " QPS（容量 " + downstreamCap
                + "，超载 " + (noBackoff / downstreamCap) + " 倍）");

        System.out.println("== 2. 指数退避 + 全抖动 ==");
        // 3 次重试分散到 1s / 2s / 4s，每个窗口只承接 1/3 的客户端
        int peak = clients + clients / 3;
        System.out.println("下游峰值 ≈ " + peak + " QPS（接近容量，不再放大）");

        System.out.println("== 3. 再加熔断 ==");
        System.out.println("下游请求 = 0 QPS（open 状态本地快速失败，直接返回降级结果）");

        System.out.println("== 4. 结论 ==");
        System.out.println("重试必须有上限、退避和抖动；三个都没有，就是给故障加杠杆");
        System.out.println("熔断才是止损的关键：下游已经挂了，再发请求只会让它更晚恢复");
    }
}
