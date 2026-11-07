import java.util.Optional;

public class Main {
    // 自定义「受检异常」：继承 Exception，调用方必须处理或在签名声明
    static class InsufficientFundsException extends Exception {
        InsufficientFundsException(String msg) { super(msg); }
    }

    // 自定义 AutoCloseable，演示 try-with-resources 的确定性清理
    static class FakeConn implements AutoCloseable {
        private final String name;
        FakeConn(String name) { this.name = name; }
        void query() { System.out.println("  执行查询 @" + name); }
        @Override public void close() {
            System.out.println("  资源已自动关闭 @" + name);
        }
    }

    // 重载：同名不同参数表，由编译期参数类型决定
    static int add(int a, int b) { return a + b; }
    static double add(double a, double b) { return a + b; }

    // 可变参数：本质就是数组
    static int sumAll(int... xs) {
        int s = 0;
        for (int x : xs) s += x;
        return s;
    }

    // 递归
    static long factorial(int n) {
        if (n <= 1) return 1;
        return n * factorial(n - 1);
    }

    // throws 把受检异常交给调用方处理
    static void withdraw(int balance, int amount) throws InsufficientFundsException {
        if (amount > balance) {
            throw new InsufficientFundsException("余额不足: 需要 " + amount + " 只有 " + balance);
        }
        System.out.println("  取款成功 " + amount);
    }

    public static void main(String[] args) {
        System.out.println("重载 int: " + add(2, 3));
        System.out.println("重载 double: " + add(2.5, 3.5));
        System.out.println("可变参数求和: " + sumAll(1, 2, 3, 4, 5));
        System.out.println("递归 5! = " + factorial(5));

        // 受检异常：必须 try/catch 或声明 throws
        try {
            withdraw(100, 30);
            withdraw(100, 500);
        } catch (InsufficientFundsException e) {
            System.out.println("捕获受检异常: " + e.getMessage());
        } finally {
            System.out.println("finally 一定执行（常用于兜底释放）");
        }

        // 非受检异常（RuntimeException）：编译器不强制处理，运行期才除零
        int zero = 0;
        try {
            int r = 10 / zero;
            System.out.println(r);
        } catch (ArithmeticException e) {
            System.out.println("捕获非受检异常: " + e.getClass().getSimpleName());
        }

        // try-with-resources：离开块时自动调用 close()，不依赖 finally
        try (var conn = new FakeConn("db-1")) {
            conn.query();
        }

        // Optional 表达「可能没有值」，避免返回 null
        var maybe = Optional.of("有值").filter(s -> s.length() > 5);
        System.out.println("Optional 过滤后是否为空=" + maybe.isEmpty());
    }
}
