import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;

public class Main {

    static final Map<String, List<String>> ROLE_PERMS = new HashMap<>();

    static {
        ROLE_PERMS.put("admin", List.of("order:read", "order:write", "user:delete"));
        ROLE_PERMS.put("viewer", List.of("order:read"));
    }

    static String hex(byte[] b) {
        StringBuilder sb = new StringBuilder();
        for (byte x : b) {
            sb.append(String.format("%02x", x));
        }
        return sb.toString();
    }

    static String hashPassword(String password, String salt) {
        try {
            PBEKeySpec spec = new PBEKeySpec(
                    password.toCharArray(), salt.getBytes(StandardCharsets.UTF_8), 1000, 128);
            SecretKeyFactory f = SecretKeyFactory.getInstance("PBKDF2WithHmacSHA256");
            return hex(f.generateSecret(spec).getEncoded());
        } catch (Exception e) {
            throw new IllegalStateException(e);
        }
    }

    static boolean safeEqual(String a, String b) {
        // 固定时间比较，避免按字节提前返回泄露信息
        return MessageDigest.isEqual(
                a.getBytes(StandardCharsets.UTF_8), b.getBytes(StandardCharsets.UTF_8));
    }

    static boolean allowed(String role, String perm) {
        List<String> perms = ROLE_PERMS.get(role);
        return perms != null && perms.contains(perm);
    }

    public static void main(String[] args) {
        String salt = "salt-8f2a";
        String hashed = hashPassword("correct-horse", salt);

        System.out.println("== 1. 加盐哈希 ==");
        System.out.println("salt = " + salt + "，hash = " + hashed.substring(0, 16) + "...");
        System.out.println("同一密码 + 同一盐 结果一致 = "
                + hashed.equals(hashPassword("correct-horse", salt)));
        System.out.println("同一密码 + 不同盐 结果不同 = "
                + !hashed.equals(hashPassword("correct-horse", "salt-other")));

        System.out.println("== 2. 校验 ==");
        System.out.println("正确密码 = " + safeEqual(hashed, hashPassword("correct-horse", salt)));
        System.out.println("错误密码 = " + safeEqual(hashed, hashPassword("wrong", salt)));

        System.out.println("== 3. RBAC ==");
        System.out.println("admin 能 order:write = " + allowed("admin", "order:write"));
        System.out.println("viewer 能 order:write = " + allowed("viewer", "order:write"));
        System.out.println("viewer 能 user:delete = " + allowed("viewer", "user:delete"));

        System.out.println("== 4. 结论 ==");
        System.out.println("密码永远不存明文，也不要可逆加密；盐必须每用户独立随机");
        System.out.println("1000 轮只是演示，生产用 bcrypt/argon2id，成本参数按机器调");
    }
}
