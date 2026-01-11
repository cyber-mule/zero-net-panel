# Zero Network Panel API 参考

> 本文档基于 `internal/handler/routes.go` 与 `internal/types` 整理，描述当前实现的接口与字段。
> 自动生成的 Markdown 文档可通过 `./scripts/gen-api-docs.sh` 生成，输出到 `docs/api-generated/`。
> 时间字段统一为 Unix 秒（UTC）。以下字段未标注 `可选` 即默认必填，实际校验以服务端错误提示为准。

## 基础信息

- Base URL：`http(s)://<host>:<port>/api/v1`
- 管理端前缀：`/api/v1/{adminPrefix}`，默认 `admin`，由 `Admin.RoutePrefix` 配置。
- 内容类型：`Content-Type: application/json`

## 鉴权

- 登录：`POST /api/v1/auth/login` 获取 `access_token` 与 `refresh_token`。
- 刷新：`POST /api/v1/auth/refresh` 换取新令牌。
- 注册：`POST /api/v1/auth/register` 创建账号，若要求验证会返回 `requires_verification=true`。
- 验证：`POST /api/v1/auth/verify` 使用验证码激活账号并返回令牌。
- 找回：`POST /api/v1/auth/forgot` 获取验证码，`POST /api/v1/auth/reset` 完成重置。
- 鉴权方式：`Authorization: Bearer <access_token>`
- 角色约束：
  - 管理端接口需要 `admin` 角色。
  - 用户端接口需要 `user` 角色。

## 错误响应

- 业务错误返回：`{"message": "..."}`，并带有对应 HTTP 状态码。
- 常见状态码：
  - `400` 参数非法
  - `401` 未登录或令牌失效
  - `403` 权限不足或访问受限
  - `404` 资源不存在
  - `409` 冲突（并发/状态不允许）
  - `429` 超出速率限制（管理端 IP 限流/验证码频控）
  - `500` 未捕获错误

## 第三方签名与加密（可选）

当 `security_settings.third_party_api_enabled = true` 且 `api_key/api_secret` 生效时，`/api/v1/user/*` 接口需要签名校验。

必填头部：

- `X-ZNP-API-Key`
- `X-ZNP-Timestamp`（Unix 秒）
- `X-ZNP-Nonce`（随机字符串）
- `X-ZNP-Signature`

签名规则：

```
METHOD\nPATH\nRAW_QUERY\nTIMESTAMP\nNONCE\nBASE64(BODY)
```

- `METHOD` 为大写 HTTP 方法。
- `PATH` 为请求路径（不含 host）。
- `RAW_QUERY` 为原始查询串（无则为空）。
- `BODY` 为原始请求体（空 body 也需要参与签名）。
- 使用 `HMAC-SHA256` 以 `api_secret` 计算，结果 Base64 编码后填入 `X-ZNP-Signature`。

可选加密：

- 头部：`X-ZNP-Encrypted: true`、`X-ZNP-IV: <base64>`
- 算法：AES-256-GCM
- key：`SHA256(api_secret)`
- body 为密文，服务端会在验签通过后解密。

## 分页与通用字段

- `page`、`per_page` 默认 `1/20`，最大 `100`。
- 分页响应：
  - `page`、`per_page`、`total_count`、`has_next`、`has_prev`

## 通用数据结构

### PaginationMeta

- `page` int
- `per_page` int
- `total_count` int64
- `has_next` bool
- `has_prev` bool

### AuthenticatedUser

- `id` uint64
- `email` string
- `display_name` string
- `roles` []string
- `created_at` int64
- `updated_at` int64

### BalanceSnapshot

- `user_id` uint64
- `balance_cents` int64
- `currency` string
- `updated_at` int64

### BalanceTransactionSummary

- `id` uint64
- `entry_type` string
- `amount_cents` int64
- `currency` string
- `balance_after_cents` int64
- `reference` string
- `description` string
- `metadata` object
- `created_at` int64

### CouponSummary

- `id` uint64
- `code` string
- `name` string
- `description` string
- `status` string（示例：`active`、`disabled`）
- `discount_type` string（`percent` 或 `fixed`）
- `discount_value` int64（percent 为 0~10000，fixed 为分单位）
- `currency` string（fixed 折扣必填）
- `max_redemptions` int
- `max_redemptions_per_user` int
- `min_order_cents` int64
- `starts_at` int64（可选）
- `ends_at` int64（可选）
- `created_at` int64
- `updated_at` int64

### OrderItem

- `id` uint64
- `order_id` uint64
- `item_type` string（示例：`plan`、`discount`）
- `item_id` uint64
- `name` string
- `quantity` int
- `unit_price_cents` int64
- `currency` string
- `subtotal_cents` int64
- `metadata` object（优惠券折扣条目包含 `coupon_id`、`coupon_code`、`discount_type`、`discount_value`）
- `created_at` int64

### OrderRefund

- `id` uint64
- `order_id` uint64
- `amount_cents` int64
- `reason` string
- `reference` string
- `metadata` object
- `created_at` int64

### OrderPayment

- `id` uint64
- `order_id` uint64
- `provider` string
- `method` string
- `intent_id` string
- `reference` string
- `status` string
- `amount_cents` int64
- `currency` string
- `failure_code` string
- `failure_message` string
- `metadata` object（常见字段：`pay_url`、`qr_code`、`gateway_intent_id`、`notify_url`、`return_url`）
- `created_at` int64
- `updated_at` int64

### OrderDetail

- `id` uint64
- `number` string
- `user_id` uint64
- `status` string（示例：`pending_payment`、`paid`、`payment_failed`、`cancelled`、`partially_refunded`、`refunded`）
- `payment_status` string（示例：`pending`、`succeeded`、`failed`）
- `payment_method` string（示例：`balance`、`external`、`manual`）
- `payment_intent_id` string（可选）
- `payment_reference` string（可选）
- `payment_failure_code` string（可选）
- `payment_failure_message` string（可选）
- `total_cents` int64
- `refunded_cents` int64
- `currency` string
- `plan_id` uint64（可选）
- `plan_snapshot` object（可选）
- `metadata` object（可选）
- `paid_at` int64（可选）
- `cancelled_at` int64（可选）
- `refunded_at` int64（可选）
- `created_at` int64
- `updated_at` int64
- `items` []OrderItem
- `refunds` []OrderRefund
- `payments` []OrderPayment

## 接口参考

### 系统

#### GET /api/v1/ping

