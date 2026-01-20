# API 模块与业务逻辑说明

本文档汇总 Zero Net Panel 已实现的 REST API 模块，并补充关键业务的端到端流程、错误码与排障建议，方便前后端协作与第三方集成。

## 管理端模块

| 模块 | 路径 | 说明 |
| ---- | ---- | ---- |
| 仪表盘 | `/api/v1/{admin}/dashboard` | 展示模块导航、权限控制 |
| 用户管理 | `/api/v1/{admin}/users` | 用户列表、创建、禁用、角色调整、重置密码、强制下线 |
| 节点管理 | `/api/v1/{admin}/nodes` | 节点查询、创建、更新、禁用、删除（软删除）、协议内核同步、状态同步 |
| 订阅模板 | `/api/v1/{admin}/subscription-templates` | 模板 CRUD、发布、历史追溯 |
| 订阅管理 | `/api/v1/{admin}/subscriptions` | 订阅列表、创建、调整、禁用、延长有效期 |
| 套餐管理 | `/api/v1/{admin}/plans` | 套餐列表、创建、更新，字段涵盖价格、时长、流量限制等 |
| 套餐计费选项 | `/api/v1/{admin}/plans/{plan_id}/billing-options` | 为套餐维护多周期/多价格选项（小时/天/月/年） |
| 优惠券管理 | `/api/v1/{admin}/coupons` | 优惠券创建、启停、限额与有效期维护 |
| 公告中心 | `/api/v1/{admin}/announcements` | 公告列表、创建、发布，支持置顶与可见时间窗 |
| 安全配置 | `/api/v1/{admin}/security-settings` | 读取与更新第三方签名/加密开关、凭据与时间窗口 |
| 审计日志 | `/api/v1/{admin}/audit-logs` | 审计日志检索与导出 |
| 订单管理 | `/api/v1/{admin}/orders` | 检索、查看订单，支持多支付方式、外部流水追踪、手动标记支付/取消与余额退款 |

> `{admin}` 为可配置的后台前缀，默认为 `admin`，可通过 `Admin.RoutePrefix` 自定义。

## 用户端模块

> 注册/找回/验证接口已开放，需在配置中开启注册开关并配置邮件发送与验证码策略。

- `/api/v1/user/subscriptions`：用户订阅列表、预览、模板切换。
- `/api/v1/subscriptions/{token}`：订阅拉取地址（免登录，按 `User-Agent` 选择模板）。
- `/api/v1/user/plans`：面向终端的套餐列表，返回价格、特性、流量限制与 `billing_options`。
- `/api/v1/user/nodes`：用户侧节点运行状态（脱敏展示）。
- `/api/v1/user/announcements`：按受众过滤当前有效公告，支持置顶排序与限量返回。
- `/api/v1/user/account/balance`：返回当前余额、币种以及流水历史。
- `/api/v1/user/account/profile`：用户资料查询与更新。
- `/api/v1/user/account/password`：用户自主改密。
- `/api/v1/user/account/email`：用户自主改邮箱（验证码流程）。
- `/api/v1/user/orders`：创建、查询订单并支持取消待支付或零元订单，返回计划快照、条目与余额快照（可选 `billing_option_id`）。
  - 用户侧默认不返回 `disabled` 状态订阅，`expired` 仍可展示用于续费。

### 订单操作补充说明

- 用户端 `POST /api/v1/user/orders` 新增 `payment_method`、`payment_channel`、`payment_return_url`、`coupon_code` 字段：
  - 默认 `payment_method = balance`，系统直接扣减余额、记录 `balance_transactions`，订单状态立即变为 `paid`、`payment_status = succeeded`。
  - 当 `payment_method = external` 且金额大于零时，会生成 `pending_payment` 订单，创建 `order_payments` 预订单记录，并按支付通道 `config` 发起支付；响应包含 `payment_intent_id` 与 `payments`，其中 `payments[].metadata.pay_url`/`qr_code` 可用于跳转或展示二维码，余额不会变动。
  - 当 `payment_method = manual` 时，会生成待支付订单；需管理员通过 `/api/v1/{admin}/orders/{id}/pay` 标记已支付。
