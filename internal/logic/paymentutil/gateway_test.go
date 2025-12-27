package paymentutil

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

func TestInitiateHTTP(t *testing.T) {
	var payload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"data":{"pay_url":"https://pay.test/redirect","reference":"ref-123"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	channel := repository.PaymentChannel{
		Code:     "stripe",
		Provider: "stripe",
		Config: map[string]any{
			"mode":       "http",
			"notify_url": "https://notify.test/callback?order_id={{order_id}}&payment_id={{payment_id}}",
			"return_url": "https://return.test/orders/{{order_number}}",
			"http": map[string]any{
				"endpoint":  server.URL,
				"method":    "POST",
				"body_type": "json",
				"payload": map[string]any{
					"order_no":   "{{order_number}}",
					"amount":     "{{amount}}",
					"notify_url": "{{notify_url}}",
					"return_url": "{{return_url}}",
				},
			},
			"response": map[string]any{
				"pay_url":   "data.pay_url",
				"reference": "data.reference",
			},
		},
	}

	order := repository.Order{
		ID:              10,
		Number:          "ORD-10",
		UserID:          5,
		TotalCents:      1234,
		Currency:        "CNY",
		PaymentIntentID: "intent-10",
	}
	payment := repository.OrderPayment{
		ID: 22,
	}

	result, err := Initiate(context.Background(), InitiateParams{
		Channel:  channel,
		Order:    order,
		Payment:  payment,
		Quantity: 1,
		PlanID:   3,
		PlanName: "Basic",
	})
	require.NoError(t, err)
	require.Equal(t, "https://pay.test/redirect", result.Metadata["pay_url"])
	require.Equal(t, "ref-123", result.Reference)

	require.Equal(t, "ORD-10", payload["order_no"])
	require.Equal(t, "12.34", payload["amount"])
	require.Contains(t, payload["notify_url"].(string), "order_id=10")
	require.Contains(t, payload["notify_url"].(string), "payment_id=22")
	require.Equal(t, "https://return.test/orders/ORD-10", payload["return_url"])
}

func TestRefundHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"data":{"reference":"refund-001","status":"success"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	channel := repository.PaymentChannel{
		Code:     "stripe",
		Provider: "stripe",
		Config: map[string]any{
			"refund": map[string]any{
				"http": map[string]any{
					"endpoint":  server.URL,
					"method":    "POST",
					"body_type": "json",
					"payload": map[string]any{
						"amount": "{{refund_amount}}",
					},
				},
				"response": map[string]any{
					"reference": "data.reference",
					"status":    "data.status",
				},
				"status_map": map[string]any{
					"success": "succeeded",
				},
			},
		},
	}

	result, err := Refund(context.Background(), RefundParams{
		Channel:           channel,
		Order:             repository.Order{ID: 1, Number: "ORD-1", TotalCents: 990, Currency: "CNY"},
		Payment:           repository.OrderPayment{ID: 2},
		RefundAmountCents: 990,
		RefundReason:      "test",
	})
	require.NoError(t, err)
	require.Equal(t, "refund-001", result.Reference)
	require.Equal(t, "refund-001", result.Metadata["gateway_refund_reference"])
}

func TestReconcileHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status":"paid","reference":"pay-001"}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	channel := repository.PaymentChannel{
		Code:     "stripe",
		Provider: "stripe",
		Config: map[string]any{
			"reconcile": map[string]any{
				"http": map[string]any{
					"endpoint":  server.URL,
					"method":    "POST",
					"body_type": "json",
					"payload": map[string]any{
						"payment_id": "{{payment_id}}",
					},
				},
				"response": map[string]any{
					"status":    "status",
					"reference": "reference",
				},
				"status_map": map[string]any{
					"paid": "succeeded",
				},
			},
		},
	}

	result, err := Reconcile(context.Background(), ReconcileParams{
		Channel: channel,
		Order:   repository.Order{ID: 2, Number: "ORD-2", TotalCents: 1200, Currency: "CNY"},
		Payment: repository.OrderPayment{ID: 8},
	})
	require.NoError(t, err)
	require.Equal(t, repository.OrderPaymentStatusSucceeded, result.Status)
	require.Equal(t, "pay-001", result.Reference)
}

func TestVerifyWebhookSignature(t *testing.T) {
	channel := repository.PaymentChannel{
		Code: "stripe",
		Config: map[string]any{
			"webhook": map[string]any{
				"signature_type":   "hmac_sha256",
				"signature_header": "X-Pay-Signature",
				"secret":           "secret",
			},
		},
	}

	body := []byte(`{"payment_id":1}`)
	mac := hmac.New(sha256.New, []byte("secret"))
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))

	headers := http.Header{}
	headers.Set("X-Pay-Signature", signature)

	err := VerifyWebhookSignature(channel, body, headers)
	require.NoError(t, err)

	headers.Set("X-Pay-Signature", "invalid")
	err = VerifyWebhookSignature(channel, body, headers)
	require.ErrorIs(t, err, repository.ErrUnauthorized)
}
