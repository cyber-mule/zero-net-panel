package paymentutil

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

const defaultGatewayTimeout = 10 * time.Second

// ChannelConfig describes how to initiate external payments.
type ChannelConfig struct {
	Mode      string          `json:"mode"`
	NotifyURL string          `json:"notify_url"`
	ReturnURL string          `json:"return_url"`
	HTTP      HTTPGateway     `json:"http"`
	Response  ResponseMapping `json:"response"`
	Webhook   *WebhookConfig  `json:"webhook"`
	Refund    *ActionConfig   `json:"refund"`
	Reconcile *ActionConfig   `json:"reconcile"`
}

// HTTPGateway defines request settings for external payment initiation.
type HTTPGateway struct {
	Endpoint       string            `json:"endpoint"`
	Method         string            `json:"method"`
	BodyType       string            `json:"body_type"`
	Headers        map[string]string `json:"headers"`
	Payload        map[string]any    `json:"payload"`
	TimeoutSeconds int               `json:"timeout_seconds"`
}

// ResponseMapping declares how to extract fields from the gateway response.
type ResponseMapping struct {
	PayURL    string `json:"pay_url"`
	QRCode    string `json:"qr_code"`
	IntentID  string `json:"intent_id"`
	Reference string `json:"reference"`
}

// WebhookConfig defines callback signature verification settings.
type WebhookConfig struct {
	SignatureType   string `json:"signature_type"`
	SignatureHeader string `json:"signature_header"`
	Secret          string `json:"secret"`
}

// ActionConfig defines HTTP actions like refund or reconciliation.
type ActionConfig struct {
	HTTP      HTTPGateway           `json:"http"`
	Response  ActionResponseMapping `json:"response"`
	StatusMap map[string]string     `json:"status_map"`
}

// ActionResponseMapping declares how to extract action fields.
type ActionResponseMapping struct {
	Status         string `json:"status"`
	Reference      string `json:"reference"`
	FailureCode    string `json:"failure_code"`
	FailureMessage string `json:"failure_message"`
}

// InitiateParams carries context for initiating an external payment.
type InitiateParams struct {
	Channel   repository.PaymentChannel
	Order     repository.Order
	Payment   repository.OrderPayment
	Quantity  int
	PlanID    uint64
	PlanName  string
	ReturnURL string
}

// InitiateResult summarizes the initiation response.
type InitiateResult struct {
	Metadata  map[string]any
	Reference string
}

// RefundParams carries refund context.
type RefundParams struct {
	Channel           repository.PaymentChannel
	Order             repository.Order
	Payment           repository.OrderPayment
	RefundAmountCents int64
	RefundReason      string
}

// RefundResult summarizes refund action response.
type RefundResult struct {
	Reference string
	Metadata  map[string]any
}

// ReconcileParams carries reconciliation context.
type ReconcileParams struct {
	Channel repository.PaymentChannel
	Order   repository.Order
	Payment repository.OrderPayment
}

// ReconcileResult summarizes reconciliation response.
type ReconcileResult struct {
	Status         string
	Reference      string
	FailureCode    string
	FailureMessage string
	Metadata       map[string]any
}

// Initiate triggers an external payment request based on channel configuration.
func Initiate(ctx context.Context, params InitiateParams) (InitiateResult, error) {
	cfg, err := parseChannelConfig(params.Channel.Config)
	if err != nil {
		return InitiateResult{}, err
	}
	cfg.normalize()

	switch cfg.Mode {
	case "", "http":
		return initiateHTTP(ctx, cfg, params)
	default:
		return InitiateResult{}, repository.ErrInvalidArgument
	}
}

// VerifyWebhookSignature validates gateway callback signatures when configured.
func VerifyWebhookSignature(channel repository.PaymentChannel, body []byte, headers http.Header) error {
	if len(channel.Config) == 0 {
		return nil
	}
	cfg, err := parseChannelConfig(channel.Config)
	if err != nil {
		return err
	}
	cfg.normalize()
	if cfg.Webhook == nil {
		return nil
	}
	if cfg.Webhook.Secret == "" || cfg.Webhook.SignatureHeader == "" {
		return nil
	}
	signature := strings.TrimSpace(headers.Get(cfg.Webhook.SignatureHeader))
	if signature == "" {
		return repository.ErrUnauthorized
	}
	expected, err := computeSignature(cfg.Webhook.SignatureType, cfg.Webhook.Secret, body)
	if err != nil {
		return err
	}
	if !signatureEqual(signature, expected) {
		return repository.ErrUnauthorized
	}
	return nil
}

