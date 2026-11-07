# Go · advanced 生产场景

这里不是语法练习，是生产上真会咬人的问题。每个场景统一结构：

1. **问题复现**：先让故障真的发生（有并发 / 自测脚本）
2. **错误写法**：展示直觉上会写错的版本
3. **正确写法**：可运行的实现
4. **怎么验证**：给出命令与预期输出

场景编号与目录名全局统一，完整对照表见 `../../99-cross-language/scenario-matrix.md`。
公共底座概念在 `../../00-foundations/`，先读那一层再来做场景，否则只会抄代码。

进度：**28/28**（由 `scripts/refresh_indexes.py` 自动扫描生成）

### A · 流量与稳定性

| 编号 | 目录 | 场景 | 状态 |
|---|---|---|---|
| A01 | `A01-rate-limiting` | 限流：令牌桶 / 滑动窗口，单机与分布式 | 已实现 |
| A02 | `A02-timeout-retry-backoff` | 超时、重试、指数退避与抖动 | 已实现 |
| A03 | `A03-circuit-breaker` | 熔断与降级 | 已实现 |
| A04 | `A04-graceful-shutdown` | 优雅关闭、健康检查、上游摘除 | 已实现 |
| A05 | `A05-backpressure` | 背压与并发上限 | 已实现 |

### B · 一致性与并发

| 编号 | 目录 | 场景 | 状态 |
|---|---|---|---|
| B01 | `B01-transactions-and-isolation` | 事务与隔离级别、乐观锁 | 已实现 |
| B02 | `B02-idempotency` | 幂等键与重复请求 | 已实现 |
| B03 | `B03-distributed-lock` | 分布式锁 | 已实现 |
| B04 | `B04-race-and-deadlock` | 竞态、死锁与调试 | 已实现 |
| B05 | `B05-context-propagation` | 上下文取消与跨服务传播 | 已实现 |

### C · 数据与缓存

| 编号 | 目录 | 场景 | 状态 |
|---|---|---|---|
| C01 | `C01-connection-pool` | 连接池与资源耗尽 | 已实现 |
| C02 | `C02-n-plus-one` | N+1 查询与批量写 | 已实现 |
| C03 | `C03-cursor-pagination` | 游标分页替代 OFFSET | 已实现 |
| C04 | `C04-cache-stampede` | 缓存穿透 / 击穿 / 雪崩 | 已实现 |
| C05 | `C05-cache-db-consistency` | 缓存与 DB 双写一致性 | 已实现 |

### D · 消息与异步

| 编号 | 目录 | 场景 | 状态 |
|---|---|---|---|
| D01 | `D01-at-least-once-delivery` | 至少一次投递与重复消费 | 已实现 |
| D02 | `D02-dead-letter-queue` | 死信队列与重试策略 | 已实现 |
| D03 | `D03-ordering` | 顺序性保证 | 已实现 |
| D04 | `D04-streaming` | SSE / WebSocket / gRPC 流式 | 已实现 |

### E · 安全与工程化

| 编号 | 目录 | 场景 | 状态 |
|---|---|---|---|
| E01 | `E01-jwt-refresh-revoke` | JWT 刷新与吊销 | 已实现 |
| E02 | `E02-password-rbac` | 密码哈希与 RBAC | 已实现 |
| E03 | `E03-validation-error-code` | 输入校验与统一错误码 | 已实现 |
| E04 | `E04-observability` | 结构化日志 / trace id / 指标 | 已实现 |
| E05 | `E05-config-secrets-migration` | 配置、密钥与 DB migration | 已实现 |
| E06 | `E06-testing-containerization` | 单测、集成测试与容器化 | 已实现 |

### F · 故障复盘

| 编号 | 目录 | 场景 | 状态 |
|---|---|---|---|
| F01 | `F01-slow-query` | 慢查询拖垮连接池 | 已实现 |
| F02 | `F02-retry-storm` | 重试风暴 | 已实现 |
| F03 | `F03-message-backlog` | 消息积压 | 已实现 |


## 完成一个场景的判据

能不看代码回答三个问题：故障是怎么发生的？这门语言用哪个机制解决？换成另一门语言时机制换成什么？
