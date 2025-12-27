package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type gatewayConfig struct {
	signatureSecret string
	signatureHeader string
	webhookToken    string
	callbackDelay   time.Duration
}

func main() {
	var (
		addr            string
		signatureSecret string
		signatureHeader string
		webhookToken    string
		callbackDelay   time.Duration
	)
	flag.StringVar(&addr, "addr", ":9099", "listen address")
	flag.StringVar(&signatureSecret, "signature-secret", "", "HMAC secret for callback signature")
	flag.StringVar(&signatureHeader, "signature-header", "X-Pay-Signature", "callback signature header name")
	flag.StringVar(&webhookToken, "webhook-token", "", "X-ZNP-Webhook-Token header")
	flag.DurationVar(&callbackDelay, "callback-delay", time.Second, "delay before sending callback")
	flag.Parse()

	cfg := gatewayConfig{
		signatureSecret: signatureSecret,
		signatureHeader: signatureHeader,
		webhookToken:    webhookToken,
		callbackDelay:   callbackDelay,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pay", cfg.handlePay)
	mux.HandleFunc("/refund", cfg.handleRefund)
	mux.HandleFunc("/query", cfg.handleQuery)

	log.Printf("mock gateway listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func (cfg gatewayConfig) handlePay(w http.ResponseWriter, r *http.Request) {
	payload := readPayload(r)
	notifyURL := getString(payload, "notify_url")
	orderID := getString(payload, "order_id")
	paymentID := getString(payload, "payment_id")
	reference := "pay-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	response := map[string]any{
		"data": map[string]any{
			"pay_url":   "https://mock.pay/redirect/" + reference,
			"reference": reference,
			"status":    "created",
		},
	}
	writeJSON(w, response)

	if notifyURL != "" && orderID != "" && paymentID != "" {
		go cfg.sendCallback(notifyURL, orderID, paymentID, reference)
	}
}

func (cfg gatewayConfig) handleRefund(w http.ResponseWriter, r *http.Request) {
	reference := "refund-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	response := map[string]any{
		"data": map[string]any{
			"reference": reference,
			"status":    "success",
		},
	}
	writeJSON(w, response)
}

func (cfg gatewayConfig) handleQuery(w http.ResponseWriter, r *http.Request) {
	payload := readPayload(r)
	reference := getString(payload, "payment_ref")
	if reference == "" {
		reference = getString(payload, "payment_id")
	}
	if reference == "" {
		reference = "unknown"
	}
	response := map[string]any{
		"data": map[string]any{
			"status":    "paid",
			"reference": reference,
		},
	}
	writeJSON(w, response)
}

func (cfg gatewayConfig) sendCallback(notifyURL, orderID, paymentID, reference string) {
	time.Sleep(cfg.callbackDelay)

	body := map[string]any{
		"order_id":   toUint(orderID),
		"payment_id": toUint(paymentID),
		"status":     "succeeded",
		"reference":  reference,
		"paid_at":    time.Now().UTC().Unix(),
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, notifyURL, bytes.NewReader(data))
	if err != nil {
		log.Printf("callback request error: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.webhookToken != "" {
		req.Header.Set("X-ZNP-Webhook-Token", cfg.webhookToken)
	}
	if cfg.signatureSecret != "" {
		signature := signBody(cfg.signatureSecret, data)
		req.Header.Set(cfg.signatureHeader, signature)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("callback error: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("callback status %d", resp.StatusCode)
}

func signBody(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func readPayload(r *http.Request) map[string]any {
	payload := map[string]any{}
	if r.Body == nil {
		return payload
	}
	defer r.Body.Close()

	ct := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if strings.Contains(ct, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return payload
		}
		for key, values := range r.PostForm {
			if len(values) > 0 {
				payload[key] = values[0]
			}
		}
		return payload
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return payload
	}
	return payload
}

func writeJSON(w http.ResponseWriter, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func getString(payload map[string]any, key string) string {
	value, ok := payload[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	default:
		return strings.TrimSpace(toString(value))
	}
}

func toUint(value string) uint64 {
	parsed, _ := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	return parsed
}

func toString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []string:
		if len(typed) > 0 {
			return typed[0]
		}
	case []any:
		if len(typed) > 0 {
			if s, ok := typed[0].(string); ok {
				return s
			}
		}
	case url.Values:
		for _, values := range typed {
			if len(values) > 0 {
				return values[0]
			}
		}
	}
	return ""
}