// Refund triggers an external refund action.
func Refund(ctx context.Context, params RefundParams) (RefundResult, error) {
	if params.RefundAmountCents <= 0 {
		return RefundResult{}, repository.ErrInvalidArgument
	}
	cfg, err := parseChannelConfig(params.Channel.Config)
	if err != nil {
		return RefundResult{}, err
	}
	cfg.normalize()
	if cfg.Refund == nil || cfg.Refund.HTTP.Endpoint == "" {
		return RefundResult{}, repository.ErrInvalidArgument
	}

	vars := buildBaseVars(params.Channel, params.Order, params.Payment, 0, "", 0)
	vars["refund_amount_cents"] = strconv.FormatInt(params.RefundAmountCents, 10)
	vars["refund_amount"] = fmt.Sprintf("%.2f", float64(params.RefundAmountCents)/100.0)
	vars["refund_reason"] = strings.TrimSpace(params.RefundReason)

	endpoint := applyTemplate(cfg.Refund.HTTP.Endpoint, vars)
	payload := expandPayload(cfg.Refund.HTTP.Payload, vars)
	payloadMap, ok := payload.(map[string]any)
	if !ok || payloadMap == nil {
		payloadMap = map[string]any{}
	}

	actionResult, err := executeAction(ctx, *cfg.Refund, endpoint, payloadMap, vars)
	if err != nil {
		return RefundResult{}, err
	}
	status := normalizeActionStatus(actionResult.Status, *cfg.Refund)
	if actionResult.Status != "" && status != repository.OrderPaymentStatusSucceeded {
		return RefundResult{}, repository.ErrInvalidArgument
	}

	metadata := map[string]any{}
	if status != "" {
		metadata["gateway_refund_status"] = status
	}
	if actionResult.Reference != "" {
		metadata["gateway_refund_reference"] = actionResult.Reference
	}

	return RefundResult{
		Reference: actionResult.Reference,
		Metadata:  metadata,
	}, nil
}

// Reconcile queries payment status from an external gateway.
func Reconcile(ctx context.Context, params ReconcileParams) (ReconcileResult, error) {
	cfg, err := parseChannelConfig(params.Channel.Config)
	if err != nil {
		return ReconcileResult{}, err
	}
	cfg.normalize()
	if cfg.Reconcile == nil || cfg.Reconcile.HTTP.Endpoint == "" {
		return ReconcileResult{}, repository.ErrInvalidArgument
	}

	vars := buildBaseVars(params.Channel, params.Order, params.Payment, 0, "", 0)
	endpoint := applyTemplate(cfg.Reconcile.HTTP.Endpoint, vars)
	payload := expandPayload(cfg.Reconcile.HTTP.Payload, vars)
	payloadMap, ok := payload.(map[string]any)
	if !ok || payloadMap == nil {
		payloadMap = map[string]any{}
	}

	actionResult, err := executeAction(ctx, *cfg.Reconcile, endpoint, payloadMap, vars)
	if err != nil {
		return ReconcileResult{}, err
	}
	status := normalizeActionStatus(actionResult.Status, *cfg.Reconcile)
	if status == "" {
		return ReconcileResult{}, repository.ErrInvalidArgument
	}

	metadata := map[string]any{}
	if actionResult.Status != "" {
		metadata["gateway_status"] = actionResult.Status
	}
	if actionResult.Reference != "" {
		metadata["gateway_reference"] = actionResult.Reference
	}

	return ReconcileResult{
		Status:         status,
		Reference:      actionResult.Reference,
		FailureCode:    actionResult.FailureCode,
		FailureMessage: actionResult.FailureMessage,
		Metadata:       metadata,
	}, nil
}

