# 可观测性

出问题时查不到是哪一步慢，比出 bug 更常见。

## 三根支柱

| 支柱 | 落地形态 | 最小要求 |
|---|---|---|
| 日志 | 结构化 JSON 到 stdout | 每行带 `ts` / `level` / `trace_id` / `msg` / 关键业务字段 |
| 指标 | RED（Rate / Errors / Duration）+ 资源 | 按 handler + status 统计 QPS、错误率、P95/P99 延迟 |
| 追踪 | trace id + span | 跨网关、应用、DB、MQ 传递同一个 id |

## trace id 传递链

```
Nginx 生成/透传 X-Request-ID
  -> 应用中间件读入，放进请求上下文（Go context / Node AsyncLocalStorage / Python contextvars / Java ThreadLocal+MDC / Rust tower Extension）
  -> 打日志时自动带上
  -> 调用下游时写进 HTTP header / MQ message headers
  -> 慢查询日志、外部调用日志也带同一个 id
  -> 响应头回传，用户报障时给一个 id 就能定位
```

关键点：**不要把 trace id 当参数一层层手动传**。每门语言都有对应的请求上下文机制，这本身就是一个很好的跨语言对照主题。

## 健康检查区分两个接口

- `/healthz`（liveness）：进程还活着吗？只做极轻的自检，不查依赖
- `/readyz`（readiness）：能接流量了吗？检查 DB / 缓存 / 依赖是否可用

把依赖检查放进 liveness 会导致数据库一抖动，编排系统就把所有实例重启一遍——这是很常见的自伤配置。

## 日志不要做的事

- 不打密码、token、身份证、完整请求体
- 不打无上下文的 `fmt.Println("error")`
- 不在热路径里同步写文件 / 远程日志（会变成新的瓶颈）
- 不用字符串拼接当结构化字段

对应 `advanced/E04`。
