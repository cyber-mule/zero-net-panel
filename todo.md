# TODO（首版闭环缺陷追踪）

## 阶段 0（阻断上线/高风险）
- 安全与访问控制：管理端已有 IP 白名单/限流，但缺 CORS 开关与用户侧/全局请求级限流/防刷能力。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md, internal/middleware/accessmiddleware.go)
- 运维工具补齐：已有探活/备份/systemd/Docker 基础，仍缺日志轮转示例、定时巡检/告警脚本、集成化备份/恢复流程。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md, docs/operations.md, scripts/healthcheck.sh, scripts/backup-db.sh)
- 账户安全：多因子认证等安全强化。 (Refs: docs/missing-capabilities.md)

## 阶段 1（核心能力补齐）
- gRPC 服务：仅有内核 gRPC provider 与配置，HTTP/gRPC 一体化服务与健康检查未落地。 (Refs: docs/ROADMAP.md, pkg/kernel/grpc_provider.go, internal/config/config.go)
- API 规格与对接文档：已有 `docs/api-reference.md`/`scripts/gen-api-docs.sh`，但仍缺 Swagger/OpenAPI、统一错误码/状态枚举与请求示例汇总。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md, docs/api-reference.md, scripts/gen-api-docs.sh)
- 通知体系：支付/退款等业务通知（邮件/短信/站内信）缺失。 (Refs: docs/missing-capabilities.md)
- 支付与结算扩展：通用外部网关已覆盖发起/回调/退款/订单级对账，但仍缺更多网关适配、开票/发票管理。 (Refs: docs/missing-capabilities.md, internal/logic/paymentutil/gateway.go)
- 计费后续能力：分账/多币种支持。 (Refs: docs/ROADMAP.md)

## 阶段 2（体验/运维增强）
- 安装向导：非交互模式仍标记为“未来支持”。 (Refs: docs/installation-wizard.md)
- 运维命令补齐：`znp cache purge` 仍为后续计划。 (Refs: docs/api-overview.md)
- 闭环验收清单/脚本：补充“创建绑定→发布入口→绑定套餐→下单→支付成功→订阅预览”的最小验证流程。 (Refs: docs/operations.md, docs/getting-started.md)