- 若命中优惠券，`order.items` 会追加 `item_type=discount` 条目并在 `order.metadata` 回填折扣信息。
- 用户端 `POST /api/v1/user/orders/{id}/cancel` 仅允许取消待支付或零金额订单，不触发余额回滚。
- 用户端 `GET /api/v1/user/orders/{id}/payment-status` 用于前端轮询确认支付结果。
- 管理端提供 `POST /api/v1/{admin}/orders/{id}/pay`、`/cancel`、`/refund` 与 `/orders/payments/reconcile`，需管理员角色；外部支付退款会按通道 `config.refund` 发起并记录退款流水。
- 所有用户端接口默认需要 JWT 鉴权，同时可选启用第三方加密认证中间件，对请求进行签名验证与 AES-GCM 解密。
- 外部支付回调可按以下流程接入：
  1. 网关回调携带支付状态后，通过内部逻辑 `PaymentCallbackLogic`（或后续开放的专用接口）调用 `UpdatePaymentState`、`UpdatePaymentRecord`，将订单状态从 `pending_payment` 更新为 `paid`/`payment_failed`，并填充 `payment_reference`、`payment_failure_*` 字段。
  2. 回调完成后，`GET /api/v1/user/orders/:id` 与 `/admin/orders/:id` 均会返回最新的 `payment_status`、`payments` 明细，方便前端落地扫码/轮询场景。

## 端到端流程

### 第三方签名校验流程

1. 管理员通过 `GET /api/v1/{admin}/security-settings` 查询当前开关与密钥。
2. 根据需要调用 `PATCH /api/v1/{admin}/security-settings` 设置 `thirdPartyAPIEnabled`、`apiKey`、`apiSecret` 与 `nonceTTLSeconds`。
3. 第三方客户端在调用任何受保护接口时，按照 `timestamp + "\n" + nonce + "\n" + body` 规则生成 HMAC-SHA256 签名，并随请求携带以下头：
   - `X-ZNP-API-Key`
   - `X-ZNP-Timestamp`（Unix 秒）
   - `X-ZNP-Nonce`（建议 16 字节随机值）
   - `X-ZNP-Signature`
   - `X-ZNP-Encrypted: true` 与 `X-ZNP-IV`（可选，当启用 AES-256-GCM 加密时必填）
4. 服务端校验签名、时间窗口与随机数重复使用情况，必要时进行解密后再继续路由。

| 接口 | 错误码 | 说明 | 排障建议 |
| ---- | ------ | ---- | -------- |
| `PATCH /api/v1/{admin}/security-settings` | `400100` | 参数缺失或 TTL 小于 60 秒 | 校验请求体字段是否齐全，确认 `nonceTTLSeconds >= 60`。 |
| 同上 | `409002` | 存在并发更新冲突 | 使用最新版 `updatedAt` 再次提交，或开启重试机制。 |
| 受保护接口（任意） | `401001` | 签名不一致 | 确保使用 `apiSecret` 计算 HMAC，检查换行与大小写是否匹配。 |
| 受保护接口（任意） | `403001` | 时间戳超出窗口 | 对齐客户端时间，必要时缩短网络传输延迟或增大 `nonceTTLSeconds`。 |
| 受保护接口（任意） | `403002` | Nonce 重复使用 | 确认客户端在重试时生成全新随机数。 |

### 节点管理能力

- `POST /api/v1/{admin}/nodes`：创建节点基础信息与控制面信息（`control_endpoint`/AK-SK/`control_token`）。
- `PATCH /api/v1/{admin}/nodes/{id}`：更新节点元数据、控制面地址与标签（节点控制面必填，不再回退全局）。
- `status_sync_enabled`：控制是否允许节点状态自动同步（默认 true）。
- 启用后，服务会定时调用内核 `GET /v1/status`，将节点 `status` 更新为 `online`/`offline`。
- `POST /api/v1/{admin}/nodes/{id}/disable`：下线/禁用节点（替代物理删除）。
- `POST /api/v1/{admin}/nodes/{id}/kernels/sync`：触发节点内核配置同步（`protocol` 可选）。
- `POST /api/v1/{admin}/nodes/status/sync`：手动同步节点状态（仅处理指定 `node_ids`）。
- `POST /api/v1/{admin}/protocol-bindings/status/sync`：手动反向同步协议健康状态（仅处理指定 `node_ids`）。
- `GET /api/v1/user/nodes`：用户侧查看节点运行状态，隐藏内核端点与配置等敏感信息。

