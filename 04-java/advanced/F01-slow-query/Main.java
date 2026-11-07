public class Main {

    // 简化模型：所有请求同时到达，按到达顺序 FIFO 抢池子里的连接
    static int[] simulate(int[] durations, int pool, int waitTimeout) {
        int[] freeAt = new int[pool];
        int ok = 0;
        int timeout = 0;
        int total = 0;

        for (int d : durations) {
            int best = 0;
            for (int i = 1; i < pool; i++) {
                if (freeAt[i] < freeAt[best]) {
                    best = i;
                }
            }
            int wait = freeAt[best];
            if (wait > waitTimeout) {
                timeout++;
                continue;
            }
            freeAt[best] = wait + d;
            if (freeAt[best] > total) {
                total = freeAt[best];
            }
            ok++;
        }
        return new int[]{ok, timeout, total};
    }

    public static void main(String[] args) {
        int[] normal = new int[10];
        for (int i = 0; i < normal.length; i++) {
            normal[i] = 30;
        }

        System.out.println("== 1. 正常：全部 30ms ==");
        int[] r = simulate(normal, 2, 100);
        System.out.println("成功 = " + r[0] + "，超时 = " + r[1] + "，总耗时 = " + r[2] + "ms");

        System.out.println("== 2. 混入 1 个 300ms 慢查询（等待上限 100ms） ==");
        int[] slow = normal.clone();
        slow[0] = 300;
        r = simulate(slow, 2, 100);
        System.out.println("成功 = " + r[0] + "，超时 = " + r[1] + "，总耗时 = " + r[2] + "ms");

        System.out.println("== 3. 结论 ==");
        System.out.println("池子只有 2 个连接，1 个慢查询就占掉一半；等待超时把雪崩挡在入口");
        System.out.println("慢查询必须配 statement_timeout + 索引，光加连接数只是把问题推后");
    }
}
