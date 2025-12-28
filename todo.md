# TODO（首版闭环缺陷追踪）

## 严重
无

## 主要
- gRPC 服务：HTTP/gRPC 一体化协议实现与健康检查尚未完成。 (Refs: docs/ROADMAP.md)
- API 规格与对接文档：缺少 Swagger/OpenAPI、统一错误码/状态枚举与请求示例汇总。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md)
- 运维工具补齐：日志轮转示例、定时巡检/告警脚本、集成化备份/恢复流程。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md)
- 安全与访问控制：CORS 开关、请求级限流/防刷能力。 (Refs: docs/ROADMAP.md, docs/missing-capabilities.md)
- 账户安全：多因子认证等安全强化。 (Refs: docs/missing-capabilities.md)
- 支付与结算扩展：更多网关适配、支付结果通知渠道、对账/开票/发票管理。 (Refs: docs/missing-capabilities.md)
- 通知体系：支付/退款等业务通知（邮件/短信/站内信）缺失。 (Refs: docs/missing-capabilities.md)
- 计费后续能力：分账/多币种支持。 (Refs: docs/ROADMAP.md)
- 安装向导：非交互模式标记为“未来支持”。 (Refs: docs/installation-wizard.md)
- 运维命令补齐：`znp cache purge` 标记为后续计划。 (Refs: docs/api-overview.md)

## 已完成
- 用户与权限增强：资料维护、自助改密/邮箱、密码策略、审计日志检索/导出。 (Refs: docs/missing-capabilities.md)
- 多协议内核：协议配置/节点协议绑定管理、手动下发、状态回调/轮询、流量倍数核算与查询、订阅渲染与用户节点状态调整。 (Refs: protocol.md, api/admin/protocol_configs.api, api/admin/protocol_bindings.api, internal/logic/admin/protocolbindings, internal/logic/kernel, docs/kernel-integration.md)
- 订单支付完成后未触发订阅创建/续期，订阅下发链路缺失。参考 xboard：支付成功（回调/人工标记）应创建或续期订阅、同步到期/流量、触发订阅下发（更新订阅内容/节点同步）并记录流水。 (Refs: internal/logic/user/order/createlogic.go, internal/logic/admin/orders/paymentcallbacklogic.go, internal/logic/admin/orders/markpaidlogic.go, internal/repository/subscription_repository.go)
- 缺少后台节点管理接口（新增/编辑/禁用、内核端点配置），仅有列表/同步。推进：补齐节点 CRUD + 状态切换，支持维护 HTTP/GRPC 内核端点（protocol/endpoint/token/revision/status/config），并与 `core.yaml` 的内核可访问 API 文档保持一致。 (Refs: api/admin/nodes.api, internal/handler/routes.go, internal/repository/node_repository.go, pkg/kernel/*, core.yaml)
- 缺少用户侧支付通道列表接口，前端无法从配置渲染可用通道。 (Refs: api/user/payment_channels.api, internal/handler/routes.go, internal/logic/user/paymentchannels/listlogic.go)
- 下单时未校验支付通道是否存在/启用，外部支付可使用无效通道。 (Refs: internal/logic/user/order/createlogic.go, internal/repository/payment_channel_repository.go)
- 外部支付网关发起支付仍为占位实现，补齐通用发起能力（仅发起支付 + 回调）并完善文档说明。 (Refs: internal/logic/paymentutil/gateway.go, internal/logic/user/order/createlogic.go, docs/api-reference.md)
- 外部支付退款/对账/签名校验未覆盖（当前仅支持发起支付与回调）。 (Refs: internal/logic/paymentutil/gateway.go, internal/logic/admin/orders/refundlogic.go, internal/logic/admin/orders/reconcilelogic.go)
- 补充外部支付联调示例与 mock 网关脚本。 (Refs: docs/payment-gateway-demo.md, scripts/mock-payment-gateway.go)
- 计费后续能力：优惠券/折扣。 (Refs: api/admin/coupons.api, internal/logic/admin/coupons, internal/repository/coupon_repository.go)
