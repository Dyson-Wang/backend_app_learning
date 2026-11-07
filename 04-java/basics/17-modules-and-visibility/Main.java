public class Main {
    // 可见性四级：public / protected / 包私有 / private
    public static class PublicThing {
        public String pub = "public";
        protected String prot = "protected";
        String pkg = "包私有";
        private String priv = "private";

        private String secret() { return "只有本类内部能调用"; }
        String readPrivate() { return priv + " / " + secret(); }  // 本类内部可访问 private
    }

    // 两个类互相引用：Java 允许编译单元内的循环依赖
    static class A {
        B b;
        A(B b) { this.b = b; }
        String viaB() { return "A 通过 B 拿到 " + b.name(); }
    }
    static class B {
        A a;
        String name() { return "B"; }
        void attach(A a) { this.a = a; }
    }

    public static void main(String[] args) {
        var t = new PublicThing();
        System.out.println("可见性四级(由宽到窄): " + t.pub + " > " + t.prot + " > " + t.pkg + " > private");
        System.out.println("private 通过本类方法暴露: " + t.readPrivate());
        System.out.println("目录即包: 包名='" + Main.class.getPackageName() + "'（空=默认包）");

        // 命名冲突：java.util.Date 与 java.sql.Date 同名，Java 没有 import 别名，只能用全限定名
        System.out.println("同名类用全限定名区分: " + java.util.Date.class.getName()
                + " vs " + java.sql.Date.class.getName());

        // 循环依赖：A 依赖 B，B 又持有 A；同一编译单元内 javac 能解析
        var a = new A(new B());
        B bb = a.b;
        bb.attach(a);
        System.out.println("循环依赖: " + a.viaB() + ", B.a 指向同一 A=" + (bb.a == a));
    }
}
