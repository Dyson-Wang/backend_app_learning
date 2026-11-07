public class Main {
    private static int sideCalls = 0;

    // 有副作用的方法，用来观察短路时右边是否被跳过
    static boolean side() {
        sideCalls++;
        return true;
    }

    public static void main(String[] args) {
        // 算术、整除与优先级
        int a = 7, b = 3;
        System.out.println("算术: a+b=" + (a + b) + " a/b(整除)=" + (a / b) + " a%b=" + (a % b));
        System.out.println("优先级: 2+3*4=" + (2 + 3 * 4) + " (2+3)*4=" + ((2 + 3) * 4));
        System.out.println("浮点除: 7.0/2=" + (7.0 / 2));

        // 比较与逻辑，含短路求值
        System.out.println("比较: a>b=" + (a > b) + " a==b=" + (a == b));
        sideCalls = 0;
        boolean r1 = false && side();   // 左边 false，右边不执行
        System.out.println("短路 &&: r1=" + r1 + " side 调用=" + sideCalls);
        boolean r2 = true || side();    // 左边 true，右边不执行
        System.out.println("短路 ||: r2=" + r2 + " side 调用=" + sideCalls);
        boolean r3 = true & side();     // 单个 & 不短路，右边照样执行
        System.out.println("非短路 &: r3=" + r3 + " side 调用=" + sideCalls);

        // 位运算与移位
        int x = 0b1100, y = 0b1010;
        System.out.println("位与 x&y=" + Integer.toBinaryString(x & y));
        System.out.println("位或 x|y=" + Integer.toBinaryString(x | y)
                + " 异或 x^y=" + Integer.toBinaryString(x ^ y));
        System.out.println("移位: 1<<4=" + (1 << 4) + " 12>>2=" + (12 >> 2) + " -1>>>28=" + (-1 >>> 28));

        // 赋值与复合赋值；复合赋值自带隐式窄化 cast
        int c = 10;
        c += 5;
        c *= 2;
        System.out.println("复合赋值: c=" + c);
        byte small = 0;
        small += 200;                    // 等价于 small = (byte)(small + 200)
        System.out.println("byte 复合赋值隐式窄化: small=" + small);
        // small = small + 200;         // 这样会编译失败：需要显式 (byte) 强转

        // 整数溢出：二进制补码静默回绕，不抛异常
        int max = Integer.MAX_VALUE;
        System.out.println("整数溢出: Integer.MAX_VALUE+1=" + (max + 1));

        // 三元表达式
        String grade = a > b ? "a 更大" : "b 更大";
        System.out.println("三元表达式: " + grade);
    }
}
