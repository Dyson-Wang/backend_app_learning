import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;
import java.lang.reflect.Field;
import java.lang.reflect.Method;

public class Main {
    // 自定义注解：RetentionPolicy.RUNTIME 才能在运行期用反射读到
    @Retention(RetentionPolicy.RUNTIME)
    @Target(ElementType.METHOD)
    @interface Tag {
        String value();
        String author() default "dyson";
    }

    static class Box {
        private int count = 7;
        public String hello(String who) { return "你好 " + who; }
    }

    @Tag("需要鉴权")
    public void secureOp() { }

    public static void main(String[] args) throws Exception {
        // 1) 运行时类型信息（RTTI）
        Class<?> clazz = Box.class;
        System.out.println("类名=" + clazz.getName() + " 父类=" + clazz.getSuperclass().getSimpleName());
        System.out.println("声明的方法数=" + clazz.getDeclaredMethods().length);

        // 2) 按名字访问字段 / 方法，并动态调用
        Box box = new Box();
        Field f = clazz.getDeclaredField("count");
        f.setAccessible(true);                     // 突破 private
        System.out.println("按名读取 private 字段 count=" + f.getInt(box));
        f.setInt(box, 100);
        System.out.println("反射修改后 count=" + f.getInt(box));

        Method hi = clazz.getDeclaredMethod("hello", String.class);
        Object result = hi.invoke(box, "dyson");
        System.out.println("按名动态调用 hello(dyson)=" + result);

        // 3) 读取运行时注解
        Method op = Main.class.getDeclaredMethod("secureOp");
        Tag tag = op.getAnnotation(Tag.class);
        System.out.println("读到注解: value=" + tag.value() + " author=" + tag.author());
        System.out.println("注解类型=" + tag.annotationType().getSimpleName() + "（RUNTIME 保留）");

        // 4) 代价与边界
        System.out.println("反射绕过编译期检查、有性能开销、受模块 opens 限制；谨慎用于框架");
        System.out.println("Lombok / APT 属于编译期代码生成；Java 没有宏，运行期靠反射、编译期靠注解处理器");
    }
}
