import java.nio.charset.StandardCharsets;
import java.util.Arrays;

public class Main {
    public static void main(String[] args) {
        // String 不可变：concat 返回新对象，原串不变
        String s = "后端";
        String s2 = s.concat("路线图");
        System.out.println("不可变性: 原串=" + s + " concat 结果=" + s2 + " 同一对象=" + (s == s2));

        // 拼接与格式化
        String joined = String.join("/", "2024", "01", "02");
        String formatted = String.format("%s 有 %d 课", "Java", 24);
        System.out.println("拼接: " + joined + " | 格式化: " + formatted);

        // 大小写 / 查找 / 切分
        String raw = "Hello,Backend,Roadmap";
        System.out.println("大小写: " + raw.toLowerCase() + " 含 Backend=" + raw.contains("Backend"));
        System.out.println("切分: " + Arrays.toString(raw.split(",")));

        // 码点 vs UTF-16 码元 vs UTF-8 字节：用一个非 BMP 的 emoji（占 2 个 char）
        String emoji = new String(Character.toChars(0x1F600));  // U+1F600 😀
        System.out.println("emoji=" + emoji + " UTF-16长度=" + emoji.length()
                + " 码点数=" + emoji.codePointCount(0, emoji.length())
                + " UTF-8字节数=" + emoji.getBytes(StandardCharsets.UTF_8).length);
        System.out.print("按码点遍历: ");
        emoji.codePoints().forEach(cp -> System.out.print(Integer.toHexString(cp) + " "));
        System.out.println();

        // 编码 / 解码：必须用同一字符集
        byte[] bytes = "中文ABC".getBytes(StandardCharsets.UTF_8);
        String decoded = new String(bytes, StandardCharsets.UTF_8);
        System.out.println("UTF-8 往返: 字节数=" + bytes.length + " 还原=" + decoded);

        // 转义字符串与 text block（Java 的「原始字符串」近似物）
        String escaped = "第一行\n第二行\t带制表";
        String block = """
                第一行
                第二行
                """.strip();
        System.out.println("转义字符串=" + escaped.replace("\n", "\\n").replace("\t", "\\t"));
        System.out.println("text block 行数=" + block.lines().count());

        // strip / repeat
        System.out.println("strip='" + "  两头  ".strip() + "' repeat=" + "ab ".repeat(3));
    }
}