- 说明：健康检查
- 响应：`PingResponse`
  - `status` string
  - `service` string
  - `version` string
  - `site_name` string
  - `logo_url` string
  - `timestamp` int64

### 认证

#### POST /api/v1/auth/login

- 说明：用户登录并获取访问令牌
- 请求体：
  - `email` string
  - `password` string
- 响应：
  - `access_token` string
  - `refresh_token` string
  - `token_type` string
  - `expires_in` int64
  - `refresh_expires_in` int64
  - `user` AuthenticatedUser

#### POST /api/v1/auth/refresh

- 说明：刷新访问令牌
- 请求体：
  - `refresh_token` string
- 响应：同 `auth/login`

#### POST /api/v1/auth/register

- 说明：注册账号
- 请求体：
  - `email` string
  - `password` string
  - `display_name` string（可选）
  - `invite_code` string（可选）
- 备注：当 `Auth.Registration.InviteOnly=true` 时，必须提供 `invite_code`；缺失返回 `400`，未命中白名单返回 `403`。
- 响应：
  - `requires_verification` bool
  - `access_token` string（可选）
  - `refresh_token` string（可选）
  - `token_type` string（可选）
  - `expires_in` int64（可选）
  - `refresh_expires_in` int64（可选）
  - `user` AuthenticatedUser

#### POST /api/v1/auth/verify

- 说明：邮箱验证码验证
- 请求体：
  - `email` string
  - `code` string
- 响应：同 `auth/login`

#### POST /api/v1/auth/forgot

- 说明：发送密码重置验证码
- 请求体：
  - `email` string
- 响应：
  - `message` string

#### POST /api/v1/auth/reset

- 说明：使用验证码重置密码
- 请求体：
  - `email` string
  - `code` string
  - `password` string
- 响应：
  - `message` string

### 管理端（需要 admin 权限）

> 实际路径：`/api/v1/{adminPrefix}`

#### GET /api/v1/{adminPrefix}/dashboard

- 说明：获取后台模块清单
- 响应：
  - `modules` []AdminModule
    - `key` string
    - `name` string
    - `description` string
    - `icon` string
    - `route` string
    - `permissions` []string

#### GET /api/v1/{adminPrefix}/users

- 说明：用户列表
- 查询参数：`page`、`per_page`、`q`、`status`、`role`
- 响应：
  - `users` []AdminUserSummary
  - `pagination` PaginationMeta

AdminUserSummary 字段：

- `id`、`email`、`display_name`、`roles`、`status`
- `email_verified_at`（可选）、`failed_login_attempts`
- `locked_until`（可选）、`last_login_at`（可选）
- `created_at`、`updated_at`

#### POST /api/v1/{adminPrefix}/users

- 说明：创建用户
- 请求体：
  - `email` string
  - `password` string
  - `display_name` string（可选）
  - `roles` []string（可选）
  - `status` string（可选）
  - `email_verified` bool（可选）
- 响应：
  - `user` AdminUserSummary

#### PATCH /api/v1/{adminPrefix}/users/{id}/status

- 说明：更新用户状态
- 请求体：
  - `status` string（示例：`active`、`disabled`、`pending`）
- 响应：
  - `user` AdminUserSummary

#### PATCH /api/v1/{adminPrefix}/users/{id}/roles

- 说明：更新用户角色
- 请求体：
  - `roles` []string
- 响应：
  - `user` AdminUserSummary

#### POST /api/v1/{adminPrefix}/users/{id}/reset-password

- 说明：重置用户密码
- 请求体：
  - `password` string
- 响应：
  - `message` string

#### POST /api/v1/{adminPrefix}/users/{id}/force-logout

- 说明：强制下线用户
- 响应：
  - `message` string

#### POST /api/v1/{adminPrefix}/users/{id}/credentials/rotate

- 说明：手动轮换用户协议鉴权凭据
- 响应：
  - `user_id` uint64
  - `credential` CredentialSummary

#### GET /api/v1/{adminPrefix}/nodes

- 说明：节点列表
- 查询参数：`page`、`per_page`、`sort`、`direction`、`q`、`status`、`protocol`
- `sort` 可选：`name`、`region`、`last_synced_at`、`capacity_mbps`
- 响应：
  - `nodes` []NodeSummary
  - `pagination` PaginationMeta

NodeSummary 字段：

- `id`、`name`、`region`、`country`、`isp`、`status`、`tags`
- `capacity_mbps`、`description`、`access_address`、`control_endpoint`
- `status_sync_enabled`（是否允许节点状态自动同步）
- `last_synced_at`、`updated_at`
备注：
- `status` 为管理端维护字段，手动禁用时为 `disabled`；运行态健康度请看协议绑定健康状态。
- 当 `status_sync_enabled=true` 且能访问节点控制面时，服务会自动将 `status` 更新为 `online`/`offline`。
- 节点控制面必须配置 `control_endpoint`，不再回退全局 `Kernel.HTTP`。
- 控制面鉴权优先级：`control_access_key` + `control_secret_key` → `control_token`（无全局兜底）

#### POST /api/v1/{adminPrefix}/nodes

- 说明：创建节点
- 请求体：
  - `name` string
  - `region` string（可选）
  - `country` string（可选）
  - `isp` string（可选）
  - `status` string（可选）
  - `tags` []string（可选）
  - `capacity_mbps` int（可选）
  - `description` string（可选）
  - `access_address` string（可选，客户端对外地址）
  - `control_endpoint` string（必填，节点控制面地址）
  - `control_access_key` string（可选，节点控制面 AK，写入不回显）
  - `control_secret_key` string（可选，节点控制面 SK，写入不回显）
  - `ak` string（可选，兼容字段，等同 control_access_key）
  - `sk` string（可选，兼容字段，等同 control_secret_key）
  - `control_token` string（可选，节点控制面鉴权 token，写入不回显）
  - `status_sync_enabled` bool（可选，是否允许节点状态自动同步，默认 true）
- 响应：
  - `node` NodeSummary
注：`control_token` 可直接填写 `Basic <base64(ak:sk)>` 或 `Bearer <token>`，无前缀按 `Bearer` 处理。
- 示例请求体：
```json
{
  "name": "hk-edge-1",
  "region": "hk",
  "country": "HK",
  "isp": "HKT",
  "status": "online",
  "tags": ["edge"],
  "capacity_mbps": 1000,
  "description": "HK edge",
  "access_address": "hk.example.com",
  "control_endpoint": "https://kernel-hk.example.com/api"
}
```

#### PATCH /api/v1/{adminPrefix}/nodes/{id}

