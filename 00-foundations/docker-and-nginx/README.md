# Docker 与 Nginx

## 镜像与运行的最小要求

- 多阶段构建：编译产物和运行镜像分开，Go / Rust 可以直接跑在 `scratch` 或 `distroless` 上
- 非 root 用户运行
- 健康检查 `HEALTHCHECK` 与 `/healthz` 接口对齐
- 配置全部走环境变量，密钥不进镜像、不进代码库
- 日志输出到 stdout / stderr，由采集侧统一收集，别写文件

## Nginx 反代必须设的四类参数

1. **超时**：`proxy_connect_timeout` / `proxy_read_timeout` / `proxy_send_timeout`，不设的话下游卡住会一直占连接
2. **请求体**：`client_max_body_size`，配合上传场景 `advanced/D04`
3. **头透传**：`X-Request-ID`（trace id）、`X-Forwarded-For`、`X-Forwarded-Proto`
4. **上游健康检查与失败剔除**：`max_fails` / `fail_timeout`

`nginx.conf.template` 给了可直接改的最小片段。对应 `advanced/A04`（优雅关闭时要先从上游摘除，再停进程）。

## 优雅关闭的顺序

```
摘除上游/停止接收新流量 -> 等待存量请求处理完（带上限） -> 关闭连接池/消费者 -> 退出
```

`SIGTERM` 到来直接 `os.Exit` 是生产事故的经典来源。5 门语言各自的信号处理方式会写进 `advanced/A04`。

## 与 homelab 对接

`~/workwork2/docker-services` 已有统一反代和外部网络。示例服务只需暴露一个端口并挂到同一外部网络，Nginx 配置片段里用占位符 `<UPSTREAM_HOST>`，实际值按项目 `.env` 填，不写进本仓库。
