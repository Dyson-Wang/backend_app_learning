import java.util.List;
import java.util.function.Function;
import java.util.function.IntUnaryOperator;
import java.util.function.Predicate;

public class Main {
    // 高阶函数：返回一个函数，闭包捕获了参数 n（effectively final）
    static IntUnaryOperator adder(int n) {
        return x -> x + n;
    }

    // 高阶函数：接收函数并连续应用两次
    static int applyTwice(IntUnaryOperator f, int v) {
        return f.applyAsInt(f.applyAsInt(v));
    }

    // this 绑定：lambda 里的 this 是词法捕获的外层对象，不像 JS 普通函数动态绑定
    static class Self {
        private final String name;
        Self(String name) { this.name = name; }
        Runnable callback() {
            return () -> System.out.println("  callback 里的 this.name=" + this.name);
        }
    }

    public static void main(String[] args) {
        // 函数作为值（方法引用）
        Function<String, Integer> len = String::length;
        System.out.println("函数作为值: 长度=" + len.apply("dyson"));

        // 返回函数 + 捕获
        var plus5 = adder(5);
        System.out.println("闭包捕获 n=5: plus5(10)=" + plus5.applyAsInt(10));
        System.out.println("高阶函数 applyTwice(+5, 0)=" + applyTwice(plus5, 0));

        // Predicate 作为过滤条件（高阶函数）
        Predicate<Integer> even = x -> x % 2 == 0;
        System.out.println("Predicate 过滤偶数: " + List.of(1, 2, 3, 4, 5).stream().filter(even).toList());

        // 捕获变量必须是 effectively final：被重新赋值后就无法捕获
        int base = 10;
        // base = 20;  // 取消注释后，下面这行会编译失败
        IntUnaryOperator f = x -> x + base;
        System.out.println("捕获 effectively final base=10: " + f.applyAsInt(1));

        // 捕获可变状态的坑：把「引用」保持 final，用数组/原子类改内部状态
        int[] holder = {0};
        Runnable tick = () -> holder[0]++;
        tick.run();
        tick.run();
        System.out.println("可变捕获(final 数组引用)累加后=" + holder[0]);

        // 生命周期：闭包捕获值/引用副本，可脱离创建它的栈帧存活
        var maker = adder(100);
        System.out.println("闭包脱离原栈帧仍可用: " + maker.applyAsInt(1));

        // this 绑定
        new Self("dyson").callback().run();
        System.out.println("lambda 的 this 词法捕获外层对象，不随调用方式改变");
    }
}