- 说明：更新节点
- 路径参数：`id` uint64
- 请求体：
  - `name` string（可选）
  - `region` string（可选）
  - `country` string（可选）
  - `isp` string（可选）
  - `status` string（可选）
  - `tags` []string（可选）
  - `capacity_mbps` int（可选）
  - `description` string（可选）
  - `access_address` string（可选，客户端对外地址）
  - `control_endpoint` string（可选，节点控制面地址）
  - `control_access_key` string（可选，节点控制面 AK，写入不回显）
  - `control_secret_key` string（可选，节点控制面 SK，写入不回显）
  - `ak` string（可选，兼容字段，等同 control_access_key）
  - `sk` string（可选，兼容字段，等同 control_secret_key）
  - `control_token` string（可选，节点控制面鉴权 token，写入不回显）
  - `status_sync_enabled` bool（可选，是否允许节点状态自动同步）
- 响应：
  - `node` NodeSummary
- 示例请求体：
```json
{
  "status": "maintenance",
  "tags": ["edge", "maintenance"],
  "capacity_mbps": 500
}
```

#### DELETE /api/v1/{adminPrefix}/nodes/{id}

- 说明：删除节点（软删除，同时清理关联协议绑定与内核记录）
- 路径参数：`id` uint64
- 响应：`204 No Content`

#### POST /api/v1/{adminPrefix}/nodes/{id}/disable

- 说明：禁用节点
- 路径参数：`id` uint64
- 响应：
  - `node` NodeSummary

#### GET /api/v1/{adminPrefix}/nodes/{id}/kernels

- 说明：节点内核配置列表
- 路径参数：`id` uint64
- 响应：
  - `node_id` uint64
  - `kernels` []NodeKernelSummary

NodeKernelSummary 字段：

- `protocol`、`endpoint`、`revision`、`status`、`config`、`last_synced_at`
- 备注：该接口返回内核端点与配置，属于管理端敏感信息。

#### POST /api/v1/{adminPrefix}/nodes/{id}/kernels/sync

- 说明：触发节点与内核同步
- 路径参数：`id` uint64
- 请求体：
  - `protocol` string（可选，空表示同步默认协议；当前仅支持 `http`）
- 响应：
  - `node_id` uint64
  - `protocol` string
  - `revision` string
  - `synced_at` int64
  - `message` string

#### POST /api/v1/{adminPrefix}/nodes/status/sync

- 说明：手动触发节点状态同步（仅同步指定节点）
- 请求体：
  - `node_ids` []uint64（必填，节点 ID 列表）
- 响应：
  - `results` []NodeStatusSyncResult

NodeStatusSyncResult 字段：

- `node_id`、`status`、`message`、`synced_at`
- `status` 可能为 `online` / `offline` / `skipped` / `error`
- `skipped` 表示节点已 `disabled`
- `error` 表示节点不存在或控制面地址缺失

#### GET /api/v1/{adminPrefix}/protocol-entries

- 说明：协议发布列表（对外入口）
- 查询参数：`page`、`per_page`、`sort`、`direction`、`q`、`protocol`、`status`、`binding_id`
- 响应：
  - `entries` []ProtocolEntrySummary
  - `pagination` PaginationMeta

ProtocolEntrySummary 字段：

- `id`、`name`、`binding_id`、`binding_name`、`node_id`、`node_name`
- `protocol`、`status`、`binding_status`、`health_status`
- `entry_address`、`entry_port`、`tags`、`description`、`profile`
- `created_at`、`updated_at`

说明：
- `entry_address/entry_port` 为对外入口地址，可与绑定监听不一致。
- `status` 仅影响用户可见性；`binding_status`/`health_status` 来自绑定健康状态。

#### POST /api/v1/{adminPrefix}/protocol-entries

- 说明：创建协议发布
- 请求体：
  - `binding_id` uint64
  - `entry_address` string
  - `entry_port` int
  - `protocol` string（可选，默认继承绑定协议）
  - `status` string（可选）
  - `tags` []string（可选）
  - `description` string（可选）
  - `profile` map（可选，对外公开配置）
- 响应：
  - ProtocolEntrySummary

#### PATCH /api/v1/{adminPrefix}/protocol-entries/{id}

- 说明：更新协议发布
- 路径参数：`id` uint64
- 请求体：同创建（均可选）
- 响应：
  - ProtocolEntrySummary

#### DELETE /api/v1/{adminPrefix}/protocol-entries/{id}

- 说明：删除协议发布
- 路径参数：`id` uint64
- 响应：204

#### GET /api/v1/{adminPrefix}/protocol-bindings

- 说明：协议绑定列表
- 查询参数：`page`、`per_page`、`sort`、`direction`、`q`、`status`、`protocol`、`node_id`
- 响应：
  - `bindings` []ProtocolBindingSummary
  - `pagination` PaginationMeta

ProtocolBindingSummary 字段：

- `id`、`name`、`node_id`、`node_name`、`protocol`
- `role`、`listen`、`connect`、`access_port`、`status`、`kernel_id`（字符串）
- `kernel_id` 需与内核侧协议 ID 一致，通常不是数字
- `sync_status`、`health_status`、`last_synced_at`、`last_heartbeat_at`、`last_sync_error`
- `tags`、`description`、`profile`、`metadata`
- `created_at`、`updated_at`

说明：
- `listen` 为空或仅端口时，会用 `access_port` 归一化为 `0.0.0.0:<port>` 供内核使用。

#### POST /api/v1/{adminPrefix}/protocol-bindings

- 说明：创建协议绑定
- 请求体：
  - `node_id` uint64
  - `protocol` string
  - `role` string
  - `profile` map（必填，内核实际配置）
  - `listen` string（可选）
  - `connect` string（可选）
  - `access_port` int（可选，内核监听端口）
  - `status` string（可选）
  - `kernel_id` string（必填，内核协议标识，通常为字符串）
  - `tags` []string（可选）
  - `description` string（可选）
  - `metadata` map（可选）
- 响应：
  - ProtocolBindingSummary

#### PATCH /api/v1/{adminPrefix}/protocol-bindings/{id}

- 说明：更新协议绑定
- 路径参数：`id` uint64
- 请求体：同创建（均可选）
- 响应：
  - ProtocolBindingSummary

#### DELETE /api/v1/{adminPrefix}/protocol-bindings/{id}

- 说明：删除协议绑定
- 路径参数：`id` uint64
- 响应：204

#### POST /api/v1/{adminPrefix}/protocol-bindings/{id}/sync

