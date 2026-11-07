# E06 单测、集成测试与容器化

## 问题复现

测试连真库，CI 一半时间等容器；或者没测试，改一行靠手工点。

## 错误写法

```java
@Test
void createUser() {
    jdbc.execute("DELETE FROM users");    // 单测连真库还污染数据
}
```

## 正确写法

```java
interface Store { void create(User u); }
// 测试用 FakeStore，生产用 PgStore；Spring 里 @MockBean 自动注入
```

集成测试用 `@SpringBootTest` + Testcontainers，和单测分开跑。

容器化：`eclipse-temurin:21-jre` + 非 root + `HEALTHCHECK`。

## 怎么验证

```bash
javac Main.java && java Main
```

预期：前两条 `PASS`，第三条 `FAIL`（故意），汇总 `2 通过，1 失败`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 测试框架 | `testing` | `node:test` | `unittest` | JUnit | `#[test]` |
| mock | 接口 + fake | 手写 | `unittest.mock` | **Mockito** | trait + fake |
| 容器镜像 | distroless | `node:slim` | `python:slim` | JRE（分层 jar） | distroless |
