# E05 配置、密钥与 DB migration

## 问题复现

配置散落在常量、`.env`、启动脚本里；密钥写进仓库；migration 手工跑，某台机器漏一个版本。

## 错误写法

```java
private static final String DB_URL = "postgres://root:pass@prod/db";
jdbc.execute("ALTER TABLE ...");
```

## 正确写法

```java
if (env.containsKey(key))  return new Source(env.get(key), "env");
if (file.containsKey(key)) return new Source(file.get(key), "file");
return new Source(def, "default");
```

Spring 里就是 `application.yml` < `application-prod.yml` < 环境变量；migration 用 Flyway 的 `flyway_schema_history` 表。

## 怎么验证

```bash
javac Main.java && java Main
```

预期：`port = 8080（来源：env）`、`db.host = localhost（来源：default）`；第二次 `0 个待应用`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 配置库 | `envconfig` | `dotenv` | `pydantic-settings` | `@ConfigurationProperties` | `figment` |
| migration | `golang-migrate` | `prisma migrate` | `alembic` | **Flyway / Liquibase** | `sqlx::migrate!` |
| 陷阱 | 忘 env 覆盖 | 全是字符串 | `.env` 未加载 | profile 混用 | build.rs 时机 |
