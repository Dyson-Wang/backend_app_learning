public class Main {

    static class Event {
        final String id;
        final String kind;

        Event(String id, String kind) {
            this.id = id;
            this.kind = kind;
        }
    }

    static String[] apply(String state, Event e) {
        if (e.kind.equals("create") && state.isEmpty()) {
            return new String[]{"CREATED", ""};
        }
        if (e.kind.equals("pay") && state.equals("CREATED")) {
            return new String[]{"PAID", ""};
        }
        if (e.kind.equals("ship") && state.equals("PAID")) {
            return new String[]{"SHIPPED", ""};
        }
        String shown = state.isEmpty() ? "INIT" : state;
        return new String[]{state, e.kind + " 在 " + shown + " 状态下不合法"};
    }

    static String[] consume(Event[] events) {
        String state = "";
        for (Event e : events) {
            String[] r = apply(state, e);
            state = r[0];
            if (!r[1].isEmpty()) {
                return new String[]{state, r[1]};
            }
        }
        return new String[]{state, ""};
    }

    static int partition(String key, int n) {
        int h = 0;
        for (char ch : key.toCharArray()) {
            h = (h * 31 + ch) % n;
        }
        return h;
    }

    public static void main(String[] args) {
        Event create = new Event("c1", "create");
        Event pay = new Event("c2", "pay");
        Event ship = new Event("c3", "ship");

        System.out.println("== 1. 同一 key 被并发消费（顺序打乱） ==");
        String[] r = consume(new Event[]{ship, create, pay});
        System.out.println("最终状态 = " + (r[0].isEmpty() ? "INIT" : r[0]) + "，错误 = " + r[1]);

        System.out.println("== 2. 按 key 路由到同一分区、单消费者串行处理 ==");
        r = consume(new Event[]{create, pay, ship});
        System.out.println("最终状态 = " + r[0] + "，错误 = " + (r[1].isEmpty() ? "无" : r[1]));

        System.out.println("== 3. 分区函数：同一 order_id 必然落到同一分区 ==");
        for (String k : new String[]{"order-42", "order-43", "order-42"}) {
            System.out.println(k + " -> 分区 " + partition(k, 4));
        }

        System.out.println("== 4. 结论 ==");
        System.out.println("顺序性不是全局的，是「同一分区内」的；跨 key 的顺序本来就没有意义");
        System.out.println("要顺序就牺牲并发：单分区单消费者，或者消费端按 key 做串行队列");
    }
}
