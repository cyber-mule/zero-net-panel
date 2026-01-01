# 前端项目开发指南

本文档面向需要对接 Zero Network Panel（ZNP）后端的前端团队，覆盖环境配置、鉴权策略、接口映射和常见坑位。

## 1. 项目定位与角色

- **管理端（Admin Console）**：节点、模板、套餐、公告、订单、第三方安全配置等运维/运营功能。
- **用户端（User Portal）**：订阅与套餐展示、公告、余额、订单购买与取消。

后端区分角色访问：管理端接口要求 `admin` 角色，用户端接口要求 `user` 角色。

## 2. 环境配置建议

建议在前端通过环境变量管理 API 地址与路由前缀：

- `API_BASE_URL`：如 `http://localhost:8888`
- `API_PREFIX`：固定为 `/api/v1`
- `ADMIN_PREFIX`：默认 `admin`（需要与后端 `Admin.RoutePrefix` 一致）

示例拼接规则：

- 管理端：`${API_BASE_URL}/api/v1/${ADMIN_PREFIX}`
- 用户端：`${API_BASE_URL}/api/v1/user`

## 3. API 客户端设计

推荐封装统一的请求层：

- 自动拼接 base URL 与前缀
- 自动注入 `Authorization: Bearer <token>`
- 全局处理 `401/403/429/5xx`
- 支持单飞刷新（避免并发刷新导致令牌覆盖）

示例（伪代码）：

```ts
async function request(url, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) };
  const token = getAccessToken();
  if (token) headers.Authorization = `Bearer ${token}`;

  const res = await fetch(`${API_BASE_URL}${url}`, { ...options, headers });
  if (res.status === 401) {
    await refreshTokenOnce();
    return request(url, options);
  }
  return res.json();
}
```

## 4. 鉴权与刷新

- 登录接口：`POST /api/v1/auth/login`
- 刷新接口：`POST /api/v1/auth/refresh`
- 注册接口：`POST /api/v1/auth/register`
- 验证接口：`POST /api/v1/auth/verify`
- 找回密码：`POST /api/v1/auth/forgot`、`POST /api/v1/auth/reset`

注册与验证流程由配置开关控制：

- `Auth.Registration.Enabled`：是否开放注册。
- `Auth.Registration.InviteOnly/InviteCodes`：邀请制与邀请码白名单。
- `Auth.Registration.DefaultRoles`：注册成功后的默认角色。
- `Auth.Registration.RequireEmailVerification`：是否要求邮箱验证。

> 启用邀请制时，前端注册表单需提供 `invite_code` 字段。

建议流程：

1. 登录成功后缓存 `access_token` 与 `refresh_token`。
2. 访问接口时注入 `Authorization` 头。
3. 遇到 `401` 时用 `refresh_token` 换新令牌，再重试一次原请求。

> 建议 `access_token` 存于内存（减少 XSS 风险），`refresh_token` 放在更安全的存储（如 HttpOnly Cookie 或受控存储）。

账号生命周期提示：

- 注册后若 `RequireEmailVerification=true`，返回 `requires_verification=true`，账号处于 `pending` 状态，需调用 `POST /auth/verify` 激活。
- 账号被禁用或锁定时会返回 `403`，需提示联系管理员或稍后重试。
- 账号被重置密码或强制下线时会返回 `401`，前端需清理令牌并引导重新登录。

验证码与前端对接提示：

- 注册后需展示邮箱验证码输入框并调用 `POST /auth/verify` 完成激活。
- 找回密码流程为 `POST /auth/forgot` 获取验证码，再用 `POST /auth/reset` 设置新密码。
- 验证码发送受冷却与频控限制（`Auth.Verification`/`Auth.PasswordReset`），建议前端增加倒计时提示。

## 5. 页面与接口映射

### 管理端