- 说明：同步单条协议绑定
- 路径参数：`id` uint64
- 响应：
  - ProtocolBindingSyncResult

#### POST /api/v1/{adminPrefix}/protocol-bindings/sync

- 说明：批量同步协议绑定
- 请求体：
  - `binding_ids` []uint64（可选）
  - `node_ids` []uint64（可选）
- 响应：
  - `results` []ProtocolBindingSyncResult

#### POST /api/v1/{adminPrefix}/protocol-bindings/status/sync

- 说明：手动反向同步协议健康状态
- 请求体：
  - `node_ids` []uint64（必填，节点 ID 列表）
- 响应：
  - `results` []ProtocolBindingStatusSyncResult

ProtocolBindingStatusSyncResult 字段：

- `node_id`、`status`、`message`、`synced_at`、`updated`
- `status` 可能为 `synced` / `error` / `skipped`

#### GET /api/v1/{adminPrefix}/subscriptions

- 说明：订阅列表
- 查询参数：`page`、`per_page`、`q`、`status`、`user_id`、`plan_name`、`plan_id`、`template_id`
- 响应：
  - `subscriptions` []AdminSubscriptionSummary
  - `pagination` PaginationMeta

AdminSubscriptionUserSummary 字段：

- `id`、`email`、`display_name`

AdminSubscriptionSummary 字段：

- `id`、`user`
- `name`、`plan_name`、`plan_id`、`plan_snapshot`、`status`
- `template_id`、`available_template_ids`
- `token`、`expires_at`
- `traffic_total_bytes`、`traffic_used_bytes`
- `devices_limit`、`last_refreshed_at`
- `created_at`、`updated_at`

#### GET /api/v1/{adminPrefix}/subscriptions/{id}

- 说明：订阅详情
- 路径参数：`id` uint64
- 响应：
  - `subscription` AdminSubscriptionSummary

#### POST /api/v1/{adminPrefix}/subscriptions

- 说明：创建订阅
- 请求体：
  - `user_id` uint64
  - `name` string
  - `plan_name` string（可选）
  - `plan_id` uint64
  - `status` string（可选）
  - `template_id` uint64
  - `available_template_ids` []uint64（可选）
  - `token` string（可选）
  - `expires_at` int64
  - `traffic_total_bytes` int64
  - `traffic_used_bytes` int64（可选）
  - `devices_limit` int
- 响应：
  - `subscription` AdminSubscriptionSummary

#### PATCH /api/v1/{adminPrefix}/subscriptions/{id}

- 说明：更新订阅
- 路径参数：`id` uint64
- 请求体（字段均可选）：
  - `name`、`plan_name`、`plan_id`、`status`
  - `template_id`、`available_template_ids`
  - `token`、`expires_at`
  - `traffic_total_bytes`、`traffic_used_bytes`
  - `devices_limit`
- 响应：
  - `subscription` AdminSubscriptionSummary

#### POST /api/v1/{adminPrefix}/subscriptions/{id}/disable

- 说明：禁用订阅
- 路径参数：`id` uint64
- 请求体：
  - `reason` string（可选）
- 响应：
  - `subscription` AdminSubscriptionSummary

#### POST /api/v1/{adminPrefix}/subscriptions/{id}/extend

- 说明：延长订阅有效期（`extend_days`/`extend_hours` 与 `expires_at` 二选一）
- 路径参数：`id` uint64
- 请求体：
  - `extend_days` int（可选）
  - `extend_hours` int（可选）
  - `expires_at` int64（可选）
- 响应：
  - `subscription` AdminSubscriptionSummary

#### GET /api/v1/{adminPrefix}/subscription-templates

- 说明：订阅模板列表
- 查询参数：`page`、`per_page`、`sort`、`direction`、`q`、`client_type`、`format`、`include_drafts`
- `sort` 可选：`name`、`client_type`、`version`、`created_at`
- 响应：
  - `templates` []SubscriptionTemplateSummary
  - `pagination` PaginationMeta

TemplateVariable 字段：

- `value_type` string
- `required` bool
- `description` string
- `default_value` interface{}

SubscriptionTemplateSummary 字段：

- `id`、`name`、`description`、`client_type`、`format`
- `content`（可选）
- `variables` map[string]TemplateVariable
- `is_default` bool
- `version` uint32
- `updated_at` int64
- `published_at` int64
- `last_published_by` string

#### POST /api/v1/{adminPrefix}/subscription-templates

- 说明：创建订阅模板
- 请求体：
  - `name` string
  - `description` string（可选）
  - `client_type` string
  - `format` string
  - `content` string
  - `variables` map[string]TemplateVariable（可选）
  - `is_default` bool（可选）
- 响应：SubscriptionTemplateSummary

#### PATCH /api/v1/{adminPrefix}/subscription-templates/{id}

- 说明：更新订阅模板
- 路径参数：`id` uint64
- 请求体：
  - `name` string（可选）
  - `description` string（可选）
  - `format` string（可选）
  - `content` string（可选）
  - `variables` map[string]TemplateVariable（可选）
  - `is_default` bool（可选）
- 响应：SubscriptionTemplateSummary

#### POST /api/v1/{adminPrefix}/subscription-templates/{id}/publish

- 说明：发布订阅模板
- 路径参数：`id` uint64
- 请求体：
  - `changelog` string（可选）
  - `operator` string（可选）
- 响应：
  - `template` SubscriptionTemplateSummary
  - `history` SubscriptionTemplateHistoryEntry

SubscriptionTemplateHistoryEntry 字段：

- `version` uint32
- `changelog` string
- `published_at` int64
- `published_by` string
- `variables` map[string]TemplateVariable

#### GET /api/v1/{adminPrefix}/subscription-templates/{id}/history

- 说明：查看模板发布历史
- 路径参数：`id` uint64
- 响应：
  - `template_id` uint64
  - `history` []SubscriptionTemplateHistoryEntry

#### GET /api/v1/{adminPrefix}/plans

- 说明：套餐列表
- 查询参数：`page`、`per_page`、`sort`、`direction`、`q`、`status`、`visible`
- `sort` 可选：`price`、`name`、`updated`
- 响应：
  - `plans` []PlanSummary
  - `pagination` PaginationMeta

PlanSummary 字段：

- `id`、`name`、`slug`、`description`、`tags`、`features`
- `binding_ids`
- `billing_options`
- `price_cents`、`currency`、`duration_days`
- `traffic_limit_bytes`、`traffic_multipliers`、`devices_limit`
- `sort_order`、`status`、`visible`
- `created_at`、`updated_at`

