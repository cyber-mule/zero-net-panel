# TODO（首版闭环缺陷追踪）

## 严重
- 缺少用户注册与后台用户管理；当前仅有登录/刷新，无用户新增/禁用/重置等管理能力。 (Refs: api/auth/auth.api, internal/repository/user_repository.go, internal/bootstrap/seed/seed.go)
  - auth 模块需新增注册/找回/重置/验证 API，并配套验证码/邮件发送与频控策略。
  - admin 模块需新增用户管理 API（列表/搜索/创建/禁用/角色调整/重置密码/强制下线）。
  - 数据层补齐 users 认证/验证/锁定字段与审计日志表，扩展仓储 CRUD 与事务能力。
  - 配置与文档同步补充注册开关、邀请制/默认角色、账号生命周期说明与前端对接提示。
- 订单支付完成后未触发订阅创建/续期，订阅下发链路缺失。 (Refs: internal/logic/user/order/createlogic.go, internal/logic/admin/orders/paymentcallbacklogic.go, internal/logic/admin/orders/markpaidlogic.go, internal/repository/subscription_repository.go)
- 缺少后台订阅管理接口（创建/调整/禁用/延长等），仅有套餐与模板管理。 (Refs: api/admin/plans.api, api/admin/templates.api, api/user/subscriptions.api)
- 缺少后台节点管理接口（新增/编辑/禁用、内核端点配置），仅有列表/同步。 (Refs: api/admin/nodes.api, internal/handler/routes.go)

## 主要
- 下单时未校验支付通道是否存在/启用，外部支付可使用无效通道。 (Refs: internal/logic/user/order/createlogic.go, internal/repository/payment_channel_repository.go)
- 缺少用户侧支付通道列表接口，前端无法从配置渲染可用通道。 (Refs: api/admin/payment_channels.api)
- 外部支付网关仍为占位实现，无创建意图/签名校验/退款/对账等完整流程。 (Refs: docs/missing-capabilities.md, internal/logic/user/order/createlogic.go)
