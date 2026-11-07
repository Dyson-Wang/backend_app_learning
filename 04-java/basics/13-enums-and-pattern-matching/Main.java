public class Main {
    // 枚举：可带字段、构造器、方法
    enum Level {
        LOW(1), MID(2), HIGH(3);
        private final int code;
        Level(int code) { this.code = code; }
        int code() { return code; }
    }

    // sealed 接口 + permits：只允许下面两个实现，构成封闭的 ADT
    sealed interface Shape permits Circle, Rect {}

    // record 天然 final，正好作为「带数据的枚举变体」
    record Circle(double r) implements Shape {}
    record Rect(double w, double h) implements Shape {}

    // 用 instanceof 模式匹配做「解构」（Java 16+ 稳定）
    // 注意：switch 的模式匹配在 17/21 仍是 preview，这里绝不使用
    static double area(Shape s) {
        if (s instanceof Circle c) return Math.PI * c.r() * c.r();
        if (s instanceof Rect r) return r.w() * r.h();
        throw new IllegalStateException("sealed 已覆盖所有实现，理论上到不了这里");
    }

    static String describe(Shape s) {
        if (s instanceof Circle c) return "圆 半径=" + c.r();
        if (s instanceof Rect r) return "矩形 " + r.w() + "x" + r.h();
        throw new IllegalStateException("unreachable");
    }

    public static void main(String[] args) {
        // 枚举取值 + switch 表达式（对枚举是稳定特性）
        for (var lv : Level.values()) {
            String zh = switch (lv) {
                case LOW -> "低";
                case MID -> "中";
                case HIGH -> "高";
                default -> "未知";
            };
            System.out.println("枚举: " + lv + " code=" + lv.code() + " 中文=" + zh);
        }
        System.out.println("Level.valueOf(\"HIGH\")=" + Level.valueOf("HIGH")
                + " ordinal=" + Level.HIGH.ordinal());

        // 构造 ADT 实例并用 instanceof 模式匹配解构
        Shape[] shapes = { new Circle(2.0), new Rect(3.0, 4.0) };
        for (var s : shapes) {
            System.out.println("模式匹配: " + describe(s) + " 面积=" + String.format("%.2f", area(s)));
        }

        // 穷尽性：sealed 让编译期知道只有 Circle/Rect 两种实现
        // 但 if/else 链编译器不做穷尽检查，所以末尾必须自己 throw 兜底
        System.out.println("穷尽性靠 sealed 约束实现集合，而非 switch 的编译器检查");
    }
}
