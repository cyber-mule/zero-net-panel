# 外部支付联调示例（发起/回调/退款/对账）

本示例使用内置的 mock 网关脚本完成全链路联调，适合本地调试与接口验证。

## 1) 启动 mock 网关

```bash
go run ./scripts/mock-payment-gateway.go -addr :9099 -signature-secret mock-secret
```

## 2) 创建支付通道（管理端）

```bash
curl -X POST http://localhost:8888/api/v1/admin/payment-channels \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <ADMIN_TOKEN>' \
  -d '{
    "name": "Mock Gateway",
    "code": "mockpay",
    "provider": "mockpay",
    "enabled": true,
    "sort_order": 1,
    "config": {
      "mode": "http",
      "notify_url": "http://localhost:8888/api/v1/payments/callback",
      "return_url": "http://localhost:3000/orders/{{order_number}}",
      "http": {
        "endpoint": "http://localhost:9099/pay",
        "method": "POST",
        "body_type": "json",
        "payload": {
          "order_id": "{{order_id}}",
          "payment_id": "{{payment_id}}",
          "amount": "{{amount}}",
          "notify_url": "{{notify_url}}",
          "return_url": "{{return_url}}"
        }
      },
      "response": {
        "pay_url": "data.pay_url",
        "reference": "data.reference"
      },
      "webhook": {
        "signature_type": "hmac_sha256",
        "signature_header": "X-Pay-Signature",
        "secret": "mock-secret"
      },
      "refund": {
        "http": {
          "endpoint": "http://localhost:9099/refund",
          "method": "POST",
          "body_type": "json",
          "payload": {
            "payment_ref": "{{payment_reference}}",
            "amount": "{{refund_amount}}",
            "reason": "{{refund_reason}}"
          }
        },
        "response": {
          "reference": "data.reference",
          "status": "data.status"
        },
        "status_map": {
          "success": "succeeded",
          "failed": "failed"
        }
      },
      "reconcile": {
        "http": {
          "endpoint": "http://localhost:9099/query",
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
  }'
```

提示：
- 若启用了全局 Webhook 校验（`Webhook.SharedToken` 或 `Webhook.Stripe.SigningSecret`），请按配置补齐对应请求头，或在本地关闭该校验。

## 3) 下单并获取支付链接

```bash
curl -X POST http://localhost:8888/api/v1/user/orders \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <USER_TOKEN>' \
  -d '{
    "plan_id": 1,
    "quantity": 1,
    "payment_method": "external",
    "payment_channel": "mockpay"
  }'
```

响应中的 `payments[0].metadata.pay_url` 即为跳转地址。

## 4) 退款（外部支付）

```bash
curl -X POST http://localhost:8888/api/v1/admin/orders/<ORDER_ID>/refund \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <ADMIN_TOKEN>' \
  -d '{
    "amount_cents": 1500,
    "reason": "test refund"
  }'
```

## 5) 对账（外部支付）

```bash
curl -X POST http://localhost:8888/api/v1/admin/orders/payments/reconcile \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer <ADMIN_TOKEN>' \
  -d '{
    "order_id": <ORDER_ID>,
    "payment_id": <PAYMENT_ID>
  }'
```

## 6) 观察结果

- `GET /api/v1/user/orders/{id}` / `GET /api/v1/{adminPrefix}/orders/{id}` 查看 `payment_status`、`payments[].metadata` 与退款记录。