创建节点字段提示：

- 必填：`name`
- 必填：`control_endpoint`（节点控制面地址）
- 常用：`control_access_key`/`control_secret_key` 或 `control_token`

### 节点同步流程

1. 管理端列表接口 `GET /api/v1/{admin}/nodes` 返回节点详情与最新同步时间。
2. 运维人员选择目标节点，调用 `POST /api/v1/{admin}/nodes/{id}/kernels/sync` 触发与内核的即时同步（`protocol` 可选，空表示默认协议）。
3. 服务端拉取内核配置并更新记录，立即返回 `revision` 与 `synced_at` 等结果字段。
4. 若开启 Prometheus，观察 `znp_node_sync_operations_total` 与 `znp_node_sync_duration_seconds` 判断成功率与耗时。

### 节点状态同步流程

1. 定时任务按 `Kernel.StatusPollInterval` 轮询 `/v1/status`，成功则更新节点 `status=online`，失败更新为 `offline`。
2. 需要即时更新时，可调用 `POST /api/v1/{admin}/nodes/status/sync` 并传入 `node_ids`。
3. 响应中返回每个节点的 `status` 与 `message`，便于定位控制面鉴权或地址问题。

### 协议健康反向同步流程

1. 调用 `POST /api/v1/{admin}/protocol-bindings/status/sync` 并传入 `node_ids`。
2. 服务端按节点控制面分组拉取 `/v1/status`，将协议绑定 `health_status` 更新为 `healthy/degraded/unhealthy/offline/unknown`。
3. 响应中返回每个节点的同步结果与更新数量。

| 接口 | HTTP 状态码 | 说明 | 排障建议 |
| ---- | ----------- | ---- | -------- |
| `GET /api/v1/{admin}/nodes` | `400` | 过滤条件非法 | 确认查询参数（如 `protocol`、`status`）是否在允许范围内。 |
| `POST /api/v1/{admin}/nodes/{id}/kernels/sync` | `400` | 协议不支持 | 确认 `protocol` 参数与内核 Provider 配置一致。 |
| 同上 | `404` | 节点不存在 | 检查节点是否被删除，确认 `Admin.RoutePrefix` 与 URL 中的 `{id}` 是否正确。 |
| 同上 | `500` | 内核同步失败 | 检查节点 `control_endpoint` 与鉴权信息，必要时抓取内核 HTTP 日志。 |

### 协议绑定与发布流程

1. 创建协议绑定 `POST /api/v1/{admin}/protocol-bindings`：必填 `node_id`、`protocol`、`role`（`listener`/`connector`）、`kernel_id`（字符串，需与内核协议 ID 对齐）、`profile`（内核实际配置）；常用字段 `listen`、`connect`、`access_port`。
2. 创建协议发布 `POST /api/v1/{admin}/protocol-entries`：必填 `binding_id`、`entry_address`、`entry_port`；`profile` 填写对外公开配置（如 reality 公钥、short_id 等）。
3. 更新绑定或发布（可选）：`PATCH /api/v1/{admin}/protocol-bindings/{id}` / `PATCH /api/v1/{admin}/protocol-entries/{id}`。

补充说明：
- `entry_address/entry_port` 为对外入口地址，可与绑定的 `listen/access_port` 不一致，用于中转或分流场景。
- 协议发布 `status` 仅影响用户可见性；健康状态以协议绑定 `health_status` 为准。
- 绑定 `listen` 为空或仅端口时，会用 `access_port` 归一化为 `0.0.0.0:<port>` 供内核使用。

### 协议绑定同步流程

