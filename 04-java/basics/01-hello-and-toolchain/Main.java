import java.util.List;

public class Main {
    public static void main(String[] args) {
        System.out.println("hello, backend roadmap");
        System.out.println("java:   " + System.getProperty("java.version"));
        System.out.println("vendor: " + System.getProperty("java.vendor"));
        System.out.println("os:     " + System.getProperty("os.name") + " " + System.getProperty("os.arch"));
        System.out.println("cpus:   " + Runtime.getRuntime().availableProcessors());

        // var 是局部变量类型推断（Java 10+），不是动态类型，编译期就确定
        var name = "dyson";
        final int maxRetries = 3;
        var nums = List.of(1, 2, 3, 4);
        var sum = 0;
        for (var n : nums) {
            sum += n;
        }
        System.out.printf("name=%s maxRetries=%d sum=%d%n", name, maxRetries, sum);

        // 受检异常：编译器强制你处理或声明，这是 Java 独有的错误处理方式
        try {
            Thread.sleep(1);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        // String 不可变；频繁拼接用 StringBuilder
        var builder = new StringBuilder("a");
        for (int i = 0; i < 3; i++) {
            builder.append(i);
        }
        System.out.println("builder=" + builder);
    }
}
