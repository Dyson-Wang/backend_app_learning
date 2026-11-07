# E06 单测、集成测试与容器化

## 问题复现

测试里连真实数据库，CI 一半时间在等容器起来；或者根本没测试，改一行代码靠手工点。

## 错误写法

```go
func TestCreateUser(t *testing.T) {
    db := connectRealDB()          // 单测连真库
    db.Exec("DELETE FROM users")   // 还会污染数据
}
```

## 正确写法

单测只测纯逻辑，依赖用接口替换：

```go
type Store interface { Create(u User) error }
// 测试里传 fakeStore，生产传 pgStore
```

集成测试单独一个 build tag / 目录，CI 里分开跑。

容器化三条硬要求：

```dockerfile
FROM golang:1.22 AS build
COPY . /src && WORKDIR /src
RUN CGO_ENABLED=0 go build -o /app ./cmd/server

FROM gcr.io/distroless/static
COPY --from=build /app /app
USER 65532
HEALTHCHECK CMD ["/app", "-healthcheck"]
ENTRYPOINT ["/app"]
```

## 怎么验证

```bash
go run .
```

预期：前两条 `PASS`，第三条 `FAIL`（故意），汇总 `2 通过，1 失败`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 测试框架 | `testing` 内建 | `node:test` 内建 | `unittest` 内建 | JUnit | `#[test]` 内建 |
| mock | 接口 + 手写 fake | `mock` | `unittest.mock` | Mockito | trait + fake |
| 容器 | 静态二进制最小 | `node:slim` | `python:slim` | JRE + 分层 jar | 静态二进制最小 |