PlanBillingOptionSummary 字段：

- `id`、`plan_id`、`name`
- `duration_value`、`duration_unit`
- `price_cents`、`currency`
- `sort_order`、`status`、`visible`
- `created_at`、`updated_at`

#### POST /api/v1/{adminPrefix}/plans

- 说明：创建套餐
- 请求体：
  - `name` string
  - `slug` string（可选）
  - `description` string（可选）
  - `tags` []string（可选）
  - `features` []string（可选）
  - `binding_ids` []uint64（可选，套餐绑定的协议）
  - `price_cents` int64
  - `currency` string
  - `duration_days` int
  - `traffic_limit_bytes` int64（可选）
  - `traffic_multipliers` map（可选，协议流量倍数）
  - `devices_limit` int（可选）
  - `sort_order` int（可选）
  - `status` string（可选，默认 draft）
  - `visible` bool（可选）
- 响应：PlanSummary

#### PATCH /api/v1/{adminPrefix}/plans/{id}

- 说明：更新套餐
- 路径参数：`id` uint64
- 请求体（字段均可选）：
  - `name`、`slug`、`description`、`tags`、`features`、`binding_ids`
  - `price_cents`、`currency`、`duration_days`
  - `traffic_limit_bytes`、`traffic_multipliers`、`devices_limit`
  - `sort_order`、`status`、`visible`
- 响应：PlanSummary

#### GET /api/v1/{adminPrefix}/plans/{plan_id}/billing-options

- 说明：套餐计费选项列表
- 路径参数：`plan_id` uint64
- 查询参数：`status`（可选）、`visible`（可选）
- 响应：
  - `options` []PlanBillingOptionSummary

#### POST /api/v1/{adminPrefix}/plans/{plan_id}/billing-options

- 说明：创建套餐计费选项
- 路径参数：`plan_id` uint64
- 请求体：
  - `name` string（可选）
  - `duration_value` int
  - `duration_unit` string（hour/day/month/year）
  - `price_cents` int64
  - `currency` string（可选）
  - `sort_order` int（可选）
  - `status` string（可选，默认 draft）
  - `visible` bool（可选）
- 响应：PlanBillingOptionSummary

#### PATCH /api/v1/{adminPrefix}/plans/{plan_id}/billing-options/{id}

- 说明：更新套餐计费选项
- 路径参数：`plan_id` uint64、`id` uint64
- 请求体（字段均可选）：
  - `name`、`duration_value`、`duration_unit`
  - `price_cents`、`currency`
  - `sort_order`、`status`、`visible`
- 响应：PlanBillingOptionSummary

#### GET /api/v1/{adminPrefix}/coupons

- 说明：优惠券列表
- 查询参数：`page`、`per_page`、`q`、`status`、`sort`、`direction`
- `sort` 可选：`code`、`status`、`created_at`、`updated_at`、`starts_at`、`ends_at`
- 响应：
  - `coupons` []CouponSummary
  - `pagination` PaginationMeta

#### POST /api/v1/{adminPrefix}/coupons

- 说明：创建优惠券
- 请求体：
  - `code` string
  - `name` string
  - `description` string（可选）
  - `status` string（可选，默认 active）
  - `discount_type` string（percent 或 fixed）
  - `discount_value` int64
  - `currency` string（可选，fixed 折扣必填）
  - `max_redemptions` int（可选）
  - `max_redemptions_per_user` int（可选）
  - `min_order_cents` int64（可选）
  - `starts_at` int64（可选）
  - `ends_at` int64（可选）
- 响应：CouponSummary

#### PATCH /api/v1/{adminPrefix}/coupons/{id}

- 说明：更新优惠券
- 路径参数：`id` uint64
- 请求体（字段均可选）：同创建接口字段
- 响应：CouponSummary

#### DELETE /api/v1/{adminPrefix}/coupons/{id}

- 说明：删除优惠券
- 路径参数：`id` uint64
- 响应：`{"message":"ok"}`

#### GET /api/v1/{adminPrefix}/payment-channels

- 说明：支付通道列表
- 查询参数：`page`、`per_page`、`q`、`provider`、`enabled`、`sort`、`direction`
- `sort` 可选：`name`、`created`、`updated`
- 响应：
  - `channels` []PaymentChannelSummary
  - `pagination` PaginationMeta

PaymentChannelSummary 字段：

- `id`、`name`、`code`、`provider`
- `enabled`、`sort_order`、`config`
- `created_at`、`updated_at`

支付通道 `config`（外部支付发起）示例：

```json
{
  "mode": "http",
  "notify_url": "https://example.com/api/v1/payments/callback?order_id={{order_id}}&payment_id={{payment_id}}",
  "return_url": "https://example.com/orders/{{order_number}}",
  "http": {
    "endpoint": "https://gateway.example.com/pay",
    "method": "POST",
    "body_type": "json",
    "headers": {
      "Content-Type": "application/json"
    },
    "payload": {
      "order_no": "{{order_number}}",
      "amount": "{{amount}}",
      "notify_url": "{{notify_url}}",
      "return_url": "{{return_url}}"
    }
  },
  "response": {
    "pay_url": "data.pay_url",
    "qr_code": "data.qr_code",
    "reference": "data.reference"
  },
  "webhook": {
    "signature_type": "hmac_sha256",
    "signature_header": "X-Pay-Signature",
    "secret": "your-signing-secret"
  },
  "refund": {
    "http": {
      "endpoint": "https://gateway.example.com/refund",
      "method": "POST",
      "body_type": "json",
      "payload": {
        "payment_ref": "{{payment_reference}}",
        "amount": "{{refund_amount}}",
        "reason": "{{refund_reason}}"
      }
    },
    "response": {
      "reference": "data.refund_id",
      "status": "data.status"
    },
    "status_map": {
      "success": "succeeded",
      "failed": "failed"
    }
  },
  "reconcile": {
    "http": {
      "endpoint": "https://gateway.example.com/query",
      "method": "POST",
      "body_type": "json",
      "payload": {
        "payment_ref": "{{payment_reference}}"
      }
    },
    "response": {
      "status": "data.status",
      "reference": "data.reference"
    },
    "status_map": {
      "paid": "succeeded",
      "failed": "failed",
      "processing": "pending"
    }
  }
}
```

