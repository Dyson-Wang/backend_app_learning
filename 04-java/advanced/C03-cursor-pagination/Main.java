import java.util.ArrayList;
import java.util.Comparator;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class Main {

    static class Item {
        final int id;
        final int score;

        Item(int id, int score) {
            this.id = id;
            this.score = score;
        }
    }

    static final Comparator<Item> ORDER =
            Comparator.comparingInt((Item i) -> -i.score).thenComparingInt(i -> i.id);

    static List<Item> ordered(List<Item> items) {
        List<Item> out = new ArrayList<>(items);
        out.sort(ORDER);
        return out;
    }

    static List<Item> offsetPage(List<Item> items, int offset, int limit) {
        List<Item> s = ordered(items);
        if (offset >= s.size()) {
            return List.of();
        }
        return new ArrayList<>(s.subList(offset, Math.min(offset + limit, s.size())));
    }

    static List<Item> cursorPage(List<Item> items, Item after, int limit) {
        List<Item> out = new ArrayList<>();
        for (Item it : ordered(items)) {
            if (after != null) {
                boolean later = it.score < after.score
                        || (it.score == after.score && it.id > after.id);
                if (!later) {
                    continue;
                }
            }
            out.add(it);
            if (out.size() == limit) {
                break;
            }
        }
        return out;
    }

    static List<Integer> ids(List<Item> items) {
        List<Integer> out = new ArrayList<>();
        for (Item i : items) {
            out.add(i.id);
        }
        return out;
    }

    static int overlap(List<Item> a, List<Item> b) {
        Set<Integer> seen = new HashSet<>();
        for (Item i : a) {
            seen.add(i.id);
        }
        int n = 0;
        for (Item i : b) {
            if (seen.contains(i.id)) {
                n++;
            }
        }
        return n;
    }

    public static void main(String[] args) {
        List<Item> base = new ArrayList<>();
        for (int i = 0; i < 10; i++) {
            base.add(new Item(i + 1, 100 - i * 5));
        }

        System.out.println("== 1. OFFSET 分页：翻页期间有新数据插入 ==");
        List<Item> p1 = offsetPage(base, 0, 3);
        List<Item> grown = new ArrayList<>(base);
        grown.add(new Item(11, 97));
        List<Item> p2 = offsetPage(grown, 3, 3);
        System.out.println("第一页 = " + ids(p1));
        System.out.println("插入一条 Score=97 后，第二页 = " + ids(p2)
                + "，与第一页重复 " + overlap(p1, p2) + " 条");

        System.out.println("== 2. 游标分页：用 (score, id) 做锚点 ==");
        List<Item> c1 = cursorPage(base, null, 3);
        Item last = c1.get(c1.size() - 1);
        List<Item> c2 = cursorPage(grown, last, 3);
        System.out.println("第一页 = " + ids(c1)
                + "，锚点 = (score=" + last.score + ", id=" + last.id + ")");
        System.out.println("插入同一条后，第二页 = " + ids(c2)
                + "，与第一页重复 " + overlap(c1, c2) + " 条");

        System.out.println("== 3. 结论 ==");
        System.out.println("OFFSET 按位置翻页，前面插入数据会让整体后移，出现重复或漏读");
        System.out.println("游标按值翻页，锚点不变结果就不变；代价是不能跳页");
    }
}