func parseChannelConfig(raw map[string]any) (ChannelConfig, error) {
	if len(raw) == 0 {
		return ChannelConfig{}, repository.ErrInvalidArgument
	}
	payload, err := json.Marshal(raw)
	if err != nil {
		return ChannelConfig{}, err
	}
	var cfg ChannelConfig
	if err := json.Unmarshal(payload, &cfg); err != nil {
		return ChannelConfig{}, err
	}
	return cfg, nil
}

func (c *ChannelConfig) normalize() {
	c.Mode = strings.ToLower(strings.TrimSpace(c.Mode))
	c.NotifyURL = strings.TrimSpace(c.NotifyURL)
	c.ReturnURL = strings.TrimSpace(c.ReturnURL)

	normalizeHTTPGateway(&c.HTTP)
	if c.Webhook != nil {
		c.Webhook.normalize()
	}
	if c.Refund != nil {
		normalizeActionConfig(c.Refund)
	}
	if c.Reconcile != nil {
		normalizeActionConfig(c.Reconcile)
	}
}

func normalizeHTTPGateway(cfg *HTTPGateway) {
	cfg.Endpoint = strings.TrimSpace(cfg.Endpoint)
	cfg.Method = strings.ToUpper(strings.TrimSpace(cfg.Method))
	if cfg.Method == "" {
		cfg.Method = http.MethodPost
	}
	cfg.BodyType = strings.ToLower(strings.TrimSpace(cfg.BodyType))
	if cfg.BodyType == "" {
		cfg.BodyType = "json"
	}
	if cfg.BodyType != "json" && cfg.BodyType != "form" {
		cfg.BodyType = "json"
	}
	if cfg.Headers == nil {
		cfg.Headers = map[string]string{}
	}
}

func normalizeActionConfig(cfg *ActionConfig) {
	normalizeHTTPGateway(&cfg.HTTP)
	if cfg.StatusMap == nil {
		return
	}
	normalized := make(map[string]string, len(cfg.StatusMap))
	for key, value := range cfg.StatusMap {
		key = strings.ToLower(strings.TrimSpace(key))
		if key == "" {
			continue
		}
		normalized[key] = strings.ToLower(strings.TrimSpace(value))
	}
	cfg.StatusMap = normalized
}

func (w *WebhookConfig) normalize() {
	w.SignatureType = strings.ToLower(strings.TrimSpace(w.SignatureType))
	w.SignatureHeader = strings.TrimSpace(w.SignatureHeader)
	w.Secret = strings.TrimSpace(w.Secret)
	if w.SignatureType == "" {
		w.SignatureType = "hmac_sha256"
	}
}

