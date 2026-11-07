# D03 顺序性保证

## 问题复现

`create → pay → ship` 并发消费，`ship` 先到，状态机报错。

## 错误写法

```java
for (Msg m : queue) {
    executor.submit(() -> handle(m));   // 线程池乱序
}
```

## 正确写法

1. 分区内串行：`hash(key) % N`，每分区一个消费者。
2. 消费端按 key 串行：`ConcurrentHashMap<key, CompletableFuture链>`。

```java
int h = 0;
for (char ch : key.toCharArray()) h = (h * 31 + ch) % n;
```

## 怎么验证

```bash
javac Main.java && java Main
```

预期：乱序报 `ship 在 INIT 状态下不合法`；顺序 `最终状态 = SHIPPED`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 分区 | 手写 hash | 同 | 同 | 同 | 同 |
| 串行手段 | 单 goroutine | 单线程 | 单线程 | 单线程 executor | 单任务 |
| 陷阱 | channel 多消费者 | await 顺序 | GIL 不保证 | **线程池乱序** | 多线程乱序 |
