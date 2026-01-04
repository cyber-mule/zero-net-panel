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
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
	subtemplate "github.com/zero-net-panel/zero-net-panel/pkg/subscription/template"
)

// PreviewLogic 渲染订阅预览。
type PreviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPreviewLogic 构造函数。
func NewPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewLogic {
	return &PreviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Preview 生成订阅预览。
func (l *PreviewLogic) Preview(req *types.UserSubscriptionPreviewRequest) (*types.UserSubscriptionPreviewResponse, error) {
	user, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrForbidden
	}

	sub, err := l.svcCtx.Repositories.Subscription.Get(l.ctx, req.SubscriptionID)
	if err != nil {
		return nil, err
	}

	if sub.UserID != user.ID {
		return nil, repository.ErrForbidden
	}
	if strings.EqualFold(sub.Status, "disabled") {
		return nil, repository.ErrNotFound
	}

	templateID := req.TemplateID
	if templateID == 0 {
		templateID = sub.TemplateID
	} else {
		if !isTemplateAllowed(sub, templateID) {
			return nil, repository.ErrForbidden
		}
	}

	tpl, err := l.svcCtx.Repositories.SubscriptionTemplate.Get(l.ctx, templateID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	isActive := strings.EqualFold(sub.Status, "active")
	identityData := map[string]any{
		"version":    0,
		"status":     "",
		"account_id": "",
		"account":    "",
		"password":   "",
		"id":         "",
		"uuid":       "",
		"username":   "",
		"secret":     "",
	}
	if isActive {
		credential, err := credentialutil.EnsureActiveCredential(l.ctx, l.svcCtx.Repositories, l.svcCtx.Credentials, user.ID)
		if err != nil {
			return nil, err
		}
		identity, err := credentialutil.BuildIdentity(l.svcCtx.Credentials, user.ID, credential)
		if err != nil {
			return nil, err
		}
		identityData = map[string]any{
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
	}

	entries := []repository.ProtocolEntry{}
	if isActive {
		entries, err = subscriptionutil.LoadSubscriptionEntries(l.ctx, l.svcCtx.Repositories, sub)
		if err != nil {
			return nil, err
		}
	}

	entryContext := normalizeEntryContext(entries)
	if !isActive {
		entryContext = []map[string]any{}
	}
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
		return nil, err
	}

	hash := sha256.Sum256([]byte(content))
	etag := hex.EncodeToString(hash[:])

	contentType := "text/plain; charset=utf-8"
	switch tpl.Format {
	case "json":
		contentType = "application/json"
	}

	return &types.UserSubscriptionPreviewResponse{
		SubscriptionID: sub.ID,
		TemplateID:     templateID,
		Content:        content,
		ContentType:    contentType,
		ETag:           etag,
		GeneratedAt:    now.Unix(),
	}, nil
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

func isTemplateAllowed(sub repository.Subscription, templateID uint64) bool {
	if templateID == sub.TemplateID {
		return true
	}
	for _, id := range sub.AvailableTemplateIDs {
		if id == templateID {
			return true
		}
	}
	return false
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func entryVisible(entry repository.ProtocolEntry) bool {
	if strings.ToLower(entry.Status) != "active" {
		return false
	}
	if strings.ToLower(entry.Binding.Status) != "active" {
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