func initiateHTTP(ctx context.Context, cfg ChannelConfig, params InitiateParams) (InitiateResult, error) {
	if cfg.HTTP.Endpoint == "" {
		return InitiateResult{}, repository.ErrInvalidArgument
	}

	vars := buildBaseVars(params.Channel, params.Order, params.Payment, params.PlanID, params.PlanName, params.Quantity)
	returnURL := strings.TrimSpace(params.ReturnURL)
	if returnURL == "" {
		returnURL = cfg.ReturnURL
	}
	notifyURL := cfg.NotifyURL

	if returnURL != "" {
		returnURL = applyTemplate(returnURL, vars)
	}
	if notifyURL != "" {
		notifyURL = applyTemplate(notifyURL, vars)
	}
	vars["return_url"] = returnURL
	vars["notify_url"] = notifyURL

	endpoint := applyTemplate(cfg.HTTP.Endpoint, vars)
	payload := expandPayload(cfg.HTTP.Payload, vars)
	payloadMap, ok := payload.(map[string]any)
	if !ok || payloadMap == nil {
		payloadMap = map[string]any{}
	}

	req, err := buildGatewayRequest(ctx, cfg.HTTP, endpoint, payloadMap, vars)
	if err != nil {
		return InitiateResult{}, err
	}

	client := &http.Client{Timeout: gatewayTimeout(cfg.HTTP)}

	resp, err := client.Do(req)
	if err != nil {
		return InitiateResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return InitiateResult{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return InitiateResult{}, fmt.Errorf("payment gateway status %d", resp.StatusCode)
	}

	rawBody := strings.TrimSpace(string(body))
	if cfg.Response.PayURL == "$" {
		if rawBody == "" {
			return InitiateResult{}, errors.New("payment gateway response empty")
		}
		metadata := buildGatewayMetadata(rawBody, "", "", notifyURL, returnURL)
		return InitiateResult{
			Metadata: metadata,
		}, nil
	}

	var decoded any
	if rawBody != "" {
		if err := json.Unmarshal(body, &decoded); err != nil {
			return InitiateResult{}, err
		}
	} else {
		decoded = map[string]any{}
	}

	payURL := extractByPath(decoded, cfg.Response.PayURL)
	qrCode := extractByPath(decoded, cfg.Response.QRCode)
	intentID := extractByPath(decoded, cfg.Response.IntentID)
	reference := extractByPath(decoded, cfg.Response.Reference)

	if payURL == "" && qrCode == "" {
		return InitiateResult{}, errors.New("payment gateway response missing pay url")
	}

	metadata := buildGatewayMetadata(payURL, qrCode, intentID, notifyURL, returnURL)

	return InitiateResult{
		Metadata:  metadata,
		Reference: reference,
	}, nil
}

func buildGatewayRequest(ctx context.Context, cfg HTTPGateway, endpoint string, payload map[string]any, vars map[string]string) (*http.Request, error) {
	method := cfg.Method
	if method == "" {
		method = http.MethodPost
	}

	var (
		bodyReader  io.Reader
		contentType string
	)

	requestURL := endpoint
	if method == http.MethodGet {
		values := toValues(payload)
		if len(values) > 0 {
			parsed, err := url.Parse(endpoint)
			if err != nil {
				return nil, err
			}
			query := parsed.Query()
			for key, vals := range values {
				for _, v := range vals {
					query.Add(key, v)
				}
			}
			parsed.RawQuery = query.Encode()
			requestURL = parsed.String()
		}
	} else {
		switch cfg.BodyType {
		case "form":
			values := toValues(payload)
			bodyReader = strings.NewReader(values.Encode())
			contentType = "application/x-www-form-urlencoded"
		default:
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(payloadBytes)
			contentType = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return nil, err
	}

	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	for key, value := range cfg.Headers {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		req.Header.Set(key, applyTemplate(value, vars))
	}

	return req, nil
}

func buildBaseVars(channel repository.PaymentChannel, order repository.Order, payment repository.OrderPayment, planID uint64, planName string, quantity int) map[string]string {
	amount := fmt.Sprintf("%.2f", float64(order.TotalCents)/100.0)
	currency := strings.TrimSpace(order.Currency)
	if currency == "" {
		currency = "CNY"
	}

	vars := map[string]string{
		"order_id":          strconv.FormatUint(order.ID, 10),
		"order_number":      order.Number,
		"order_status":      order.Status,
		"payment_id":        strconv.FormatUint(payment.ID, 10),
		"payment_intent_id": order.PaymentIntentID,
		"payment_reference": order.PaymentReference,
		"payment_status":    order.PaymentStatus,
		"amount_cents":      strconv.FormatInt(order.TotalCents, 10),
		"amount":            amount,
		"currency":          currency,
		"user_id":           strconv.FormatUint(order.UserID, 10),
		"plan_id":           strconv.FormatUint(planID, 10),
		"plan_name":         planName,
		"quantity":          strconv.Itoa(quantity),
		"payment_channel":   channel.Code,
		"payment_provider":  channel.Provider,
	}

	return vars
}

func expandPayload(value any, vars map[string]string) any {
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, val := range typed {
			result[key] = expandPayload(val, vars)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for i, val := range typed {
			result[i] = expandPayload(val, vars)
		}
		return result
	case string:
		return applyTemplate(typed, vars)
	default:
		return value
	}
}

func applyTemplate(input string, vars map[string]string) string {
	result := input
	for key, value := range vars {
		token := "{{" + key + "}}"
		result = strings.ReplaceAll(result, token, value)
	}
	return result
}

func toValues(payload map[string]any) url.Values {
	values := url.Values{}
	if payload == nil {
		return values
	}

	keys := make([]string, 0, len(payload))
	for key := range payload {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		values.Set(key, fmt.Sprint(payload[key]))
	}
	return values
}

func extractByPath(payload any, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}

	value, ok := navigatePath(payload, path)
	if !ok || value == nil {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func navigatePath(payload any, path string) (any, bool) {
	parts := strings.Split(path, ".")
	current := payload
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		switch typed := current.(type) {
		case map[string]any:
			next, ok := typed[part]
			if !ok {
				return nil, false
			}
			current = next
		case []any:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(typed) {
				return nil, false
			}
			current = typed[index]
		default:
			return nil, false
		}
	}
	return current, true
}

func buildGatewayMetadata(payURL, qrCode, intentID, notifyURL, returnURL string) map[string]any {
	metadata := map[string]any{}
	if payURL != "" {
		metadata["pay_url"] = payURL
	}
	if qrCode != "" {
		metadata["qr_code"] = qrCode
	}
	if intentID != "" {
		metadata["gateway_intent_id"] = intentID
	}
	if notifyURL != "" {
		metadata["notify_url"] = notifyURL
	}
	if returnURL != "" {
		metadata["return_url"] = returnURL
	}
	return metadata
}

type actionResult struct {
	Status         string
	Reference      string
	FailureCode    string
	FailureMessage string
}

func executeAction(ctx context.Context, cfg ActionConfig, endpoint string, payload map[string]any, vars map[string]string) (actionResult, error) {
	req, err := buildGatewayRequest(ctx, cfg.HTTP, endpoint, payload, vars)
	if err != nil {
		return actionResult{}, err
	}

	timeout := gatewayTimeout(cfg.HTTP)
	client := &http.Client{Timeout: timeout}

	resp, err := client.Do(req)
	if err != nil {
		return actionResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return actionResult{}, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return actionResult{}, fmt.Errorf("payment gateway status %d", resp.StatusCode)
	}

	rawBody := strings.TrimSpace(string(body))
	var decoded any
	if rawBody != "" {
		if err := json.Unmarshal(body, &decoded); err != nil {
			return actionResult{}, err
		}
	} else {
		decoded = map[string]any{}
	}

	status := extractActionField(rawBody, decoded, cfg.Response.Status)
	reference := extractActionField(rawBody, decoded, cfg.Response.Reference)
	failureCode := extractActionField(rawBody, decoded, cfg.Response.FailureCode)
	failureMessage := extractActionField(rawBody, decoded, cfg.Response.FailureMessage)

	return actionResult{
		Status:         status,
		Reference:      reference,
		FailureCode:    failureCode,
		FailureMessage: failureMessage,
	}, nil
}

func extractActionField(raw string, decoded any, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if path == "$" {
		return strings.TrimSpace(raw)
	}
	return extractByPath(decoded, path)
}

func normalizeActionStatus(status string, cfg ActionConfig) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized == "" {
		return ""
	}
	if len(cfg.StatusMap) > 0 {
		if mapped, ok := cfg.StatusMap[normalized]; ok {
			normalized = mapped
		}
	}

	switch normalized {
	case repository.OrderPaymentStatusSucceeded, "success", "paid", "ok":
		return repository.OrderPaymentStatusSucceeded
	case repository.OrderPaymentStatusFailed, "fail", "error", "invalid":
		return repository.OrderPaymentStatusFailed
	case repository.OrderPaymentStatusPending, "processing", "in_progress":
		return repository.OrderPaymentStatusPending
	default:
		return ""
	}
}

func gatewayTimeout(cfg HTTPGateway) time.Duration {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		return defaultGatewayTimeout
	}
	return timeout
}

func computeSignature(signatureType, secret string, body []byte) (string, error) {
	switch strings.ToLower(strings.TrimSpace(signatureType)) {
	case "", "hmac_sha256":
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		return hex.EncodeToString(mac.Sum(nil)), nil
	default:
		return "", repository.ErrInvalidArgument
	}
}

func signatureEqual(signature, expected string) bool {
	left := strings.ToLower(strings.TrimSpace(signature))
	right := strings.ToLower(strings.TrimSpace(expected))
	if left == "" || right == "" {
		return false
	}
	if len(left) != len(right) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(left), []byte(right)) == 1
}
