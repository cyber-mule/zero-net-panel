package subscription

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/logic/credentialutil"
	subscriptionutil "github.com/zero-net-panel/zero-net-panel/internal/logic/subscriptionutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	subtemplate "github.com/zero-net-panel/zero-net-panel/pkg/subscription/template"
)

// DownloadResult carries rendered subscription output.
type DownloadResult struct {
	Content     string
	ContentType string
	ETag        string
	TemplateID  uint64
}

// DownloadLogic renders public subscription output.
type DownloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDownloadLogic constructs DownloadLogic.
func NewDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadLogic {
	return &DownloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Download renders a subscription using the client User-Agent to pick templates.
func (l *DownloadLogic) Download(token, userAgent string) (DownloadResult, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return DownloadResult{}, repository.ErrInvalidArgument
	}

	sub, err := l.svcCtx.Repositories.Subscription.GetByToken(l.ctx, token)
	if err != nil {
		return DownloadResult{}, err
	}

	now := time.Now().UTC()
	if !isSubscriptionActive(sub, now) {
		return DownloadResult{}, repository.ErrNotFound
	}

	tpl, err := l.resolveTemplate(sub, userAgent)
	if err != nil {
		return DownloadResult{}, err
	}

	entries, err := subscriptionutil.LoadSubscriptionEntries(l.ctx, l.svcCtx.Repositories, sub)
	if err != nil {
		return DownloadResult{}, err
	}

	credential, err := credentialutil.EnsureActiveCredential(l.ctx, l.svcCtx.Repositories, l.svcCtx.Credentials, sub.UserID)
	if err != nil {
		return DownloadResult{}, err
	}
	identity, err := credentialutil.BuildIdentity(l.svcCtx.Credentials, sub.UserID, credential)
	if err != nil {
		return DownloadResult{}, err
	}

	identityData := map[string]any{
		"version":    credential.Version,
		"status":     credential.Status,
		"account_id": identity.AccountID,
		"account":    identity.AccountID,
		"password":   identity.Password,
		"id":         identity.ID,
		"uuid":       identity.UUID,
		"username":   identity.Username,
		"secret":     identity.Secret,
	}

	entryContext := normalizeEntryContext(entries)
	data := map[string]any{
		"subscription": map[string]any{
			"id":                      sub.ID,
			"name":                    sub.Name,
			"plan":                    sub.PlanName,
			"plan_id":                 sub.PlanID,
			"plan_snapshot":           sub.PlanSnapshot,
			"status":                  sub.Status,
			"token":                   sub.Token,
			"expires_at":              sub.ExpiresAt.Format(time.RFC3339),
			"traffic_total_bytes":     sub.TrafficTotalBytes,
			"traffic_used_bytes":      sub.TrafficUsedBytes,
			"traffic_remaining_bytes": maxInt64(sub.TrafficTotalBytes-sub.TrafficUsedBytes, 0),
			"devices_limit":           sub.DevicesLimit,
			"available_template_ids":  sub.AvailableTemplateIDs,
		},
		"nodes":             entryContext,
		"protocol_bindings": entryContext,
		"user_identity":     identityData,
		"template": map[string]any{
			"id":      tpl.ID,
			"name":    tpl.Name,
			"format":  tpl.Format,
			"version": tpl.Version,
		},
		"generated_at": now.Format(time.RFC3339),
	}

	content, err := subtemplate.Render(tpl.Format, tpl.Content, data)
	if err != nil {
		return DownloadResult{}, err
	}

	hash := sha256.Sum256([]byte(content))
	etag := hex.EncodeToString(hash[:])

	contentType := "text/plain; charset=utf-8"
	if tpl.Format == "json" {
		contentType = "application/json"
	}

	return DownloadResult{
		Content:     content,
		ContentType: contentType,
		ETag:        etag,
		TemplateID:  tpl.ID,
	}, nil
}

func (l *DownloadLogic) resolveTemplate(sub repository.Subscription, userAgent string) (repository.SubscriptionTemplate, error) {
	clientType := detectClientType(userAgent)
	available := append([]uint64(nil), sub.AvailableTemplateIDs...)
	if sub.TemplateID != 0 && !containsUint64(available, sub.TemplateID) {
		available = append(available, sub.TemplateID)
	}
	if len(available) == 0 {
		return repository.SubscriptionTemplate{}, repository.ErrNotFound
	}

	templates, err := l.svcCtx.Repositories.SubscriptionTemplate.ListByIDs(l.ctx, available)
	if err != nil {
		return repository.SubscriptionTemplate{}, err
	}

	template, ok := selectTemplate(templates, available, sub.TemplateID, clientType)
	if !ok {
		return repository.SubscriptionTemplate{}, repository.ErrNotFound
	}
	return template, nil
}