`notify_url`/`return_url`/`payload` 支持模板变量：`{{order_id}}`、`{{order_number}}`、`{{order_status}}`、`{{payment_id}}`、`{{payment_intent_id}}`、`{{payment_reference}}`、`{{payment_status}}`、`{{amount_cents}}`、`{{amount}}`、`{{currency}}`、`{{user_id}}`、`{{plan_id}}`、`{{plan_name}}`、`{{quantity}}`、`{{payment_channel}}`、`{{payment_provider}}`、`{{refund_amount_cents}}`、`{{refund_amount}}`、`{{refund_reason}}`。

`response` 字段支持点路径（如 `data.pay_url`），`pay_url` 设为 `$` 可直接使用原始响应体字符串。

`webhook` 签名默认使用 `hmac_sha256`，签名体为原始回调请求体（body）。

外部支付联调示例见 `docs/payment-gateway-demo.md`。

#### GET /api/v1/{adminPrefix}/payment-channels/{id}

- 说明：支付通道详情
- 路径参数：`id` uint64
- 响应：PaymentChannelSummary

#### POST /api/v1/{adminPrefix}/payment-channels

- 说明：创建支付通道
- 请求体：
  - `name` string
  - `code` string
  - `provider` string（可选）
  - `enabled` bool（可选）
  - `sort_order` int（可选）
  - `config` object（可选）
- 响应：PaymentChannelSummary

#### PATCH /api/v1/{adminPrefix}/payment-channels/{id}

- 说明：更新支付通道
- 路径参数：`id` uint64
- 请求体：
  - `name` string（可选）
  - `code` string（可选）
  - `provider` string（可选）
  - `enabled` bool（可选）
  - `sort_order` int（可选）
  - `config` object（可选）
- 响应：PaymentChannelSummary

#### GET /api/v1/{adminPrefix}/announcements

- 说明：公告列表
- 查询参数：`page`、`per_page`、`status`、`category`、`audience`、`q`、`sort`、`direction`
- `sort` 可选：`created`、`title`、`priority`
- 响应：
  - `announcements` []AnnouncementSummary
  - `pagination` PaginationMeta

AnnouncementSummary 字段：

- `id`、`title`、`content`、`category`、`status`、`audience`
- `is_pinned`、`priority`
- `visible_from`、`visible_to`（可选）
- `published_at`（可选）
- `published_by`、`created_by`、`updated_by`
- `created_at`、`updated_at`

#### POST /api/v1/{adminPrefix}/announcements

- 说明：创建公告
- 请求体：
  - `title` string
  - `content` string
  - `category` string（可选）
  - `audience` string（可选）
  - `is_pinned` bool（可选）
  - `priority` int（可选）
  - `created_by` string（可选）
- 响应：AnnouncementSummary

#### POST /api/v1/{adminPrefix}/announcements/{id}/publish

- 说明：发布公告
- 路径参数：`id` uint64
- 请求体：
  - `visible_to` int64（可选）
  - `operator` string（可选）
- 响应：AnnouncementSummary

#### GET /api/v1/{adminPrefix}/site-settings

- 说明：查询站点配置
- 响应：
  - `setting` SiteSetting

SiteSetting 字段：

- `id`、`name`、`logo_url`、`access_domain`
- `created_at`、`updated_at`

#### PATCH /api/v1/{adminPrefix}/site-settings

- 说明：更新站点配置
- 请求体：
  - `name` string（可选）
  - `logo_url` string（可选）
  - `access_domain` string（可选）
- 响应：同 GET

#### GET /api/v1/{adminPrefix}/security-settings

- 说明：查询第三方安全配置
- 响应：
  - `setting` SecuritySetting

SecuritySetting 字段：

- `id`、`third_party_api_enabled`
- `api_key`、`api_secret`
- `encryption_algorithm`
- `nonce_ttl_seconds`
- `created_at`、`updated_at`

#### PATCH /api/v1/{adminPrefix}/security-settings

- 说明：更新第三方安全配置
- 请求体：
  - `third_party_api_enabled` bool（可选）
  - `api_key` string（可选）
  - `api_secret` string（可选）
  - `encryption_algorithm` string（可选）
  - `nonce_ttl_seconds` int（可选）
- 响应：同 GET

#### GET /api/v1/{adminPrefix}/audit-logs

- 说明：审计日志列表
- 查询参数：`page`、`per_page`、`actor_id`、`action`、`resource_type`、`resource_id`、`since`、`until`
- `since`/`until` 为 Unix 秒
- 响应：
  - `logs` []AuditLogSummary
  - `pagination` PaginationMeta

AuditLogSummary 字段：

- `id`、`actor_id`、`actor_email`、`actor_roles`
- `action`、`resource_type`、`resource_id`
- `source_ip`、`metadata`
- `created_at`

#### GET /api/v1/{adminPrefix}/audit-logs/export

- 说明：导出审计日志
- 查询参数：`page`、`per_page`、`actor_id`、`action`、`resource_type`、`resource_id`、`since`、`until`、`format`
- `format` 可选：`json`、`csv`（默认 `json`）
- `per_page` 导出上限为 5000（默认 1000）
- 响应：
  - `json`：`logs` []AuditLogSummary + `total_count` + `exported_at`
  - `csv`：CSV 文件下载

#### GET /api/v1/{adminPrefix}/orders

- 说明：订单列表
- 查询参数：`page`、`per_page`、`status`、`payment_method`、`payment_status`、`number`、`sort`、`direction`、`user_id`
- `sort` 可选：`updated`、`total`
- 响应：
  - `orders` []AdminOrderDetail
  - `pagination` PaginationMeta

AdminOrderDetail 字段：

- 字段同 OrderDetail，另含 `user`（`id`、`email`、`display_name`）

#### GET /api/v1/{adminPrefix}/orders/{id}

- 说明：订单详情
- 路径参数：`id` uint64
- 响应：
  - `order` AdminOrderDetail

#### POST /api/v1/{adminPrefix}/orders/{id}/pay

- 说明：人工标记订单已支付
- 路径参数：`id` uint64
- 请求体：
  - `payment_method` string（可选，线下支付可用 `manual`）
  - `paid_at` int64（可选）
  - `note` string（可选）
  - `reference` string（可选）
  - `charge_balance` bool（可选，是否影响余额）
- 响应：
  - `order` AdminOrderDetail

#### POST /api/v1/{adminPrefix}/orders/{id}/cancel

- 说明：取消订单
- 路径参数：`id` uint64
- 请求体：
  - `reason` string（可选）
  - `cancelled_at` int64（可选）
- 响应：
  - `order` AdminOrderDetail

#### POST /api/v1/{adminPrefix}/orders/{id}/refund

