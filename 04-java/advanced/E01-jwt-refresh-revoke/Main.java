import java.nio.charset.StandardCharsets;
import java.util.HashSet;
import java.util.Set;
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class Main {

    static final byte[] SECRET = "s3cr3t".getBytes(StandardCharsets.UTF_8);

    static String hex(byte[] b) {
        StringBuilder sb = new StringBuilder();
        for (byte x : b) {
            sb.append(String.format("%02x", x));
        }
        return sb.toString();
    }

    static String sign(String payload) {
        try {
            Mac mac = Mac.getInstance("HmacSHA256");
            mac.init(new SecretKeySpec(SECRET, "HmacSHA256"));
            byte[] out = mac.doFinal(("HS256." + payload).getBytes(StandardCharsets.UTF_8));
            return hex(out).substring(0, 32);
        } catch (Exception e) {
            throw new IllegalStateException(e);
        }
    }

    static String issue(String sub, String jti, long exp) {
        String payload = "sub=" + sub + ";jti=" + jti + ";exp=" + exp;
        return "HS256." + payload + "." + sign(payload);
    }

    static String field(String payload, String key) {
        for (String kv : payload.split(";")) {
            if (kv.startsWith(key + "=")) {
                return kv.substring(key.length() + 1);
            }
        }
        return "";
    }

    static String[] verify(String token, long now, Set<String> revoked) {
        String[] parts = token.split("\\.");
        if (parts.length != 3 || !parts[0].equals("HS256")) {
            return new String[]{"", "bad format or alg"};
        }
        if (!parts[2].equals(sign(parts[1]))) {
            return new String[]{"", "signature mismatch"};
        }
        if (now >= Long.parseLong(field(parts[1], "exp"))) {
            return new String[]{"", "token expired"};
        }
        if (revoked.contains(field(parts[1], "jti"))) {
            return new String[]{"", "token revoked"};
        }
        return new String[]{field(parts[1], "sub"), ""};
    }

    public static void main(String[] args) {
        long now = 1000L;
        Set<String> revoked = new HashSet<>();

        System.out.println("== 1. 签发 ==");
        String access = issue("u1", "jti-a1", now + 60);
        String refresh = issue("u1", "jti-r1", now + 86400);
        System.out.println("access.jti = jti-a1，exp = +60s");
        System.out.println("refresh.jti = jti-r1，exp = +86400s");

        System.out.println("== 2. 校验正常 token ==");
        String[] r = verify(access, now, revoked);
        System.out.println("ok = " + r[1].isEmpty() + "，sub = " + r[0]);

        System.out.println("== 3. 篡改 payload ==");
        r = verify(access.replace("sub=u1", "sub=u2"), now, revoked);
        System.out.println("ok = " + r[1].isEmpty() + "，reason = " + r[1]);

        System.out.println("== 4. 刷新：轮换 refresh token ==");
        revoked.add("jti-r1");
        issue("u1", "jti-r2", now + 86400);
        System.out.println("新 refresh.jti = jti-r2，旧 jti-r1 进黑名单");
        r = verify(refresh, now, revoked);
        System.out.println("用旧 refresh 再刷：ok = " + r[1].isEmpty() + "，reason = " + r[1]);

        System.out.println("== 5. 吊销 access token ==");
        revoked.add("jti-a1");
        r = verify(access, now, revoked);
        System.out.println("ok = " + r[1].isEmpty() + "，reason = " + r[1]);

        System.out.println("== 6. 结论 ==");
        System.out.println("签名只保证没被篡改，吊销必须靠 jti 黑名单（或短 TTL + 版本号）");
        System.out.println("refresh 必须轮换：旧的用过一次就作废，否则泄露后可长期续命");
    }
}
