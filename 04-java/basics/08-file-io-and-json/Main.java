import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;

public class Main {
    // 教学用极简 JSON：只抽取字符串字段，生产请用 Jackson / Gson，别手写
    static String strField(String json, String key) {
        int k = json.indexOf("\"" + key + "\"");
        if (k < 0) return null;
        int colon = json.indexOf(':', k);
        int first = json.indexOf('"', colon);
        int last = json.indexOf('"', first + 1);
        if (colon < 0 || first < 0 || last < 0) return null;
        return json.substring(first + 1, last);
    }

    public static void main(String[] args) throws Exception {
        // Files.writeString / readString（Java 11+），全程显式 UTF-8
        Path file = Files.createTempFile("brm-", ".txt");
        try {
            String text = "第一行：dyson\n第二行：后端路线图\n";
            Files.writeString(file, text, StandardCharsets.UTF_8);
            System.out.println("写入临时文件完成, 字节数=" + Files.readAllBytes(file).length);

            String back = Files.readString(file, StandardCharsets.UTF_8);
            System.out.println("按行读回: " + back.lines().count() + " 行");
            System.out.println("首行=" + back.lines().findFirst().orElse(""));
            System.out.println("属性: 是普通文件=" + Files.isRegularFile(file) + " 大小=" + Files.size(file));
        } finally {
            Files.deleteIfExists(file);
            System.out.println("清理临时文件=" + !Files.exists(file));
        }

        // 手写 JSON 序列化：text block 让多行字符串更干净
        String json = """
                {"id":7,"name":"alice","role":"admin"}
                """.strip();
        System.out.println("构造 JSON=" + json);
        System.out.println("解析 name=" + strField(json, "name") + " role=" + strField(json, "role"));

        // 转义：JSON 字符串里的引号要写成 \"
        String quoted = "他说\"你好\"";
        System.out.println("转义演示: " + quoted.replace("\"", "\\\""));
    }
}
