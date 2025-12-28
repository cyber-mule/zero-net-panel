# TODO（首版闭环缺陷追踪）

## 严重
无

## 主要
- gRPC 服务：仅有内核 gRPC provider 与配置，HTTP/gRPC 一体化服务与健康检查未落地。 (Refs: docs/ROADMAP.md, pkg/kernel/grpc_provider.go, internal/config/config.go)
- API 规格与对接文档：已有 `docs/api-reference.md`/`scripts/gen-api-docs.sh`，但仍缺 Swagger/OpenAPI、统一错误码/状态枚举与请求示例汇总。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md, docs/api-reference.md, scripts/gen-api-docs.sh)
- 运维工具补齐：已有探活/备份/systemd/Docker 基础，仍缺日志轮转示例、定时巡检/告警脚本、集成化备份/恢复流程。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md, docs/operations.md, scripts/healthcheck.sh, scripts/backup-db.sh)
- 安全与访问控制：管理端已有 IP 白名单/限流，但缺 CORS 开关与用户侧/全局请求级限流/防刷能力。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md, internal/middleware/accessmiddleware.go)
- 账户安全：多因子认证等安全强化。 (Refs: docs/missing-capabilities.md)
- 支付与结算扩展：通用外部网关已覆盖发起/回调/退款/订单级对账，但仍缺更多网关适配、开票/发票管理。 (Refs: docs/missing-capabilities.md, internal/logic/paymentutil/gateway.go)
- 通知体系：支付/退款等业务通知（邮件/短信/站内信）缺失。 (Refs: docs/missing-capabilities.md)
- 计费后续能力：分账/多币种支持。 (Refs: docs/ROADMAP.md)
- 安装向导：非交互模式仍标记为“未来支持”。 (Refs: docs/installation-wizard.md)
- 运维命令补齐：`znp cache purge` 仍为后续计划。 (Refs: docs/api-overview.md)

## 已完成
- 用户与权限增强：资料维护、自助改密/邮箱、密码策略、审计日志检索/导出已落地。 (Refs: docs/missing-capabilities.md)
- 多协议内核：协议配置/节点协议绑定管理、手动下发、状态回调/轮询、流量倍数核算与查询、订阅渲染与用户节点状态调整已落地。 (Refs: protocol.md, api/admin/protocol_configs.api, api/admin/protocol_bindings.api, internal/logic/admin/protocolbindings, internal/logic/kernel, docs/kernel-integration.md)
- 支付成功后订阅创建/续期与订阅下发链路已补齐（回调/人工标记触发同步到期/流量、更新订阅内容/节点同步并记录流水）。 (Refs: internal/logic/user/order/createlogic.go, internal/logic/admin/orders/paymentcallbacklogic.go, internal/logic/admin/orders/markpaidlogic.go, internal/repository/subscription_repository.go)
- 节点管理接口已补齐：节点 CRUD + 状态切换，支持维护 HTTP/GRPC 内核端点配置，并与 `core.yaml` 内核 API 文档保持一致。 (Refs: api/admin/nodes.api, internal/handler/routes.go, internal/repository/node_repository.go, pkg/kernel/*, core.yaml)
- 用户侧支付通道列表接口已补齐。 (Refs: api/user/payment_channels.api, internal/handler/routes.go, internal/logic/user/paymentchannels/listlogic.go)
- 下单时支付通道存在/启用校验已补齐。 (Refs: internal/logic/user/order/createlogic.go, internal/repository/payment_channel_repository.go)
- 外部支付网关发起支付能力已落地，并补充文档说明。 (Refs: internal/logic/paymentutil/gateway.go, internal/logic/user/order/createlogic.go, docs/api-reference.md)
- 外部支付退款/对账/签名校验已覆盖。 (Refs: internal/logic/paymentutil/gateway.go, internal/logic/admin/orders/refundlogic.go, internal/logic/admin/orders/reconcilelogic.go)
- 补充外部支付联调示例与 mock 网关脚本。 (Refs: docs/payment-gateway-demo.md, scripts/mock-payment-gateway.go)
- 计费后续能力：优惠券/折扣已落地。 (Refs: api/admin/coupons.api, internal/logic/admin/coupons, internal/repository/coupon_repository.go)
