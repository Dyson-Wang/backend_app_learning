# E05 配置、密钥与 DB migration

## 问题复现

配置散落在常量、`.env`、启动脚本里；密钥写进仓库；migration 手工跑，某台机器漏一个版本。

## 错误写法

```rust
const DB_URL: &str = "postgres://root:pass@prod/db";   // 密钥进代码
sqlx::query("ALTER TABLE ...").execute(&pool).await?;  // 每次启动都执行
```

## 正确写法

```rust
if let Some(v) = env.get(key) { return (v, "env"); }
if let Some(v) = file.get(key) { return (v, "file"); }
(def, "default")
```

Rust 的杀手锏是 `sqlx::migrate!("./migrations")`：编译期就把 migration 打进二进制，运行时按 `_sqlx_migrations` 表跳过。

## 怎么验证

```bash
cargo run --quiet
```

预期：`port = 8080（来源：env）`、`db.host = localhost（来源：default）`；第二次 `0 个待应用`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 配置库 | `envconfig` | `dotenv` | `pydantic-settings` | `@ConfigurationProperties` | `figment` |
| migration | 外部工具 | `prisma migrate` | `alembic` | Flyway | **编译期内嵌** |
| 陷阱 | 忘 env 覆盖 | 全是字符串 | `.env` 未加载 | profile 混用 | 改 migration 需重新编译 |
