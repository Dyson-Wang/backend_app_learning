# E05 配置、密钥与 DB migration

## 问题复现

配置散落在代码常量、`.env`、启动脚本里，线上改了端口不知道生效没；密钥写进仓库；migration 手工执行，某台机器漏跑一个版本。

## 错误写法

```go
const DBURL = "postgres://root:pass@prod/db"   // 密钥进代码
db.Exec("ALTER TABLE ...")                      // 每次启动都执行
```

## 正确写法

三层来源，且能说清「值从哪来」：

```go
if v, ok := env[key]; ok { return Source{v, "env"} }
if v, ok := file[key]; ok { return Source{v, "file"} }
return Source{def, "default"}
```

密钥**没有默认值**，缺失就拒绝启动：

```go
if _, ok := secrets["DATABASE_URL"]; !ok { log.Fatal("DATABASE_URL 未设置") }
```

migration 靠版本表跳过已应用的：

```sql
SELECT version FROM schema_migrations;
```

## 怎么验证

```bash
go run .
```

预期：`port = 8080（来源：env）`、`db.host = localhost（来源：default）`、`db.name = brm（来源：file）`；第二次执行 `0 个待应用`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 配置库 | `envconfig` | `dotenv` | `pydantic-settings` | Spring `@ConfigurationProperties` | `figment` |
| migration | `golang-migrate` | `knex` / `prisma` | `alembic` | Flyway / Liquibase | `sqlx::migrate!` |
| 密钥 | 环境变量 / Vault | 同 | 同 | 同 | 同 |
