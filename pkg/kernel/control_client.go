package kernel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ControlClient provides access to kernel control plane APIs.
type ControlClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewControlClient constructs a control client.
func NewControlClient(opts HTTPOptions) (*ControlClient, error) {
	if strings.TrimSpace(opts.BaseURL) == "" {
		return nil, fmt.Errorf("kernel control client: base url required")
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &ControlClient{
		baseURL: strings.TrimSuffix(strings.TrimSpace(opts.BaseURL), "/"),
		token:   strings.TrimSpace(opts.Token),
		client:  &http.Client{Timeout: timeout},
	}, nil
}

// UpsertProtocol pushes protocol configuration to the kernel.
func (c *ControlClient) UpsertProtocol(ctx context.Context, req ProtocolUpsertRequest) (ProtocolSummary, error) {
	endpoint := c.buildURL("/protocols")

	payload, err := json.Marshal(req)
	if err != nil {
		return ProtocolSummary{}, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return ProtocolSummary{}, err
	}
	c.applyAuth(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return ProtocolSummary{}, err
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return ProtocolSummary{}, fmt.Errorf("kernel control: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var summary ProtocolSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return ProtocolSummary{}, err
	}
	return summary, nil
}

// GetStatus fetches runtime status snapshot (nodes only when supported).
func (c *ControlClient) GetStatus(ctx context.Context) (StatusResponse, error) {
	endpoint := c.buildURL("/status?include=nodes")

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return StatusResponse{}, err
	}
	c.applyAuth(httpReq)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return StatusResponse{}, err
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return StatusResponse{}, fmt.Errorf("kernel control: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var status StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return StatusResponse{}, err
	}
	return status, nil
}

func (c *ControlClient) buildURL(path string) string {
	if strings.HasSuffix(c.baseURL, "/v1") {
		return c.baseURL + path
	}
	return c.baseURL + "/v1" + path
}

func (c *ControlClient) applyAuth(req *http.Request) {
	if c.token == "" {
		return
	}
	token := c.token
	lower := strings.ToLower(token)
	if strings.HasPrefix(lower, "bearer ") || strings.HasPrefix(lower, "basic ") {
		req.Header.Set("Authorization", token)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
}

func closeBody(body io.ReadCloser) {
	if body == nil {
		return
	}
	if err := body.Close(); err != nil {
		fmt.Printf("kernel control client: close response body: %v\n", err)
	}
}