- 仪表盘：`GET /api/v1/{adminPrefix}/dashboard`
- 用户管理：`GET/POST/PATCH /users`、`POST /users/{id}/reset-password`、`POST /users/{id}/force-logout`
- 节点管理：`GET/POST/PATCH /nodes`、`POST /nodes/{id}/disable`、`GET /nodes/{id}/kernels`、`POST /nodes/{id}/kernels/sync`
- 订阅模板：`GET/POST/PATCH /subscription-templates`、`POST /subscription-templates/{id}/publish`
- 订阅管理：`GET/POST/PATCH /subscriptions`、`POST /subscriptions/{id}/disable`、`POST /subscriptions/{id}/extend`
- 套餐管理：`GET/POST/PATCH /plans`
- 优惠券管理：`GET/POST/PATCH/DELETE /coupons`
- 公告管理：`GET/POST /announcements`、`POST /announcements/{id}/publish`
- 安全配置：`GET/PATCH /security-settings`
- 审计日志：`GET /audit-logs`、`GET /audit-logs/export`
- 订单管理：`GET /orders`、`GET /orders/{id}`、`POST /orders/{id}/pay|cancel|refund`

### 用户端

- 订阅列表/预览/切换模板：`GET /subscriptions`、`GET /subscriptions/{id}/preview`、`POST /subscriptions/{id}/template`
- 套餐列表：`GET /plans`
- 节点状态：`GET /nodes`
- 公告列表：`GET /announcements`
- 账户资料：`GET/PATCH /account/profile`
- 自助改密：`POST /account/password`
- 自助改邮箱：`POST /account/email/code`、`POST /account/email`
- 余额与流水：`GET /account/balance`
- 订单：`POST /orders`、`GET /orders`、`GET /orders/{id}`、`GET /orders/{id}/payment-status`、`POST /orders/{id}/cancel`

完整字段说明请参考 `docs/api-reference.md`，或使用 `./scripts/gen-api-docs.sh` 生成的 `docs/api-generated/`。

## 6. 同步与订阅交互定义

### 6.1 管理端节点同步

**触发同步**

- `POST /api/v1/{adminPrefix}/nodes/{id}/kernels/sync`
- 请求体（可选）：
  - `protocol` string（空表示默认协议）

示例：

```http
POST /api/v1/admin/nodes/42/kernels/sync
Content-Type: application/json

{"protocol":"http"}
```

响应字段：

- `node_id` uint64
- `protocol` string
- `revision` string
- `synced_at` int64（Unix 秒）
- `message` string

**错误场景**

- `400`：协议不支持或参数非法
- `404`：节点不存在
- `500`：内核同步失败（检查 Kernel 地址/令牌）

### 6.2 管理端协议绑定下发

**单条下发**

- `POST /api/v1/{adminPrefix}/protocol-bindings/{id}/sync`

响应：

```json
{"binding_id":3,"status":"synced","message":"ok","synced_at":1719766500}
```

**批量下发**

- `POST /api/v1/{adminPrefix}/protocol-bindings/sync`
- 请求体：
  - `binding_ids` []uint64（可选）
  - `node_ids` []uint64（可选）

响应：

```json
{
  "results": [
    {"binding_id":3,"status":"synced","message":"ok","synced_at":1719766500},
    {"binding_id":4,"status":"error","message":"kernel control not configured","synced_at":1719766500}
  ]
}
```

**状态约定**

- `status=synced` 表示下发成功
- `status=error` 表示失败（`message` 描述原因）

### 6.3 用户侧节点状态

- `GET /api/v1/user/nodes`
- 查询参数：`page`、`per_page`、`status`、`protocol`

关键字段说明：

- `nodes[].status`：节点状态（`online`/`offline`/`maintenance`/`disabled`）
- `kernel_statuses[]`：节点同步摘要（来自最近一次同步记录）
- `protocol_statuses[]`：协议绑定健康状态
  - `health_status`：`healthy`/`degraded`/`unhealthy`/`offline`/`unknown`
- 节点范围来自当前生效订阅（`active` 且未过期）绑定的协议；无生效订阅时返回空数组
- `protocol` 过滤时仅在套餐允许的协议绑定中筛选

提示：`kernel_statuses` 表示同步记录状态，不是实时心跳。

### 6.4 订阅预览与鉴权字段

- `GET /api/v1/user/subscriptions/{id}/preview`
- 查询参数：`template_id`（可选）
  - 用户侧订阅列表默认不返回 `disabled` 状态，`expired` 仍可展示用于续费

响应字段：

- `content`：渲染后的订阅内容
- `content_type`：`text/plain` 或 `application/json`
- `etag`：内容哈希
- `generated_at`：生成时间

模板变量中的鉴权字段：

- `user_identity.account_id` / `user_identity.password`
- `user_identity.account` / `user_identity.id` / `user_identity.uuid`

当订阅 `status != active` 时：