func selectTemplate(templates []repository.SubscriptionTemplate, orderedIDs []uint64, fallbackID uint64, clientType string) (repository.SubscriptionTemplate, bool) {
	byID := make(map[uint64]repository.SubscriptionTemplate, len(templates))
	for _, tpl := range templates {
		byID[tpl.ID] = tpl
	}

	clientType = strings.ToLower(strings.TrimSpace(clientType))
	if clientType != "" {
		for _, id := range orderedIDs {
			tpl, ok := byID[id]
			if !ok || !strings.EqualFold(tpl.ClientType, clientType) || !tpl.IsDefault {
				continue
			}
			return tpl, true
		}
		for _, id := range orderedIDs {
			tpl, ok := byID[id]
			if !ok || !strings.EqualFold(tpl.ClientType, clientType) {
				continue
			}
			return tpl, true
		}
	}

	if fallbackID != 0 {
		if tpl, ok := byID[fallbackID]; ok {
			return tpl, true
		}
	}
	for _, id := range orderedIDs {
		if tpl, ok := byID[id]; ok {
			return tpl, true
		}
	}

	return repository.SubscriptionTemplate{}, false
}

func detectClientType(userAgent string) string {
	ua := strings.ToLower(strings.TrimSpace(userAgent))
	if ua == "" {
		return ""
	}

	type rule struct {
		client string
		tokens []string
	}
	rules := []rule{
		{client: "sing-box", tokens: []string{"sing-box", "singbox"}},
		{client: "clash", tokens: []string{"clash", "clash.meta", "clashmeta", "clash-verge", "clash verge", "clashx", "clashforwindows", "mihomo", "mihomo-party"}},
		{client: "surge", tokens: []string{"surge"}},
		{client: "quantumult", tokens: []string{"quantumult", "quantumult x", "quantumultx"}},
		{client: "stash", tokens: []string{"stash"}},
		{client: "shadowrocket", tokens: []string{"shadowrocket"}},
		{client: "loon", tokens: []string{"loon"}},
		{client: "hiddify", tokens: []string{"hiddify"}},
		{client: "v2rayn", tokens: []string{"v2rayn"}},
		{client: "v2rayng", tokens: []string{"v2rayng"}},
		{client: "nekobox", tokens: []string{"nekobox", "neko box"}},
		{client: "kitsunebi", tokens: []string{"kitsunebi"}},
		{client: "potatso", tokens: []string{"potatso"}},
		{client: "surfboard", tokens: []string{"surfboard"}},
	}

	for _, item := range rules {
		for _, token := range item.tokens {
			if strings.Contains(ua, token) {
				return item.client
			}
		}
	}
	return ""
}

func isSubscriptionActive(sub repository.Subscription, now time.Time) bool {
	if sub.Status != status.SubscriptionStatusActive {
		return false
	}
	if sub.ExpiresAt.IsZero() {
		return true
	}
	return sub.ExpiresAt.After(now)
}

func normalizeEntryContext(entries []repository.ProtocolEntry) []map[string]any {
	result := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		if !entryVisible(entry) {
			continue
		}
		binding := entry.Binding
		node := binding.Node
		address := selectEntryAddress(entry)
		host, port := splitHostPort(address)
		result = append(result, map[string]any{
			"id":             binding.ID,
			"binding_id":     binding.ID,
			"entry_id":       entry.ID,
			"kernel_id":      binding.KernelID,
			"protocol":       binding.Protocol,
			"role":           binding.Role,
			"hostname":       host,
			"port":           port,
			"listen":         binding.Listen,
			"connect":        binding.Connect,
			"access_address": entry.EntryAddress,
			"access_port":    entry.EntryPort,
			"entry_address":  entry.EntryAddress,
			"entry_port":     entry.EntryPort,
			"node_id":        binding.NodeID,
			"node_name":      node.Name,
			"region":         node.Region,
			"country":        node.Country,
			"status":         entry.Status,
			"binding_status": binding.Status,
			"health_status":  binding.HealthStatus,
			"profile":        cloneEntryProfile(entry.Profile),
			"updated_at":     entry.UpdatedAt.Format(time.RFC3339),
		})
	}
	return result
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func containsUint64(list []uint64, target uint64) bool {
	for _, value := range list {
		if value == target {
			return true
		}
	}
	return false
}

func entryVisible(entry repository.ProtocolEntry) bool {
	if entry.Status != status.ProtocolEntryStatusActive {
		return false
	}
	if entry.Binding.Status != status.ProtocolBindingStatusActive {
		return false
	}
	return true
}

func selectEntryAddress(entry repository.ProtocolEntry) string {
	address := strings.TrimSpace(entry.EntryAddress)
	if address != "" && entry.EntryPort > 0 {
		return net.JoinHostPort(address, strconv.Itoa(entry.EntryPort))
	}
	if address != "" {
		return address
	}
	nodeAddress := strings.TrimSpace(entry.Binding.Node.AccessAddress)
	if nodeAddress != "" && entry.EntryPort > 0 {
		return net.JoinHostPort(nodeAddress, strconv.Itoa(entry.EntryPort))
	}
	return nodeAddress
}

func splitHostPort(address string) (string, int) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", 0
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return address, 0
	}
	value, _ := strconv.Atoi(port)
	return host, value
}

func cloneEntryProfile(profile map[string]any) map[string]any {
	if profile == nil {
		return nil
	}
	cloned := make(map[string]any, len(profile))
	for key, value := range profile {
		cloned[key] = value
	}
	return cloned
}
