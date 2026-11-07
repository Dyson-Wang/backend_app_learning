public class Main {
    public static void main(String[] args) throws Exception {
        // classpath：JVM 查找 .class / 依赖 jar 的路径
        System.out.println("java.class.path=" + System.getProperty("java.class.path"));
        System.out.println("java.version=" + System.getProperty("java.version"));

        // 类加载器双亲委派：app -> platform -> bootstrap(用 null 表示)
        var loader = Main.class.getClassLoader();
        System.out.println("当前类加载器=" + loader.getName());
        System.out.println("父加载器=" + loader.getParent().getName());
        System.out.println("祖父(bootstrap, 用 null 表示)=" + loader.getParent().getParent());

        // 按名字加载类：依赖在 classpath 上才能成功
        try {
            Class.forName("java.util.ArrayDeque");
            System.out.println("Class.forName 按名加载 JDK 类=成功");
        } catch (ClassNotFoundException e) {
            System.out.println("按名加载失败: " + e.getMessage());
        }

        // 依赖缺失：不在 classpath 的类会抛 ClassNotFoundException
        try {
            Class.forName("com.example.NotOnClasspath");
        } catch (ClassNotFoundException e) {
            System.out.println("依赖不在 classpath -> ClassNotFoundException");
        }

        // 把自身当资源从 classpath 读出来，证明 classpath 就是依赖查找路径
        try (var in = Main.class.getResourceAsStream("/Main.class")) {
            System.out.println("从 classpath 读到 Main.class=" + (in != null));
        }

        // 包与模块（JPMS）：单文件在默认包 / 未命名模块
        System.out.println("包名(空=默认包)='" + Main.class.getPackageName() + "'");
        System.out.println("模块是否具名=" + Main.class.getModule().isNamed());
        System.out.println("包对象=" + Main.class.getPackage());
    }
}