- `nodes`/`protocol_bindings` 输出为空数组
- `user_identity` 字段为空字符串
  - `status=disabled` 时接口返回 `404`

当订阅 `status = active` 时：

- `nodes`/`protocol_bindings` 仅包含套餐绑定的协议

模板变量中的套餐快照字段：

- `subscription.plan_snapshot`：套餐快照（含 `binding_ids`、`traffic_multipliers` 等），用于保证订阅长期一致性

### 6.5 订阅拉取地址（客户端）

- `GET /api/v1/subscriptions/{token}`
- 根据 `User-Agent` 自动选择模板：
  - 命中 `clash` 相关客户端 → `client_type=clash`
  - 命中 `sing-box` 客户端 → `client_type=sing-box`
  - 常见识别关键字：`mihomo`、`clash-verge`、`surge`、`quantumult`、`stash`、`shadowrocket`、`loon`、`nekobox`、`v2rayn`、`v2rayng`
- 未识别则回退订阅默认模板
- 订阅非 `active` 或已过期时返回 `404`

## 7. 数据格式与展示建议

- **金额**：`*_cents` 为分单位，展示时建议 `amount_cents / 100` 并配合 `currency`。
- **流量**：`traffic_limit_bytes`、`traffic_used_bytes` 建议使用二进制单位（GB/TB）。
- **时间**：所有 `*_at` 字段为 Unix 秒（UTC），前端需本地化显示。
- **订单状态**：
  - `status`：`pending_payment`、`paid`、`payment_failed`、`cancelled`、`partially_refunded`、`refunded`
  - `payment_status`：`pending`、`succeeded`、`failed`
- **套餐状态**：`draft`、`active`（未激活套餐前端可隐藏）

## 8. 订单与支付流程提示

- `POST /user/orders` 支持 `payment_method=balance|external|manual`（manual 表示线下/人工支付）。
- `payment_method=external` 且金额大于 0 时，需要传 `payment_channel`，响应会带 `payment_intent_id` 与 `payments`，其中 `payments[].metadata.pay_url`/`qr_code` 用于跳转或展示二维码。
- `payment_method=manual` 会创建待支付订单，需管理员通过 `/api/v1/{adminPrefix}/orders/{id}/pay` 标记已支付。
- 推荐前端传 `idempotency_key`（如点击下单时生成 UUID），避免重复下单。
- 传入 `coupon_code` 命中优惠后，订单会追加 `item_type=discount` 条目并在 `order.metadata` 中返回折扣明细。

### 8.1 优惠券交互提示

- 用户在下单页输入优惠券码，前端透传 `coupon_code`（大小写不敏感，建议 trim）。
- 服务端校验失败会返回 `400`，需要展示原因（未启用/过期/次数用尽/未达最低金额）。
- 命中优惠后：
  - `order.metadata` 包含 `coupon_code`、`coupon_id`、`discount_cents`。
  - `order.items` 追加 `item_type=discount` 条目，`subtotal_cents` 为负值，可用于订单明细展示。
- 优惠券不影响订阅身份体系，仅影响订单应付金额。

## 9. 第三方签名开关

用户端路由统一挂载第三方签名中间件：

- 若 `security_settings.third_party_api_enabled=true` 且 `api_key/api_secret` 生效，前端必须附带签名头。
- 浏览器端不适合存储 `api_secret`，建议在后台关闭该开关或通过 BFF 服务代签名。

## 10. 管理端访问限制

管理端可能开启 IP 白名单与速率限制（`Admin.Access`）：

- 前端部署地址需在允许网段内。
- 被限流时返回 `429`，可做提示与退避重试。

## 11. 本地联调建议

1. 启动后端：

```bash
go run ./cmd/znp migrate --config etc/znp-sqlite.yaml --apply --seed-demo
go run ./cmd/znp serve --config etc/znp-sqlite.yaml --migrate-to latest
```

2. 默认账号：

- 管理员：`admin@example.com` / `P@ssw0rd!`
- 用户：`user@example.com` / `P@ssw0rd!`

3. API Base：`http://localhost:8888/api/v1`

## 12. 常见问题排查

- `401`：检查 token 过期、角色是否匹配。
- `403`：检查角色、IP 白名单或第三方签名是否启用。
- `409`：通常是并发更新/状态冲突，前端应提示重试。
- `422/400`：检查必填字段与格式是否匹配。
