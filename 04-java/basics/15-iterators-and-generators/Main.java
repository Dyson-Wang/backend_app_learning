import java.util.Iterator;
import java.util.NoSuchElementException;
import java.util.stream.Stream;

public class Main {
    // 自定义一次性迭代器：模拟「生成器」，每次 next() 才产出一个值（惰性）
    static final class Counter implements Iterator<Integer> {
        private int cur;
        private final int end;
        Counter(int start, int end) { this.cur = start; this.end = end; }
        @Override public boolean hasNext() { return cur < end; }
        @Override public Integer next() {
            if (!hasNext()) throw new NoSuchElementException();
            return cur++;
        }
    }

    // 消费迭代器（一次性）
    static int drain(Iterator<Integer> it) {
        int s = 0;
        while (it.hasNext()) s += it.next();
        return s;
    }

    public static void main(String[] args) {
        // 自定义迭代器 + 惰性：只在 next() 时推进
        var counter = new Counter(1, 6);
        System.out.println("自定义迭代器首个值=" + counter.next() + " 之后 hasNext=" + counter.hasNext());
        System.out.println("剩余元素求和=" + drain(counter));

        // 一次性消费陷阱：drain 过后再消费得到 0
        System.out.println("一次性消费陷阱: 再次 drain=" + drain(counter));

        // Stream 惰性管道：无限序列 + map / filter / take(limit)
        int pipelineSum = Stream.iterate(1, n -> n + 1)
                .map(n -> n * 2)
                .filter(n -> n % 3 == 0)
                .limit(4)
                .mapToInt(Integer::intValue)
                .sum();
        System.out.println("无限序列取前 4 项(map/filter/take)之和=" + pipelineSum);

        // chain：拼接两个流
        var chained = Stream.concat(Stream.of(1, 2), Stream.of(3, 4)).toList();
        System.out.println("chain 拼接=" + chained);

        // 流也只能消费一次
        Stream<Integer> oneShot = Stream.of(1, 2, 3);
        oneShot.forEach(x -> {});
        try {
            oneShot.forEach(x -> {});
        } catch (IllegalStateException e) {
            System.out.println("Stream 二次消费 -> IllegalStateException");
        }
    }
}
