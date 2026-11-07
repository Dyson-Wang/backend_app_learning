import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class Main {

    static class Event {
        final String id;
        final int amount;

        Event(String id, int amount) {
            this.id = id;
            this.amount = amount;
        }
    }

    // seen 为 null 表示消费端没有幂等表
    static boolean handle(Set<String> seen, int[] balance, Event e) {
        if (seen != null) {
            if (seen.contains(e.id)) {
                return false;
            }
            seen.add(e.id);
        }
        balance[0] += e.amount;
        return true;
    }

    public static void main(String[] args) {
        List<Event> deliveries = List.of(
                new Event("e1", 10),
                new Event("e2", 10),
                new Event("e2", 10),   // ACK 丢了，Broker 重投
                new Event("e3", 10));

        System.out.println("== 1. 消费端没有幂等表 ==");
        int[] balance = {100};
        for (Event e : deliveries) {
            handle(null, balance, e);
        }
        System.out.println("投递 " + deliveries.size() + " 次，余额 100 -> "
                + balance[0] + "（多加了 10）");

        System.out.println("== 2. 消费端按消息 ID 去重 ==");
        balance = new int[]{100};
        Set<String> seen = new HashSet<>();
        int applied = 0;
        for (Event e : deliveries) {
            if (handle(seen, balance, e)) {
                applied++;
            }
        }
        System.out.println("投递 " + deliveries.size() + " 次，实际处理 " + applied
                + " 次，余额 100 -> " + balance[0]);

        System.out.println("== 3. 结论 ==");
        System.out.println("至少一次投递下，重复是常态；正确性必须由消费端幂等保证");
        System.out.println("去重表要落库并和业务写在同一事务里，否则还是会重复");
    }
}
