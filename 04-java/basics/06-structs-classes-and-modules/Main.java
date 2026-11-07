public class Main {
    // 普通类：封装、私有字段、构造器校验
    static class BankAccount {
        private String owner;
        private int balance;

        BankAccount(String owner, int opening) {
            if (opening < 0) throw new IllegalArgumentException("开户金额不能为负");
            this.owner = owner;
            this.balance = opening;
        }
        void deposit(int amount) {
            if (amount <= 0) throw new IllegalArgumentException("存款必须为正");
            balance += amount;
        }
        int getBalance() { return balance; }
        String getOwner() { return owner; }
    }

    // record：不可变「结构体」，自动生成访问器 / equals / hashCode / toString
    record Point(int x, int y) {
        Point {                       // 紧凑构造器，做校验
            if (x < 0 || y < 0) throw new IllegalArgumentException("坐标不能为负");
        }
        int manhattan() { return Math.abs(x) + Math.abs(y); }
    }

    // interface：可带 default 与 static 方法；抽象方法隐式 public
    interface Greeter {
        String greet(String who);
        default String shout(String who) { return greet(who).toUpperCase(); }
        static String lang() { return "zh-CN"; }
    }

    // abstract class + 继承 + @Override
    abstract static class Animal {
        abstract String speak();
        String kind() { return "动物"; }
    }
    static class Dog extends Animal {
        @Override String speak() { return "汪汪"; }
    }

    // 实现接口：接口方法隐式 public，这里必须显式写 public
    static class FormalGreeter implements Greeter {
        public String greet(String who) { return "你好, " + who; }
    }

    public static void main(String[] args) {
        var acct = new BankAccount("dyson", 100);
        acct.deposit(50);
        System.out.println("类与封装: " + acct.getOwner() + " 余额=" + acct.getBalance());

        try {
            new BankAccount("x", -1);
        } catch (IllegalArgumentException e) {
            System.out.println("构造器校验: " + e.getMessage());
        }

        var p = new Point(3, 4);
        var q = new Point(3, 4);
        System.out.println("record: " + p + " 曼哈顿距离=" + p.manhattan());
        System.out.println("record 值相等(equals 自动生成): " + p.equals(q));

        Greeter g = new FormalGreeter();
        System.out.println("接口 default: " + g.shout("dyson") + " 静态方法=" + Greeter.lang());

        Animal a = new Dog();
        System.out.println("继承与多态: " + a.kind() + " 叫声=" + a.speak());

        System.out.println("可见性: private 字段只能通过 getter 访问，default 方法可被实现复用");
    }
}
