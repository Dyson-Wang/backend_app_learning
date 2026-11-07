import java.lang.ref.WeakReference;

public class Main {
    // 确定性清理：AutoCloseable 的 close() 一定在离开 try 块时执行
    static class Resource implements AutoCloseable {
        private final String name;
        Resource(String name) { this.name = name; }
        void use() { System.out.println("  使用资源 " + name + "（对象在堆上）"); }
        @Override public void close() { System.out.println("  close(): 确定性释放 " + name); }
    }

    // 引用 vs 值：对象通过引用共享
    static class Holder {
        int n;
        Holder(int n) { this.n = n; }
    }

    public static void main(String[] args) {
        // 值拷贝 vs 引用拷贝
        int x = 10;
        int y = x;            // 基本类型：值拷贝
        y = 99;
        System.out.println("值拷贝: x=" + x + " y=" + y + "（改 y 不影响 x）");

        Holder h1 = new Holder(1);
        Holder h2 = h1;       // 引用拷贝，指向同一堆对象
        h2.n = 42;
        System.out.println("引用拷贝: h1.n=" + h1.n + " h2.n=" + h2.n + " 同一对象=" + (h1 == h2));

        // == 比较引用，equals 比较内容；Integer 缓存是经典「身份陷阱」
        Integer a = 127, b = 127;
        Integer c = 128, d = 128;
        System.out.println("== 陷阱: 127(缓存内) a==b=" + (a == b) + " , 128(缓存外) c==d=" + (c == d));

        // 栈与堆、GC Roots、逃逸
        System.out.println("栈存局部变量/引用，堆存对象实体；对象经 GC Roots 可达链引用才存活");

        // WeakReference：弱引用不阻止 GC；强引用在时一定能取到
        Object obj = new Object();
        WeakReference<Object> weak = new WeakReference<>(obj);
        System.out.println("强引用存在时 weak.get() 非空=" + (weak.get() != null));
        obj = null;           // 去掉强引用
        System.gc();          // 仅「建议」GC，不保证立即回收，也不读取其结果
        System.out.println("System.gc() 只是建议，是否/何时回收都不确定，故不作结论");

        // try-with-resources：确定性清理，对应 RAII / defer
        try (var r = new Resource("db-conn")) {
            r.use();
        }

        // 常见内存陷阱：静态集合长期持有引用 => 泄漏
        System.out.println("常见陷阱: static 集合长期持有对象引用不释放 => 内存泄漏");
        System.out.println("GC Roots 包括: 栈帧局部变量、静态字段、活动线程、JNI 引用等");
    }
}
