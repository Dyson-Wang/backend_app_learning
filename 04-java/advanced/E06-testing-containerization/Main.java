public class Main {

    static int failures = 0;

    static void check(String name, boolean ok, String detail) {
        if (ok) {
            System.out.println("PASS " + name);
            return;
        }
        failures++;
        System.out.println("FAIL " + name + ": " + detail);
    }

    static int add(int a, int b) {
        return a + b;
    }

    static int divide(int a, int b) {
        if (b == 0) {
            throw new ArithmeticException("division by zero");
        }
        return a / b;
    }

    static int parseAge(String s) {
        if (s.isEmpty()) {
            throw new IllegalArgumentException("empty");
        }
        return 18;                 // 故意不校验非数字，用来触发失败用例
    }

    public static void main(String[] args) {
        System.out.println("== 1. 单元测试 ==");
        check("add(1,2)=3", add(1, 2) == 3, "期望 3");

        boolean threw = false;
        try {
            divide(1, 0);
        } catch (ArithmeticException e) {
            threw = true;
        }
        check("divide(1,0) 抛错", threw, "期望抛错");

        Throwable ageErr = null;
        try {
            parseAge("abc");
        } catch (RuntimeException e) {
            ageErr = e;
        }
        check("parseAge(\"abc\") 抛错", ageErr != null, "期望抛错，实际正常返回 18");

        System.out.println("== 2. 汇总 ==");
        System.out.println("3 个用例：" + (3 - failures) + " 通过，" + failures + " 失败");

        System.out.println("== 3. 结论 ==");
        System.out.println("最后一条是故意失败的用例，用来看失败输出长什么样");
        System.out.println("单测不碰网络/DB；集成测试才起真实依赖，两者分开跑");
        System.out.println("容器化：多阶段构建、非 root、HEALTHCHECK 打 /healthz");
    }
}
