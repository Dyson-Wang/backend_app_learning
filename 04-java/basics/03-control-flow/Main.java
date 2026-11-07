import java.util.List;

public class Main {
    public static void main(String[] args) {
        // if / else if / else
        int score = 78;
        if (score >= 90) {
            System.out.println("分支: 优秀");
        } else if (score >= 60) {
            System.out.println("分支: 及格");
        } else {
            System.out.println("分支: 不及格");
        }

        // 三元表达式：是表达式，有返回值
        String label = score >= 60 ? "通过" : "未通过";
        System.out.println("三元表达式: " + label);

        // switch 表达式（Java 14+ 稳定）：箭头、无穿透、直接返回值、必须 default
        int day = 3;
        String zh = switch (day) {
            case 1 -> "周一";
            case 2 -> "周二";
            case 3, 4, 5 -> "工作日";
            case 6, 7 -> "周末";
            default -> "非法";
        };
        System.out.println("switch 表达式: day=" + day + " -> " + zh);

        // 传统 for
        int sum = 0;
        for (int i = 1; i <= 5; i++) sum += i;
        System.out.println("for 循环 1..5 求和=" + sum);

        // 增强 for：适用于任何 Iterable 与数组
        var fruits = List.of("苹果", "香蕉", "梨");
        StringBuilder sb = new StringBuilder("增强 for:");
        for (var f : fruits) sb.append(" ").append(f);
        System.out.println(sb);

        // while 与 do-while（后者至少执行一次）
        int countdown = 3;
        while (countdown > 0) {
            System.out.println("while 倒数=" + countdown);
            countdown--;
        }
        int once = 0;
        do { once++; } while (once < 1);
        System.out.println("do-while 执行次数=" + once);

        // break / continue 与带标签的 break（精确跳出多层循环）
        outer:
        for (int r = 0; r < 3; r++) {
            for (int col = 0; col < 3; col++) {
                if (r == 1 && col == 1) break outer;
                if (col == 2) continue;      // 跳过本轮末尾输出
                System.out.println("嵌套循环 r=" + r + " col=" + col);
            }
        }
        System.out.println("已用带标签 break 跳出双层循环");
    }
}
