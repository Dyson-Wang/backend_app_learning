import java.util.List;

public class Main {
    // 成员变量（字段）才有零值；局部变量没有默认值，必须显式初始化
    static class Defaults {
        int i;
        long l;
        double d;
        boolean flag;
        char c;
        String s;
        int[] arr;
    }

    public static void main(String[] args) {
        // 8 种基本类型 + 引用类型
        byte b = 12;
        short sh = 300;
        int i = 100_000;              // 下划线只影响可读性，编译期忽略
        long l = 9_000_000_000L;      // long 字面量要加 L
        float f = 3.14f;              // float 字面量要加 f
        double d = 2.71828;
        char c = '中';
        boolean ok = true;
        String name = "dyson";        // 引用类型，零值是 null
        System.out.println("基本类型: b=" + b + " sh=" + sh + " i=" + i + " l=" + l);
        System.out.println("浮点字符布尔: f=" + f + " d=" + d + " c=" + c + " ok=" + ok + " name=" + name);

        // final 常量：不可再赋值
        final int MAX_RETRY = 3;
        // MAX_RETRY = 4;  // 取消注释会编译失败：final 不可改
        System.out.println("final 常量 MAX_RETRY=" + MAX_RETRY);

        // var 是局部变量类型推断，编译期就定型，不是动态类型
        var nums = List.of(1, 2, 3);
        var total = 0;
        for (var n : nums) total += n;
        System.out.println("var 局部变量推断: nums 求和=" + total);

        // 自动装箱 / 拆箱：基本类型 <-> 包装类
        Integer boxed = i;
        int unboxed = boxed;
        System.out.println("自动装箱拆箱: boxed=" + boxed + " unboxed=" + unboxed);

        // 零值：只有字段有，局部变量没有
        var defaults = new Defaults();
        System.out.println("字段零值(局部变量没有零值): i=" + defaults.i + " l=" + defaults.l
                + " d=" + defaults.d + " flag=" + defaults.flag + " char数值=" + (int) defaults.c
                + " s=" + defaults.s + " arr=" + defaults.arr);

        // 范围与类型转换
        System.out.println("int 范围: [" + Integer.MIN_VALUE + ", " + Integer.MAX_VALUE + "]");
        double widened = i;           // 隐式加宽（小类型 -> 大类型）
        int truncated = (int) 3.99;   // 显式窄化，静默截断小数
        System.out.println("类型转换: 隐式加宽=" + widened + " 显式窄化截断=" + truncated);
    }
}