- 说明：退款（余额或外部支付）
- 路径参数：`id` uint64
- 请求体：
  - `amount_cents` int64
  - `reason` string（可选）
  - `metadata` object（可选）
  - `refund_at` int64（可选）
  - `credit_balance` bool（可选）
- 外部支付说明：
  - 订单为 `payment_method=external` 时，会按支付通道 `config.refund` 发起退款。
  - 回调验签由通道 `config.webhook` 控制（不配置则不校验）。
- 响应：
  - `order` AdminOrderDetail

#### POST /api/v1/{adminPrefix}/orders/payments/reconcile

- 说明：外部支付对账
- 请求体：
  - `order_id` uint64
  - `payment_id` uint64
- 响应：
  - `order` AdminOrderDetail

#### POST /api/v1/{adminPrefix}/orders/payments/callback

- 说明：外部支付回调（Webhook 专用）
- 认证：`X-ZNP-Webhook-Token` 或 `Stripe-Signature`（取决于 `Webhook` 配置），或通道 `config.webhook` 签名
- 请求体：
  - `order_id` uint64
  - `payment_id` uint64
  - `status` string
  - `reference` string（可选）
  - `failure_code` string（可选）
  - `failure_message` string（可选）
  - `paid_at` int64（可选）
- 响应：
  - `order` AdminOrderDetail

#### POST /api/v1/payments/callback

- 说明：外部支付回调（免登录，Webhook 专用）
- 认证：`X-ZNP-Webhook-Token` 或 `Stripe-Signature`（取决于 `Webhook` 配置）
- 请求体：同 `/api/v1/{adminPrefix}/orders/payments/callback`
- 响应：同上

#### POST /api/v1/kernel/traffic

- 说明：内核流量回调（免登录，Webhook 专用）
- 认证：`X-ZNP-Webhook-Token`
- 请求体：`records` 数组，字段见 `KernelTrafficRecord`
- 响应：
  - `accepted`、`failed`

#### POST /api/v1/kernel/service-events

- 说明：内核服务事件回调（免登录，Webhook 专用）
- 认证：`X-ZNP-Webhook-Token`
- 请求体：`event` + `payload`（如 `user.traffic.reported`，payload 包含 `user_id` 与 `current.used`/`current.remaining`）
- 备注：面板侧优先使用 `subscription_id`，否则使用 `user_id` 更新订阅已用流量
- 响应：
  - `status`
  - `accepted`、`failed`（当事件为 `user.traffic.reported`）

#### POST /api/v1/kernel/events

- 说明：内核节点事件回调（免登录，Webhook 专用）
- 认证：`X-ZNP-Webhook-Token`
- 请求体：`event`、`id`/`node_id`、`status`、`observed_at`、`message`
- 响应：
  - `status`

### 公共订阅（免登录）

#### GET /api/v1/subscriptions/{token}

- 说明：客户端订阅拉取（免登录）
- 路径参数：`token` string
- 响应：**非 JSON**，直接返回订阅内容
  - `Content-Type`：`text/plain` 或 `application/json`（取决于模板格式）
  - `ETag`：内容哈希
- 规则：
  - 仅 `status=active` 且未过期的订阅可拉取
  - `User-Agent` 关键词匹配客户端类型，忽略大小写；命中后优先选择对应 `client_type` 的默认模板
  - 未命中则回退订阅默认模板

### 用户端（需要 user 权限）

#### GET /api/v1/user/subscriptions

- 说明：订阅列表
- 查询参数：`page`、`per_page`、`sort`、`direction`、`q`、`status`
- `sort` 可选：`name`、`plan_name`、`status`、`expires_at`、`created_at`
- 说明：
  - 用户侧默认不返回 `disabled` 状态订阅
  - `expired` 状态仍会返回，便于续费
- 响应：
  - `subscriptions` []UserSubscriptionSummary
  - `pagination` PaginationMeta

UserSubscriptionSummary 字段：

- `id`、`name`、`plan_name`、`plan_id`、`status`
- `template_id`、`available_template_ids`
- `expires_at`、`traffic_total_bytes`、`traffic_used_bytes`
- `devices_limit`、`last_refreshed_at`

#### GET /api/v1/user/subscriptions/{id}/preview

- 说明：订阅预览
- 路径参数：`id` uint64
- 查询参数：`template_id`（可选）
- 说明：`disabled` 状态订阅返回 404
- 响应：
  - `subscription_id` uint64
  - `template_id` uint64
  - `content` string
  - `content_type` string
  - `etag` string
  - `generated_at` int64

#### POST /api/v1/user/subscriptions/{id}/template

- 说明：切换订阅模板
- 路径参数：`id` uint64
- 请求体：
  - `template_id` uint64
- 说明：`disabled` 状态订阅返回 404
- 响应：
  - `subscription_id` uint64
  - `template_id` uint64
  - `updated_at` int64

#### GET /api/v1/user/subscriptions/{id}/traffic

- 说明：订阅流量明细
- 路径参数：`id` uint64
- 查询参数：`page`、`per_page`、`protocol`、`node_id`、`binding_id`、`from`、`to`
- `from`/`to` 为 Unix 秒
- 说明：`disabled` 状态订阅返回 404
- 响应：
  - `summary` UserSubscriptionTrafficSummary
  - `records` []UserTrafficUsageRecord
  - `pagination` PaginationMeta

UserSubscriptionTrafficSummary 字段：

- `raw_bytes`、`charged_bytes`

UserTrafficUsageRecord 字段：

- `id`、`protocol`、`node_id`、`binding_id`
- `bytes_up`、`bytes_down`
- `raw_bytes`、`charged_bytes`、`multiplier`
- `observed_at`

#### GET /api/v1/user/plans

- 说明：可购买套餐列表
- 查询参数：`q`（可选）
- 响应：
  - `plans` []UserPlanSummary

UserPlanSummary 字段：

- `id`、`name`、`description`、`features`
- `billing_options`
- `price_cents`、`currency`、`duration_days`
- `traffic_limit_bytes`、`devices_limit`、`tags`

#### GET /api/v1/user/nodes

- 说明：用户侧节点运行状态列表（脱敏）
- 查询参数：`page`、`per_page`、`status`、`protocol`
- 响应：
  - `nodes` []UserNodeStatusSummary
  - `pagination` PaginationMeta

UserNodeStatusSummary 字段：

- `id`、`name`、`region`、`country`、`isp`、`status`
- `tags`、`capacity_mbps`、`description`
- `last_synced_at`、`updated_at`
- `kernel_statuses` []UserNodeKernelStatusSummary
- `protocol_statuses` []UserNodeProtocolStatusSummary

