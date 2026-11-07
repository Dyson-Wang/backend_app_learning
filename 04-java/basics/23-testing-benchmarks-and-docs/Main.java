import java.util.List;

public class Main {
    // 手写断言：失败时带用例名与期望/实际，抛 AssertionError
    static void assertEquals(String caseName, int expected, int actual) {
        if (expected != actual) {
            throw new AssertionError("用例失败 [" + caseName + "]: 期望 " + expected + " 实际 " + actual);
        }
    }

    static void assertTrue(String caseName, boolean cond) {
        if (!cond) throw new AssertionError("用例失败 [" + caseName + "]: 期望 true");
    }

    /**
     * 计算整数的平方（演示 Javadoc 文档注释格式）。
     *
     * @param n 输入整数
     * @return n 的平方
     * @throws IllegalArgumentException 当 n 为负数时
     */
    static int square(int n) {
        if (n < 0) throw new IllegalArgumentException("n 不能为负");
        return n * n;
    }

    // 表驱动测试的数据：用例名、输入、期望输出
    record Case(String name, int input, int expected) {}

    public static void main(String[] args) {
        // 1) 表驱动测试：一组用例逐条断言
        List<Case> cases = List.of(
                new Case("零", 0, 0),
                new Case("一", 1, 1),
                new Case("五", 5, 25),
                new Case("九", 9, 81)
        );
        int passed = 0;
        for (var c : cases) {
            assertEquals(c.name(), c.expected(), square(c.input()));
            passed++;
        }
        System.out.println("表驱动测试通过: " + passed + "/" + cases.size());

        // 额外断言：覆盖正常与异常分支（分支覆盖的雏形）
        assertTrue("正数平方为正", square(3) > 0);
        boolean threw = false;
        try {
            square(-1);
        } catch (IllegalArgumentException e) {
            threw = true;
        }
        assertTrue("负数应抛异常", threw);
        System.out.println("额外断言(含异常路径)通过");

        // 2) 基准：System.nanoTime 测一段固定循环；耗时只作「差值」示意，不作结论
        long start = System.nanoTime();
        long acc = 0;
        for (int i = 0; i < 200_000; i++) acc += square(i % 100);
        long elapsed = System.nanoTime() - start;
        System.out.println("基准循环: 校验和=" + acc + " 耗时纳秒(仅示意,依赖机器)=" + elapsed);

        // 3) 文档注释：/** ... */ 由 javadoc 工具生成 HTML
        System.out.println("Javadoc: square 上有 /** */ 注释，可用 `javadoc Main.java` 生成文档");
        System.out.println("标准做法: 单测用 JUnit 5 + @Test，基准用 JMH(规避 JIT/GC 干扰)，覆盖率用 JaCoCo");
        System.out.println("全部断言通过，退出码 0");
    }
}
