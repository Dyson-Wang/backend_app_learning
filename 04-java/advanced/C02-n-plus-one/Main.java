import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Main {

    static class Order {
        final int id;
        final int accountId;
        final int amount;

        Order(int id, int accountId, int amount) {
            this.id = id;
            this.accountId = accountId;
            this.amount = amount;
        }
    }

    static class Account {
        final int id;
        final String name;

        Account(int id, String name) {
            this.id = id;
            this.name = name;
        }
    }

    static class FakeDB {
        final List<Order> orders = new ArrayList<>();
        final Map<Integer, Account> accounts = new HashMap<>();
        int queries = 0;

        FakeDB() {
            for (int i = 1; i <= 10; i++) {
                orders.add(new Order(i, i % 3 + 1, 100 + i));
            }
            accounts.put(1, new Account(1, "alice"));
            accounts.put(2, new Account(2, "bob"));
            accounts.put(3, new Account(3, "carol"));
        }

        List<Order> listOrders() {
            queries++;
            return orders;
        }

        Account getAccount(int id) {
            queries++;
            return accounts.get(id);
        }

        Map<Integer, Account> getAccounts(List<Integer> ids) {
            queries++;
            Map<Integer, Account> m = new HashMap<>();
            for (int id : ids) {
                m.put(id, accounts.get(id));
            }
            return m;
        }
    }

    static List<String> nPlusOne(FakeDB db) {
        List<String> out = new ArrayList<>();
        for (Order o : db.listOrders()) {                 // 1 次
            out.add(o.id + ":" + db.getAccount(o.accountId).name);   // N 次
        }
        return out;
    }

    static List<String> batched(FakeDB db) {
        List<Order> orders = db.listOrders();             // 1 次
        List<Integer> ids = new ArrayList<>(new java.util.TreeSet<>(
                orders.stream().map(o -> o.accountId).toList()));
        Map<Integer, Account> m = db.getAccounts(ids);    // 1 次
        List<String> out = new ArrayList<>();
        for (Order o : orders) {
            out.add(o.id + ":" + m.get(o.accountId).name);
        }
        return out;
    }

    public static void main(String[] args) {
        FakeDB db1 = new FakeDB();
        List<String> a = nPlusOne(db1);
        System.out.println("N+1：" + a.size() + " 条订单 -> 查询 " + db1.queries + " 次");

        FakeDB db2 = new FakeDB();
        List<String> b = batched(db2);
        System.out.println("批量：" + b.size() + " 条订单 -> 查询 " + db2.queries + " 次");

        System.out.println("结果一致：" + a.equals(b));
        System.out.println("结论：N+1 的代价是 1+N 次往返，批量把 N 压成 1");
    }
}
