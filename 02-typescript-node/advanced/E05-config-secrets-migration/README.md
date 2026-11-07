# E05 配置、密钥与 DB migration

## 问题复现

配置散落在常量、`.env`、启动脚本里，线上改了端口不知道生效没；密钥写进仓库；migration 手工跑，某台机器漏一个版本。

## 错误写法

```js
const DB_URL = 'postgres://root:pass@prod/db';   // 密钥进代码
await db.query('ALTER TABLE ...');               // 每次启动都执行
```

## 正确写法

```js
if (key in env) return { value: env[key], origin: 'env' };
if (key in file) return { value: file[key], origin: 'file' };
return { value: def, origin: 'default' };
```

密钥无默认值，缺失就退出；migration 用版本表跳过。

## 怎么验证

```bash
node index.mjs
```

预期：`port = 8080（来源：env）`、`db.host = localhost（来源：default）`；第二次 `0 个待应用`。

## 和其他语言的差异

| 机制 | Go | Node.js | Python | Java | Rust |
|---|---|---|---|---|---|
| 配置库 | `envconfig` | `dotenv` | `pydantic-settings` | `@ConfigurationProperties` | `figment` |
| migration | `golang-migrate` | `prisma migrate` | `alembic` | Flyway | `sqlx::migrate!` |
| 陷阱 | 忘了 env 覆盖 | `process.env` 全是字符串 | `.env` 未加载 | profile 混用 | build.rs 时机 |
