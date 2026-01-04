# 协议与节点设计总览（面板侧）

## 背景与目标
- 节点与协议解耦：节点仅描述机房/线路等资源，协议绑定描述内核实际配置，协议发布描述对外入口。
- 协议绑定（节点 + 实际协议配置）作为内核下发的最小单元，支持中转/非入口节点等场景。
- 订阅输出基于协议发布组装，支持模板化订阅内容。
- 下发协议为手动触发能力；节点注册仅建立面板与内核的 API 通讯关系。

## 数据模型
- 节点 `nodes`：节点元信息、状态、标签等（`internal/repository/node_repository.go`）。
- 协议绑定 `protocol_bindings`：绑定节点 + 实际协议配置（profile）+ 监听/中转信息 + 同步/健康状态（`internal/repository/protocol_binding_repository.go`）。
- 协议发布 `protocol_entries`：对外入口（域名/IP + 端口 + 公共 profile），仅控制用户可见状态，健康状态共享自绑定（`internal/repository/protocol_entry_repository.go`）。
- 流量记录 `traffic_usage_records`：原始与倍数计费流量（`internal/repository/traffic_usage_repository.go`）。
- 套餐倍数 `plans.traffic_multipliers`：协议 -> 倍数（`internal/repository/plan_repository.go`）。

## 分层规则
- 协议绑定仅用于内核下发；协议发布仅用于订阅交付与用户可见入口。
- 协议发布 `entry_address/entry_port` 为对外入口，可与绑定 `listen/access_port` 不一致，用于中转或分流场景。
- 协议发布 `status` 为虚拟状态，仅影响用户可见；绑定 `health_status` 代表真实健康并在发布列表共享展示。
- 绑定 `listen` 为空或仅端口时，会用 `access_port` 归一化为 `0.0.0.0:<port>` 供内核使用。

## 管理能力
- 协议绑定管理（增删改查）：`/api/v1/{admin}/protocol-bindings`
- 协议发布管理（增删改查）：`/api/v1/{admin}/protocol-entries`
- 协议下发（手动同步）：`/api/v1/{admin}/protocol-bindings/:id/sync` 与批量 `/sync`
- 节点内核端点维护与同步（原有能力保留）：`/api/v1/{admin}/nodes/:id/kernels`

## 内核交互
- 协议下发：面板调用内核控制面 `POST /v1/protocols`（`pkg/kernel/control_client.go`）。
- 注册/心跳：内核侧接口定义在 `core.yaml`，面板侧对接说明见 `docs/kernel-integration.md`。
- 回调接入：
  - `POST /api/v1/kernel/events`：节点健康事件回调（更新协议绑定健康）。
  - `POST /api/v1/kernel/traffic`：用户流量观测回调（记录原始/倍数流量）。
- 状态轮询：`Kernel.StatusPollInterval` 触发 `GET /v1/status` 轮询（`internal/logic/kernel/statussync.go`）。

## 订阅与展示
- 订阅渲染基于协议发布上下文输出 `entry_address/entry_port` 与公开 profile（`internal/logic/user/subscription/previewlogic.go`）。
- 用户侧节点状态：按协议发布可见状态 + 协议绑定健康过滤后展示，响应仍以绑定健康与内核同步摘要为主（`/api/v1/user/nodes`）。
- 用户侧流量查询：`/api/v1/user/subscriptions/:id/traffic`，返回原始/倍数流量与倍数系数。

## 流量计费与倍数
- 原始流量：`bytes_up + bytes_down`
- 计费流量：`raw_bytes * multiplier`，倍数取自套餐 `traffic_multipliers[protocol]`，默认 1。
- 结算以倍数后的消耗为准，同时保留原始与倍数消耗用于查询/审计。

## 用户身份与轮换
- 用户维度唯一身份（账户/密码）独立于订阅，仅在订阅生效时参与下发与订阅渲染。
- 订阅模板使用 `user_identity` 渲染鉴权字段（如 `user_identity.account_id`/`user_identity.password`），保留 `subscription.token` 作为兼容字段但不推荐用于鉴权。
- 轮换仅允许手动触发：用户自助接口与管理端用户操作接口均可发起。
- 身份信息加密存储并保留指纹与时间轴，用于延迟上报审计追溯。

## 实现进度（对照清单）
- [x] 协议发布/协议绑定/流量记录数据模型 + 迁移
- [x] 管理端协议发布/协议绑定 CRUD + 手动下发
- [x] 内核控制面下发协议能力
- [x] 节点事件回调接入 + 状态轮询能力
- [x] 流量回调接入 + 倍数核算 + 用户查询
- [x] 订阅渲染与用户侧节点状态展示调整
- [x] 用户身份（账户/密码）加密存储 + 手动轮换 + 订阅渲染接入

## 相关文件索引
- 数据模型与迁移：`internal/repository/*.go`，`internal/bootstrap/migrations/registry.go`
- 管理端 API：`api/admin/protocol_entries.api`，`api/admin/protocol_bindings.api`
- 用户端 API：`api/user/nodes.api`，`api/user/subscriptions.api`
- 内核对接文档：`docs/kernel-integration.md`，`core.yaml`
- 配置示例：`etc/znp-api.yaml`，`etc/znp-sqlite.yaml`，`etc/znp-prod.example.yaml`
- 用户身份文档：`docs/protocol-credentials.md`
