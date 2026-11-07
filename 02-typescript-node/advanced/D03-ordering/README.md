# D03 顺序性保证

## 问题复现

`create → pay → ship` 三个事件被并发消费，`ship` 先到，状态机里 `PAID` 还没出现，直接报错。

## 错误写法

```js
for await (const msg of queue) {
  handle(msg);        // 谁先拿到谁处理，顺序丢了
}
```

## 正确写法

1. 分区内串行：`hash(key) % N` 定分区，每分区一个消费者。
2. 消费端按 key 串行：`Map<key, Promise链>`，同 key 排队。

```js
let h = 0;
for (const ch of key) h = (h * 31 + ch.codePointAt(0)) % n;
```

## 怎么验证

```bash
node index.mjs
```

预期：乱序报 `ship 在 INIT 状态下不合法`；顺序 `最终状态 = SHIPPED`；`order-42` 两次同分区。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 分区 | 手写 hash | 同 | 同 | 同 | 同 |
| 串行手段 | 单 goroutine | **单线程天然串行** | 单线程 | 单线程 executor | 单任务 |
| 陷阱 | channel 多消费者 | `await` 顺序≠到达顺序 | GIL 不保证顺序 | 线程池乱序 | 多线程乱序 |
