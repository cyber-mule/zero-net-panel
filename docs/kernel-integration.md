# 内核对接文档

本项目的网络内核对外接口以仓库根目录的 `core.yaml` 为准（OpenAPI 3.0.3）。
该文件覆盖 `/v1/*` 控制面接口，并对 SSE/WebSocket 流式接口提供摘要说明，
可视为允许外部访问的 API 合同/清单。

## 使用方式

- 将 `core.yaml` 导入 Postman/Insomnia 或其他 OpenAPI 工具进行调试与对接。
- 将 `servers` 中的 `url` 替换为实际内核 HTTP 地址（应与节点 `control_endpoint` 对齐）。
- 默认使用 Basic 鉴权；若内核配置 `api.auth.allow_insecure=true`，可不携带 Authorization。

面板侧鉴权配置说明：

- 节点控制面必须配置 `control_endpoint`，面板不再回退全局 `Kernel.HTTP.BaseURL`。
- 在节点上配置 `control_access_key` + `control_secret_key` 时，面板会自动发送 Basic。
- 鉴权优先级：`control_access_key` + `control_secret_key` → `control_token`（无全局兜底）。

## 与面板对接关系

- 面板通过 `POST /api/v1/{admin}/nodes/{id}/kernels/sync` 触发节点与内核同步。
- 内核侧应实现 `core.yaml` 中的控制面接口（如 `/v1/status`、`/v1/traffic` 等），
  具体字段与错误码以 `core.yaml` 为准。
- `GET /v1/status` 返回的节点 `id` 为字符串；面板侧需将协议绑定的 `kernel_id` 与该 `id` 对齐。

## 运行状态检查

内核提供状态接口用于确认服务是否运行：

- `GET /healthz`：轻量存活探针，返回 `ok` 即表示服务可用。
- `GET /v1/status`：返回运行时快照，可用于判断节点/协议状态是否健康。

## 节点注册与心跳

当内核支持动态节点接入（多协议/多实例）时，可使用注册与心跳接口维护节点存活状态。

交互流程：

1. 节点启动后调用 `POST /v1/protocols/registrations` 注册自身信息（`id`/`role`/`protocol` 等）。
2. 内核返回 `expires_at_ms` 与 `heartbeat_interval_seconds`，节点按间隔发心跳。
3. 节点定期调用 `POST /v1/protocols/registrations/{id}/heartbeat` 上报健康状态。
4. 节点下线前调用 `DELETE /v1/protocols/registrations/{id}` 注销注册。

注册结果会反映在 `/v1/status` 的 `nodes` 列表中，必要时可结合 `node_*` 事件回调订阅状态变更。

## 注册通知回调

内核提供事件回调注册接口，用于订阅节点状态与服务级事件：

- 节点事件：`POST /v1/events/registrations`（如 `node_added`、`node_healthy`、`node_degraded` 等）。
- 服务事件：`POST /v1/service-events/registrations`（如 `user.quota.changed`、`user.traffic.reported`）。

内核会在事件发生时向 callback 地址推送通知，已注册的回调可通过
`GET /v1/events/registrations` 与 `GET /v1/service-events/registrations` 查询，删除使用对应的
`DELETE` 接口（详见 `core.yaml`）。

## 面板回调接入

面板侧新增以下回调入口以承接内核推送（受 `Webhook` 配置保护，默认使用 `X-ZNP-Webhook-Token`）：

- `POST /api/v1/kernel/events`：节点状态事件回调。
- `POST /api/v1/kernel/traffic`：用户流量观测回调。
- `POST /api/v1/kernel/service-events`：服务事件回调（如 `user.traffic.reported`）。

当前节点事件回调仅用于记录，不会自动更新协议健康状态；需要时请手动触发协议健康反向同步。

节点事件示例（`id` 或 `node_id` 至少一个）：

```json
{
  "event": "node_healthy",
  "id": "edge-hk-1-vless",
  "status": "healthy",
  "message": "ok",
  "observed_at": 1734001010
}
```

流量回调示例（批量）：

```json
{
  "records": [
    {
      "user_id": 1,
      "subscription_id": 2,
      "protocol": "vless",
      "node_id": 10,
      "binding_id": 3,
      "bytes_up": 1234,
      "bytes_down": 5678,
      "observed_at": 1734001010
    }
  ]
}
```

服务事件回调示例（`user.traffic.reported`）：

```json
{
  "event": "user.traffic.reported",
  "payload": {
    "user_id": "alice",
    "previous": {"used": 1024, "remaining": 2048},
    "current": {"used": 2048, "remaining": 1024}
  }
}
```

面板侧会优先使用 `subscription_id`，否则使用 `user_id`（需与面板用户 ID 对齐）来更新订阅已用流量（`current.used`）。

当内核无法推送事件时，可通过 `Kernel.StatusPollInterval` 启用状态轮询，面板会定期调用
`GET /v1/status` 判断节点控制面可达性，并将节点 `status` 更新为 `online`/`offline`。

如需即时刷新某些节点的在线状态，可调用管理端：

- `POST /api/v1/{admin}/nodes/status/sync`（请求体传 `node_ids`）

协议健康度不会自动反向同步，需手动触发：

- `POST /api/v1/{admin}/protocol-bindings/status/sync`（请求体传 `node_ids`）

为了避免内核不可用时刷屏，支持失败退避：

```yaml
Kernel:
  StatusPollInterval: 30s
  StatusPollBackoff:
    Enabled: true
    MaxInterval: 5m
    Multiplier: 2
    Jitter: 0.2
```

- `Enabled=true` 时，连续失败会按倍率退避，成功后恢复基准间隔
- `MaxInterval` 为退避上限
- `Multiplier` 为退避倍率
- `Jitter` 为抖动比例（0~1）

为缩短离线节点恢复感知时间，面板会对离线节点执行补偿轮询：

- 轮询间隔按 1/3/5 秒递增，跨自然日重置回 1s
- 可通过 `Kernel.OfflineProbeMaxInterval` 或站点配置 `kernel_offline_probe_max_interval_seconds` 设置最大间隔（秒），默认 0 不限制

## 维护约定

- 内核对外接口变更时需同步更新 `core.yaml`，并在 PR 中注明影响范围与兼容策略。
