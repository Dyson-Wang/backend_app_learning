# E05 配置、密钥与 DB migration

## 问题复现

配置散落在常量、`.env`、启动脚本里；密钥写进仓库；migration 手工跑，某台机器漏一个版本。

## 错误写法

```python
DB_URL = "postgres://root:pass@prod/db"     # 密钥进代码
conn.execute("ALTER TABLE ...")             # 每次启动都执行
```

## 正确写法

```python
if key in env:  return env[key], "env"
if key in file: return file[key], "file"
return default, "default"
```

密钥无默认值，缺失就 `sys.exit(1)`；migration 靠版本表跳过。

## 怎么验证

```bash
python3 main.py
```

预期：`port = 8080（来源：env）`、`db.host = localhost（来源：default）`；第二次 `0 个待应用`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 配置库 | `envconfig` | `dotenv` | `pydantic-settings` | `@ConfigurationProperties` | `figment` |
| migration | `golang-migrate` | `prisma migrate` | `alembic` | Flyway | `sqlx::migrate!` |
| 陷阱 | 忘 env 覆盖 | 全是字符串 | `.env` 未加载 | profile 混用 | build.rs 时机 |
