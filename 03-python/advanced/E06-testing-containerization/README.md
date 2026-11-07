# E06 单测、集成测试与容器化

## 问题复现

测试连真库，CI 一半时间等容器；或者没测试，改一行靠手工点。

## 错误写法

```python
def test_create_user():
    db.execute("DELETE FROM users")     # 单测连真库还污染数据
```

## 正确写法

依赖注入 + fake：

```python
class Store(Protocol):
    def create(self, u: User) -> None: ...

def make_service(store: Store): ...
# 测试传 FakeStore，生产传 PgStore
```

pytest 用 `-m "not integration"` 把两类分开跑。

容器化：`python:3.13-slim` + 非 root 用户 + `HEALTHCHECK`。

## 怎么验证

```bash
python3 main.py
```

预期：前两条 `PASS`，第三条 `FAIL`（故意），汇总 `2 通过，1 失败`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 测试框架 | `testing` | `node:test` | `unittest` / pytest | JUnit | `#[test]` |
| mock | 接口 + fake | 手写 | `unittest.mock` | Mockito | trait + fake |
| 容器镜像 | distroless | `node:slim` | `python:slim` | JRE | distroless |