1. 管理端创建协议绑定与发布（`/protocol-bindings`、`/protocol-entries`）。
2. 触发单条或批量下发：`POST /api/v1/{admin}/protocol-bindings/{id}/sync` 或 `/protocol-bindings/sync`。
3. 同步结果直接返回，包含 `binding_id`、`status`、`message`、`synced_at`。
4. 若未配置内核控制面，返回 `status=error` 且 `message` 提示配置缺失。

批量下发示例：
```json
{"binding_ids":[3,4]}
```

### 套餐发布流程

1. 管理端 `POST /api/v1/{admin}/plans` 创建套餐：必填 `name`、`price_cents`、`currency`、`duration_days`；可选 `binding_ids`、`traffic_limit_bytes`、`devices_limit`。
2. （可选）为套餐添加计费选项 `POST /api/v1/{admin}/plans/{plan_id}/billing-options`。
3. 前端或第三方调用 `GET /api/v1/user/plans` 验证套餐是否对终端可见（需 `status=active` 且 `visible=true`）。
4. 订单创建时，`POST /api/v1/user/orders` 会读取套餐快照、扣减余额并返回结果。

### 订阅创建与交付流程

1. 管理端准备订阅模板：`POST /api/v1/{admin}/subscription-templates` 创建，`POST /api/v1/{admin}/subscription-templates/{id}/publish` 发布。
2. 创建订阅 `POST /api/v1/{admin}/subscriptions`：必填 `user_id`、`name`、`plan_id`、`template_id`、`expires_at`、`traffic_total_bytes`、`devices_limit`；可选 `available_template_ids`、`token`。
3. 用户侧拉取与预览：`GET /api/v1/user/subscriptions`、`GET /api/v1/user/subscriptions/{id}/preview`；切换模板 `POST /api/v1/user/subscriptions/{id}/template`。
4. 公开订阅（免登录）：`GET /api/v1/subscriptions/{token}`。

| 接口 | 错误码 | 说明 | 排障建议 |
| ---- | ------ | ---- | -------- |
| `POST /api/v1/{admin}/subscription-templates/{id}/publish` | `404010` | 模板不存在或无权限 | 校验模板 ID 与管理员角色；检查是否已归档。 |
| 同上 | `409001` | 模板存在未发布草稿 | 先保存最新草稿，再重新发起发布或删除旧草稿。 |
| `POST /api/v1/{admin}/plans` | `400201` | 套餐字段缺失或价格非法 | 核对必填字段（`name`、`price_cents`、`currency`、`duration_days`），确保价格 > 0。 |
| 同上 | `409201` | 套餐名称已存在 | 更换名称或在更新接口中使用已有套餐 ID。 |
| `GET /api/v1/user/plans` | `503001` | 套餐缓存构建失败 | 查看缓存服务状态，必要时执行 `znp cache purge`（后续计划）或重启服务。 |
| `POST /api/v1/user/orders` | `402001` | 余额不足 | 提示用户充值或调整套餐价格。 |
| 同上 | `409301` | 套餐不可用 | 确认套餐状态为 `active` 且 `visible=true`，或检查权限配置。 |

## 第三方认证与加密

- `security_settings` 表提供全局开关，包含 `ThirdPartyAPIEnabled`、`APIKey`、`APISecret`、`NonceTTLSeconds`，可通过管理端 `GET/PATCH /security-settings` 接口调整。
- 中间件通过 `X-ZNP-API-Key`、`X-ZNP-Timestamp`、`X-ZNP-Nonce`、`X-ZNP-Signature` 校验请求。
- 当 `X-ZNP-Encrypted: true` 时，请求体需要使用 `api_secret` 派生的 AES-256-GCM 加密，IV 通过 `X-ZNP-IV` 传递。

## 业务扩展方向

1. **套餐售卖流程**：已实现余额与外部支付并行的下单流程（含流水记录、回调处理），后续可扩展续费、套餐升级与更多支付渠道。
2. **公告推送渠道**：结合 Webhook、邮件通知，将公告同步到外部 IM 渠道。
3. **余额充值**：配合支付网关实现充值、退款、自动开票功能。
4. **审计日志**：记录套餐、公告、节点变更的操作明细，满足审计与回溯需求。