UserNodeKernelStatusSummary 字段：

- `protocol`、`status`、`last_synced_at`

UserNodeProtocolStatusSummary 字段：

- `binding_id`、`protocol`、`role`、`status`
- `health_status`、`last_heartbeat_at`

说明：
- `status` 与节点状态枚举一致（`online`/`offline`/`maintenance`/`disabled`）。
- `kernel_statuses.status` 表示同步记录状态（如 `synced`、`configured`）。

响应示例：
```json
{
  "nodes": [
    {
      "id": 1,
      "name": "hk-edge-1",
      "region": "hk",
      "country": "HK",
      "isp": "HKT",
      "status": "online",
      "tags": ["edge"],
      "capacity_mbps": 1000,
      "description": "HK edge",
      "last_synced_at": 1734001010,
      "updated_at": 1734001010,
      "kernel_statuses": [
        {"protocol": "vless", "status": "synced", "last_synced_at": 1734001010},
        {"protocol": "ss", "status": "synced", "last_synced_at": 1734001001}
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "per_page": 20,
    "total_count": 1,
    "has_next": false,
    "has_prev": false
  }
}
```

备注：不返回内核端点、Revision、配置等敏感信息；需要查看详细配置请使用管理端接口。
提示：`kernel_statuses` 来自最近一次同步记录，并非内核实时心跳。

#### GET /api/v1/user/account/profile

- 说明：用户资料
- 响应：
  - `profile` UserProfile

UserProfile 字段：

- `id`、`email`、`display_name`、`status`
- `email_verified_at`（可选）
- `created_at`、`updated_at`

#### PATCH /api/v1/user/account/profile

- 说明：更新用户资料
- 请求体：
  - `display_name` string
- 响应：
  - `profile` UserProfile

#### POST /api/v1/user/account/password

- 说明：用户自主改密
- 请求体：
  - `current_password` string
  - `new_password` string
- 响应：
  - `message` string

备注：密码策略由 `Auth.PasswordPolicy` 控制。
提示：修改密码会刷新 `token_invalid_before`，旧令牌需重新登录。

#### POST /api/v1/user/account/credentials/rotate

- 说明：手动轮换用户协议鉴权凭据
- 响应：
  - `credential` CredentialSummary

CredentialSummary 字段：

- `version`、`status`
- `issued_at`
- `deprecated_at`（可选）、`revoked_at`（可选）
- `last_seen_at`（可选）

#### POST /api/v1/user/account/email/code

- 说明：发送邮箱变更验证码
- 请求体：
  - `email` string
- 响应：
  - `message` string

#### POST /api/v1/user/account/email

- 说明：变更用户邮箱
- 请求体：
  - `email` string
  - `code` string
  - `password` string
- 响应：
  - `profile` UserProfile

#### GET /api/v1/user/account/balance

- 说明：用户余额与流水
- 查询参数：`page`、`per_page`、`entry_type`
- 响应：
  - `user_id` uint64
  - `balance_cents` int64
  - `currency` string
  - `updated_at` int64
  - `transactions` []BalanceTransactionSummary
  - `pagination` PaginationMeta

#### GET /api/v1/user/announcements

- 说明：有效公告列表
- 查询参数：`audience`（可选）、`limit`（可选，默认 20，最大 100）
- 响应：
  - `announcements` []UserAnnouncementSummary

UserAnnouncementSummary 字段：

- `id`、`title`、`content`、`category`、`audience`
- `is_pinned`、`priority`
- `visible_from`、`visible_to`（可选）
- `published_at`（可选）

#### GET /api/v1/user/payment-channels

- 说明：用户侧支付通道列表（仅返回启用通道）
- 查询参数：`provider`（可选）
- 响应：
  - `channels` []UserPaymentChannelSummary

UserPaymentChannelSummary 字段：

- `id`、`name`、`code`、`provider`
- `sort_order`、`config`

#### POST /api/v1/user/orders

- 说明：下单
- 请求体：
  - `plan_id` uint64
  - `billing_option_id` uint64（可选）
  - `quantity` int
  - `payment_method` string（可选，默认 `balance`；线下可用 `manual`）
  - `payment_channel` string（可选，外部支付通道）
  - `payment_return_url` string（可选）
  - `idempotency_key` string（可选，幂等键）
  - `coupon_code` string（可选）
- 外部支付说明：
  - `payment_method=external` 且金额大于 0 时，需传启用的 `payment_channel` 且通道 `config` 已配置网关发起信息。
  - 响应 `order.payments[].metadata` 将包含 `pay_url` 或 `qr_code`，用于跳转支付页或展示二维码。
- 优惠券说明：
  - 校验失败会返回 `400`（未启用/过期/次数超限/不满足最低金额）。
  - 命中优惠时，`order.metadata` 会附带 `coupon_code`、`coupon_id`、`discount_cents`，并追加 `item_type=discount` 的订单条目。
- 响应：
  - `order` OrderDetail
  - `balance` BalanceSnapshot
  - `transaction` BalanceTransactionSummary（可选，仅余额扣费时返回）

#### POST /api/v1/user/orders/{id}/cancel

- 说明：取消用户订单
- 路径参数：`id` uint64
- 请求体：
  - `reason` string（可选）
- 响应：
  - `order` OrderDetail
  - `balance` BalanceSnapshot

#### GET /api/v1/user/orders

- 说明：用户订单列表
- 查询参数：`page`、`per_page`、`status`、`payment_method`、`payment_status`、`number`、`sort`、`direction`
- `sort` 可选：`updated`、`total`
- 响应：
  - `orders` []OrderDetail
  - `pagination` PaginationMeta

#### GET /api/v1/user/orders/{id}

- 说明：订单详情
- 路径参数：`id` uint64
- 响应：
  - `order` OrderDetail
  - `balance` BalanceSnapshot
  - `transaction` BalanceTransactionSummary（可选）

#### GET /api/v1/user/orders/{id}/payment-status

- 说明：确认订单支付状态
- 路径参数：`id` uint64
- 响应：
  - `order_id` uint64
  - `status` string
  - `payment_status` string
  - `payment_method` string
  - `payment_intent_id` string（可选）
  - `payment_reference` string（可选）
  - `payment_failure_code` string（可选）
  - `payment_failure_message` string（可选）
  - `paid_at` int64（可选）
  - `cancelled_at` int64（可选）
  - `refunded_cents` int64
  - `refunded_at` int64（可选）
  - `updated_at` int64
