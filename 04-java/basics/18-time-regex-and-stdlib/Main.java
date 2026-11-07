import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.Period;
import java.time.ZoneId;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.Collections;
import java.util.List;
import java.util.Random;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Main {
    public static void main(String[] args) {
        // 时间点、格式化与解析：用固定时间，不打印当前绝对时间
        LocalDateTime t = LocalDateTime.of(2024, 1, 2, 3, 4, 5);
        DateTimeFormatter fmt = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
        System.out.println("格式化: " + t.format(fmt));
        LocalDateTime parsed = LocalDateTime.parse("2024-06-01 12:00:00", fmt);
        System.out.println("解析回来: " + parsed);

        // 时长 / 日期差
        LocalDateTime later = t.plusDays(10).plusHours(5);
        System.out.println("Duration 小时差=" + Duration.between(t, later).toHours());
        System.out.println("Period 日期差(天)=" + Period.between(t.toLocalDate(), LocalDate.of(2024, 2, 12)).getDays());

        // 时区转换：固定时刻，无 DST 干扰
        ZonedDateTime shanghai = t.atZone(ZoneId.of("Asia/Shanghai"));
        ZonedDateTime utc = shanghai.withZoneSameInstant(ZoneId.of("UTC"));
        System.out.println("上海 " + shanghai.getHour() + " 点对应 UTC " + utc.getHour() + " 点");

        // 正则：匹配、分组捕获、替换
        String text = "订单 A-123 与 B-456 已创建";
        Pattern p = Pattern.compile("([A-Z])-(\\d+)");
        Matcher m = p.matcher(text);
        StringBuilder found = new StringBuilder("正则分组:");
        while (m.find()) found.append(" [").append(m.group(1)).append(":").append(m.group(2)).append("]");
        System.out.println(found);
        System.out.println("替换数字为 #: " + m.replaceAll("$1-#"));
        System.out.println("Pattern.matches(\\d+, 12345)=" + Pattern.matches("\\d+", "12345"));

        // 固定种子随机数：两个相同种子序列一致，只比较相等性，输出可复现
        Random r1 = new Random(42L), r2 = new Random(42L);
        System.out.println("同种子随机数序列一致=" + (r1.nextInt(100) == r2.nextInt(100)
                && r1.nextInt(100) == r2.nextInt(100)));

        // 排序、二分、频次、最值
        var nums = new ArrayList<>(List.of(5, 1, 4, 2, 3));
        Collections.sort(nums);
        System.out.println("排序=" + nums + " 最大值=" + Collections.max(nums)
                + " 二分查找 4=" + Collections.binarySearch(nums, 4));
        System.out.println("频次=" + Collections.frequency(nums, 3) + " 数组=" + Arrays.toString(nums.toArray()));

        // 常用工具：Base64、Math
        String enc = Base64.getEncoder().encodeToString("dyson".getBytes(StandardCharsets.UTF_8));
        String dec = new String(Base64.getDecoder().decode(enc), StandardCharsets.UTF_8);
        System.out.println("Base64 往返一致=" + dec.equals("dyson"));
        System.out.println("Math: 四舍五入=" + Math.round(3.6) + " 绝对值=" + Math.abs(-7));
    }
}
