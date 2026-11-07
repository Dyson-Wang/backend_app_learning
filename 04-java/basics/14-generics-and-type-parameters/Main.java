import java.util.ArrayList;
import java.util.List;

public class Main {
    // 泛型类 + 上界约束
    static class Box<T extends Comparable<T>> {
        private final T value;
        Box(T value) { this.value = value; }
        T get() { return value; }
    }

    // 泛型方法
    static <T> T first(List<T> xs) { return xs.get(0); }

    // 上界约束：T 必须是 Number 或其子类
    static <T extends Number> double asDouble(T n) { return n.doubleValue(); }

    // 通配符 ? extends：生产者，只能安全地「读」（PECS 的 P）
    static double total(List<? extends Number> xs) {
        double s = 0;
        for (Number n : xs) s += n.doubleValue();
        return s;
    }

    // 通配符 ? super：消费者，只能安全地「写」（PECS 的 C）
    static void addInts(List<? super Integer> xs) {
        xs.add(1);
        xs.add(2);
    }

    public static void main(String[] args) {
        // 泛型方法 + 类型推断
        var s = first(List.of("a", "b"));
        var n = first(List.of(10, 20));
        System.out.println("泛型方法类型推断: 字符串=" + s + " 整数翻倍=" + (n * 2));

        // 泛型类
        Box<String> box = new Box<>("dyson");
        System.out.println("泛型类 Box<String>: " + box.get());

        // 上界约束
        System.out.println("上界 <T extends Number>: " + asDouble(3.5) + " " + asDouble(7));

        // 通配符：? extends 只读、? super 只写
        System.out.println("? extends Number 求和(List<Integer>)=" + total(List.of(1, 2, 3)));
        System.out.println("? extends Number 求和(List<Double>)=" + total(List.of(1.5, 2.5)));
        List<Number> nums = new ArrayList<>();
        addInts(nums);
        System.out.println("? super Integer 写入后=" + nums);

        // 类型擦除实证：运行期根本没有 List<String> / List<Integer> 之分
        List<String> strings = new ArrayList<>(List.of("x"));
        List<Integer> ints = new ArrayList<>(List.of(1));
        System.out.println("类型擦除: strings.getClass()==ints.getClass() -> "
                + (strings.getClass() == ints.getClass()));
        System.out.println("擦除后的运行期类=" + strings.getClass().getName());
        // List<String>.class 是非法的（不能对参数化类型写类字面量），擦除后只有 List.class
        System.out.println("List.class=" + List.class.getName());
        System.out.println("泛型信息只存在于编译期，运行期被擦除为原始类型 / 上界");
    }
}
