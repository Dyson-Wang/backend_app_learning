public class Main {

    public static void main(String[] args) {
        final int backlog = 1000;
        final int perSec = 10;      // 单消费者每秒 10 条
        final int producers = 20;   // 生产端每秒 20 条

        System.out.println("== 1. 单消费者，逐条处理 ==");
        System.out.println(backlog + " 条 / " + perSec + " 条每秒 = " + (backlog / perSec) + " 秒");

        System.out.println("== 2. 扩到 3 个消费者 ==");
        int three = (backlog + perSec * 3 - 1) / (perSec * 3);
        System.out.println(backlog + " 条 / " + (perSec * 3) + " 条每秒 = " + three + " 秒");

        System.out.println("== 3. 批量 10 条一次（单消费者） ==");
        int batch = (backlog + perSec * 10 - 1) / (perSec * 10);
        System.out.println(backlog + " 条 / " + (perSec * 10) + " 条每秒 = " + batch + " 秒");

        System.out.println("== 4. 生产仍在继续（每秒 20 条） ==");
        int net = perSec * 10 - producers;
        int clear = (backlog + net - 1) / net;
        System.out.println("净消化 " + net + " 条每秒，清空 " + backlog + " 条需要 " + clear + " 秒");

        System.out.println("== 5. 结论 ==");
        System.out.println("积压先算「净消化速率」，不是「消费速率」——生产端还在往里灌");
        System.out.println("扩消费者最快但受分区数限制；批量提升单机吞吐，代价是延迟和失败重放");
    }
}
