import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Queue;
import java.util.Set;

public class Main {
    public static void main(String[] args) {
        // List：有序、可重复
        List<String> list = new ArrayList<>(List.of("香蕉", "苹果", "梨", "苹果"));
        System.out.println("List 初始=" + list + " 索引 1=" + list.get(1));

        // Iterator：想边遍历边删，必须用它自己的 remove()，且要先 next()
        Iterator<String> it = list.iterator();
        while (it.hasNext()) {
            if ("苹果".equals(it.next())) it.remove();
        }
        System.out.println("Iterator 删除所有苹果后=" + list);

        // for-each 中直接改集合会触发快速失败（≥3 个元素、删非末尾元素才会再走到 next 检查）
        var modList = new ArrayList<>(List.of("a", "b", "c"));
        try {
            for (String s : modList) {
                if ("a".equals(s)) modList.remove(s);
            }
        } catch (java.util.ConcurrentModificationException e) {
            System.out.println("for-each 中改集合 -> ConcurrentModificationException（快速失败）");
        }

        // Set：去重
        Set<String> set = new HashSet<>(List.of("a", "b", "a", "c"));
        System.out.println("Set 去重: size=" + set.size() + " 含a=" + set.contains("a"));

        // Map：按固定键序打印，避免 HashMap 顺序影响输出可复现性
        Map<String, Integer> map = new HashMap<>();
        map.put("张三", 30);
        map.put("李四", 25);
        map.put("王五", 28);
        for (String key : List.of("张三", "李四", "王五")) {
            System.out.println("  Map: " + key + " -> " + map.get(key));
        }

        // Queue / Deque：FIFO
        Queue<String> queue = new ArrayDeque<>();
        queue.offer("任务1");
        queue.offer("任务2");
        System.out.println("Queue FIFO: poll=" + queue.poll() + " peek=" + queue.peek());

        // 排序 + 二分查找
        var nums = new ArrayList<>(List.of(5, 1, 4, 2, 3));
        nums.sort(Comparator.naturalOrder());
        System.out.println("排序后=" + nums + " 二分查找 4 的位置=" + java.util.Collections.binarySearch(nums, 4));

        // 不可变集合：List.of 不支持修改
        var readOnly = List.of("x", "y");
        try {
            readOnly.add("z");
        } catch (UnsupportedOperationException e) {
            System.out.println("List.of 不可变 -> UnsupportedOperationException");
        }
    }
}
