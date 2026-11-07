import java.util.ArrayList;
import java.util.List;

public class Main {
    enum Code {
        REQUIRED, LENGTH, FORMAT, RANGE
    }

    record FieldError(Code code, String field, String message) {}

    static final class CreateUser {
        String name;
        String email;
        Integer age; // null 表示未提供

        CreateUser(String name, String email, Integer age) {
            this.name = name;
            this.email = email;
            this.age = age;
        }
    }

    static List<FieldError> validate(CreateUser u) {
        List<FieldError> errors = new ArrayList<>();

        String name = u.name == null ? "" : u.name.trim();
        if (name.isEmpty()) {
            errors.add(new FieldError(Code.REQUIRED, "name", "name 不能为空"));
        } else if (name.length() < 2 || name.length() > 20) {
            errors.add(new FieldError(Code.LENGTH, "name", "name 长度必须在 2-20 之间"));
        }

        String email = u.email == null ? "" : u.email;
        if (email.isEmpty()) {
            errors.add(new FieldError(Code.REQUIRED, "email", "email 不能为空"));
        } else if (!email.contains("@")) {
            errors.add(new FieldError(Code.FORMAT, "email", "email 格式不正确"));
        }

        if (u.age == null) {
            errors.add(new FieldError(Code.REQUIRED, "age", "age 不能为空"));
        } else if (u.age < 0 || u.age > 150) {
            errors.add(new FieldError(Code.RANGE, "age", "age 必须在 0-150 之间"));
        }

        return errors;
    }

    static void report(String label, CreateUser u) {
        List<FieldError> errors = validate(u);
        System.out.println(label + "：" + errors.size() + " 个错误");
        for (FieldError e : errors) {
            System.out.println("  [" + e.code() + "] " + e.field() + ": " + e.message());
        }
    }

    public static void main(String[] args) {
        report("全空", new CreateUser(null, null, null));
        report("短名 + 坏邮箱 + 超龄", new CreateUser("a", "bad", 200));
        report("合法", new CreateUser("ada", "ada@example.com", 36));
        System.out.println("结论：一次返回全部错误，前端可逐字段标红；code 与 message 分离便于国际化");
    }
}
