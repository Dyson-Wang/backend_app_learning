# E06 单测、集成测试与容器化

## 问题复现

测试里连真库，CI 一半时间在等容器；或者根本没测试，改一行靠手工点。

## 错误写法

```js
test('create user', async () => {
  await realDb.query('DELETE FROM users');   // 单测连真库还污染数据
});
```

## 正确写法

依赖注入 + fake：

```js
const makeService = (store) => ({ create: (u) => store.create(u) });
// 测试传 fakeStore，生产传 pgStore
```

Node 22 自带 `node:test`，不用装 jest：

```js
import { test } from 'node:test';
import assert from 'node:assert/strict';
test('add', () => assert.equal(add(1, 2), 3));
```

容器化：多阶段 + `node:slim` + `USER node` + `HEALTHCHECK`。

## 怎么验证

```bash
node index.mjs
```

预期：前两条 `PASS`，第三条 `FAIL`（故意），汇总 `2 通过，1 失败`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 测试框架 | `testing` | `node:test` | `unittest` | JUnit | `#[test]` |
| mock | 接口 + fake | 手写 / `mock` | `unittest.mock` | Mockito | trait + fake |
| 容器基础镜像 | distroless | `node:slim` | `python:slim` | JRE | distroless |
