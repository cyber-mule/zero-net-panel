package protocolbindings

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/logic/credentialutil"
	"github.com/zero-net-panel/zero-net-panel/internal/logic/subscriptionutil"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
	"github.com/zero-net-panel/zero-net-panel/pkg/kernel"
)

// SyncLogic handles protocol binding sync operations.
type SyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext

	registeredServiceEvents map[serviceEventKey]bool
	serviceEventCallback    string
	registeredNodeEvents    map[nodeEventKey]bool
	nodeEventCallback       string
}

// NewSyncLogic constructs SyncLogic.
func NewSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncLogic {
	return &SyncLogic{
		Logger:                  logx.WithContext(ctx),
		ctx:                     ctx,
		svcCtx:                  svcCtx,
		registeredServiceEvents: make(map[serviceEventKey]bool),
		registeredNodeEvents:    make(map[nodeEventKey]bool),
	}
}

// SyncSingle triggers sync for a single binding.
func (l *SyncLogic) SyncSingle(req *types.AdminSyncProtocolBindingRequest) (*types.ProtocolBindingSyncResult, error) {
	if req.BindingID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	binding, err := l.svcCtx.Repositories.ProtocolBinding.Get(l.ctx, req.BindingID)
	if err != nil {
		return nil, err
	}

	result := l.syncBinding(binding)
	return &result, nil
}

// SyncBatch triggers sync for multiple bindings.
func (l *SyncLogic) SyncBatch(req *types.AdminSyncProtocolBindingsRequest) (*types.AdminSyncProtocolBindingsResponse, error) {
	bindings, err := l.resolveBindings(req)
	if err != nil {
		return nil, err
	}
	if len(bindings) == 0 {
		return nil, repository.ErrInvalidArgument
	}

	results := make([]types.ProtocolBindingSyncResult, 0, len(bindings))
	for _, binding := range bindings {
		results = append(results, l.syncBinding(binding))
	}

	return &types.AdminSyncProtocolBindingsResponse{Results: results}, nil
}

func (l *SyncLogic) resolveBindings(req *types.AdminSyncProtocolBindingsRequest) ([]repository.ProtocolBinding, error) {
	bindingMap := make(map[uint64]repository.ProtocolBinding)

	for _, id := range req.BindingIDs {
		if id == 0 {
			continue
		}
		binding, err := l.svcCtx.Repositories.ProtocolBinding.Get(l.ctx, id)
		if err != nil {
			return nil, err
		}
		bindingMap[binding.ID] = binding
	}

	if len(req.NodeIDs) > 0 {
		nodeIDs := make([]uint64, 0, len(req.NodeIDs))
		for _, id := range req.NodeIDs {
			if id > 0 {
				nodeIDs = append(nodeIDs, id)
			}
		}
		if len(nodeIDs) > 0 {
			bindings, err := l.svcCtx.Repositories.ProtocolBinding.ListByNodeIDs(l.ctx, nodeIDs)
			if err != nil {
				return nil, err
			}
			for _, binding := range bindings {
				bindingMap[binding.ID] = binding
			}
		}
	}

	results := make([]repository.ProtocolBinding, 0, len(bindingMap))
	for _, binding := range bindingMap {
		results = append(results, binding)
	}
	return results, nil
}

func (l *SyncLogic) syncBinding(binding repository.ProtocolBinding) types.ProtocolBindingSyncResult {
	result := types.ProtocolBindingSyncResult{
		BindingID: binding.ID,
		Status:    "error",
		SyncedAt:  time.Now().UTC().Unix(),
	}

	if err := l.ensureNodeEventRegistration(binding); err != nil {
		l.Errorf("kernel event registration failed: %v", err)
	}
	if err := l.ensureServiceEventRegistration(binding); err != nil {
		l.Errorf("kernel service event registration failed: %v", err)
	}

	control, err := l.resolveControlClient(binding)
	if err != nil {
		result.Message = err.Error()
		_, _ = l.updateSyncState(binding, result.Status, result.Message)
		return result
	}

	profile := kernel.NodeProfile{
		ID:          binding.KernelID,
		Role:        binding.Role,
		Protocol:    normalizeBindingProtocol(binding),
		Tags:        mergeTags(binding.Tags),
		Description: strings.TrimSpace(binding.Description),
		Profile:     cloneBindingProfile(binding.Profile),
	}
	if len(profile.Profile) == 0 {
		profile.Profile = map[string]any{}
	}

	if profile.ID == "" {
		result.Message = "kernel_id is required"
		_, _ = l.updateSyncState(binding, result.Status, result.Message)
		return result
	}

	users, err := l.buildKernelUsers(binding)
	if err != nil {
		result.Message = err.Error()
		_, _ = l.updateSyncState(binding, result.Status, result.Message)
		return result
	}

	req := kernel.ProtocolUpsertRequest{
		Listen:  normalizeListen(binding.Listen, binding.AccessPort),
		Connect: binding.Connect,
		Users:   users,
		Profile: profile,
	}

	_, err = control.UpsertProtocol(l.ctx, req)
	if err != nil {
		result.Message = err.Error()
		_, _ = l.updateSyncState(binding, result.Status, result.Message)
		return result
	}

	result.Status = "synced"
	result.Message = "ok"
	_, _ = l.updateSyncState(binding, result.Status, "")
	return result
}

type serviceEventKey struct {
	endpoint string
	token    string
	callback string
}

type nodeEventKey struct {
	endpoint string
	token    string
	callback string
	event    string
}

func (l *SyncLogic) ensureNodeEventRegistration(binding repository.ProtocolBinding) error {
	if l.svcCtx == nil {
		return repository.ErrInvalidState
	}
	callback, err := l.resolveNodeEventCallbackURL()
	if err != nil {
		return err
	}

	endpoint := strings.TrimSpace(binding.Node.ControlEndpoint)
	if endpoint == "" {
		return fmt.Errorf("node control endpoint not configured")
	}
	token := resolveControlToken(binding.Node)
	control, err := kernel.NewControlClient(kernel.HTTPOptions{
		BaseURL: endpoint,
		Token:   token,
		Timeout: l.svcCtx.Config.Kernel.HTTP.Timeout,
	})
	if err != nil {
		return err
	}

	events := []string{"node_added", "node_removed", "node_healthy", "node_degraded", "node_unhealthy"}
	for _, event := range events {
		key := nodeEventKey{
			endpoint: endpoint,
			token:    token,
			callback: callback,
			event:    event,
		}
		if l.registeredNodeEvents[key] {
			continue
		}

		req := kernel.EventRegistrationRequest{
			Event:    event,
			Callback: callback,
		}
		if secret := strings.TrimSpace(l.svcCtx.Config.Webhook.SharedToken); secret != "" {
			req.Secret = secret
		}

		if _, err := control.RegisterEvent(l.ctx, req); err != nil {
			return err
		}
		l.registeredNodeEvents[key] = true
	}
	return nil
}

func (l *SyncLogic) ensureServiceEventRegistration(binding repository.ProtocolBinding) error {
	if l.svcCtx == nil {
		return repository.ErrInvalidState
	}
	callback, err := l.resolveServiceEventCallbackURL()
	if err != nil {
		return err
	}

	endpoint := strings.TrimSpace(binding.Node.ControlEndpoint)
	if endpoint == "" {
		return fmt.Errorf("node control endpoint not configured")
	}
	token := resolveControlToken(binding.Node)
	key := serviceEventKey{
		endpoint: endpoint,
		token:    token,
		callback: callback,
	}
	if l.registeredServiceEvents[key] {
		return nil
	}

	control, err := kernel.NewControlClient(kernel.HTTPOptions{
		BaseURL: endpoint,
		Token:   token,
		Timeout: l.svcCtx.Config.Kernel.HTTP.Timeout,
	})
	if err != nil {
		return err
	}

	req := kernel.ServiceEventRegistrationRequest{
		Event:    "user.traffic.reported",
		Callback: callback,
	}
	if secret := strings.TrimSpace(l.svcCtx.Config.Webhook.SharedToken); secret != "" {
		req.Secret = secret
	}

	if _, err := control.RegisterServiceEvent(l.ctx, req); err != nil {
		return err
	}
	l.registeredServiceEvents[key] = true
	return nil
}

func (l *SyncLogic) resolveNodeEventCallbackURL() (string, error) {
	if l.nodeEventCallback != "" {
		return l.nodeEventCallback, nil
	}
	defaults := repository.SiteSettingDefaults{
		Name:                                 l.svcCtx.Config.Site.Name,
		LogoURL:                              l.svcCtx.Config.Site.LogoURL,
		KernelOfflineProbeMaxIntervalSeconds: int(l.svcCtx.Config.Kernel.OfflineProbeMaxInterval / time.Second),
	}
	setting, err := l.svcCtx.Repositories.Site.GetSiteSetting(l.ctx, defaults)
	if err != nil {
		return "", err
	}
	if raw := strings.TrimSpace(setting.AccessDomain); raw != "" {
		callback, err := buildCallbackFromBase(raw, 0, "/api/v1/kernel/events")
		if err != nil {
			return "", err
		}
		l.nodeEventCallback = callback
		return callback, nil
	}
	callback, err := buildCallbackFromBase(l.svcCtx.Config.Host, l.svcCtx.Config.Port, "/api/v1/kernel/events")
	if err != nil {
		return "", err
	}
	l.nodeEventCallback = callback
	return callback, nil
}

func (l *SyncLogic) resolveServiceEventCallbackURL() (string, error) {
	if l.serviceEventCallback != "" {
		return l.serviceEventCallback, nil
	}
	defaults := repository.SiteSettingDefaults{
		Name:                                 l.svcCtx.Config.Site.Name,
		LogoURL:                              l.svcCtx.Config.Site.LogoURL,
		KernelOfflineProbeMaxIntervalSeconds: int(l.svcCtx.Config.Kernel.OfflineProbeMaxInterval / time.Second),
	}
	setting, err := l.svcCtx.Repositories.Site.GetSiteSetting(l.ctx, defaults)
	if err != nil {
		return "", err
	}
	if raw := strings.TrimSpace(setting.AccessDomain); raw != "" {
		callback, err := buildCallbackFromBase(raw, 0, "/api/v1/kernel/service-events")
		if err != nil {
			return "", err
		}
		l.serviceEventCallback = callback
		return callback, nil
	}
	callback, err := buildCallbackFromBase(l.svcCtx.Config.Host, l.svcCtx.Config.Port, "/api/v1/kernel/service-events")
	if err != nil {
		return "", err
	}
	l.serviceEventCallback = callback
	return callback, nil
}

func buildCallbackFromBase(raw string, port int, path string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("callback host not configured")
	}

	var base url.URL
	if strings.HasPrefix(strings.ToLower(raw), "http://") || strings.HasPrefix(strings.ToLower(raw), "https://") {
		parsed, err := url.Parse(raw)
		if err != nil {
			return "", err
		}
		base = *parsed
	} else {
		base = url.URL{
			Scheme: "http",
			Host:   raw,
		}
	}

	if port > 0 && !strings.Contains(base.Host, ":") {
		base.Host = fmt.Sprintf("%s:%d", base.Host, port)
	}

	path = strings.TrimSpace(path)
	if path == "" {
		path = "/"
	}
	base.Path = strings.TrimSuffix(base.Path, "/") + path
	return base.String(), nil
}

func (l *SyncLogic) buildKernelUsers(binding repository.ProtocolBinding) ([]kernel.User, error) {
	planIDs, err := l.svcCtx.Repositories.PlanProtocolBinding.ListPlanIDsByBindingID(l.ctx, binding.ID)
	if err != nil {
		return nil, err
	}
	if len(planIDs) == 0 {
		return []kernel.User{}, nil
	}

	subs, err := l.svcCtx.Repositories.Subscription.ListActiveByPlanIDs(l.ctx, planIDs)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	users := make([]kernel.User, 0, len(subs))
	seen := make(map[uint64]struct{}, len(subs))

	for _, sub := range subs {
		if !subscriptionutil.IsSubscriptionEffective(sub, now) {
			continue
		}
		if _, ok := seen[sub.UserID]; ok {
			continue
		}
		seen[sub.UserID] = struct{}{}

		credential, err := credentialutil.EnsureActiveCredential(l.ctx, l.svcCtx.Repositories, l.svcCtx.Credentials, sub.UserID)
		if err != nil {
			l.Errorf("kernel user credential missing user_id=%d: %v", sub.UserID, err)
			continue
		}
		identity, err := credentialutil.BuildIdentity(l.svcCtx.Credentials, sub.UserID, credential)
		if err != nil {
			l.Errorf("kernel user identity build failed user_id=%d: %v", sub.UserID, err)
			continue
		}

		users = append(users, kernel.User{
			ID:       strconv.FormatUint(sub.UserID, 10),
			Username: strings.TrimSpace(identity.Username),
			Password: strings.TrimSpace(identity.Password),
			Metadata: map[string]any{
				"subscription_id": sub.ID,
			},
		})
	}

	return users, nil
}

func (l *SyncLogic) resolveControlClient(binding repository.ProtocolBinding) (*kernel.ControlClient, error) {
	endpoint := strings.TrimSpace(binding.Node.ControlEndpoint)
	token := resolveControlToken(binding.Node)
	if endpoint == "" {
		return nil, fmt.Errorf("node control endpoint not configured")
	}

	opts := kernel.HTTPOptions{
		BaseURL: endpoint,
		Token:   token,
		Timeout: l.svcCtx.Config.Kernel.HTTP.Timeout,
	}
	return kernel.NewControlClient(opts)
}

func resolveControlToken(node repository.Node) string {
	accessKey := strings.TrimSpace(node.ControlAccessKey)
	secretKey := strings.TrimSpace(node.ControlSecretKey)
	if accessKey != "" && secretKey != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(accessKey + ":" + secretKey))
		return "Basic " + encoded
	}
	token := strings.TrimSpace(node.ControlToken)
	if token != "" {
		return token
	}
	return ""
}

func (l *SyncLogic) updateSyncState(binding repository.ProtocolBinding, status string, message string) (repository.ProtocolBinding, error) {
	ts := time.Now().UTC()
	input := repository.UpdateProtocolBindingInput{
		SyncStatus:    &status,
		LastSyncedAt:  &ts,
		LastSyncError: &message,
	}
	return l.svcCtx.Repositories.ProtocolBinding.UpdateSyncState(l.ctx, binding.ID, input)
}

func mergeTags(group ...[]string) []string {
	seen := make(map[string]struct{})
	var result []string
	appendTag := func(tag string) {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return
		}
		if _, ok := seen[tag]; ok {
			return
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}
	for _, tags := range group {
		for _, tag := range tags {
			appendTag(tag)
		}
	}
	return result
}

func normalizeListen(listen string, accessPort int) string {
	listen = strings.TrimSpace(listen)
	if listen == "" {
		if accessPort > 0 {
			return fmt.Sprintf("0.0.0.0:%d", accessPort)
		}
		return ""
	}
	if !strings.Contains(listen, ":") {
		return fmt.Sprintf("0.0.0.0:%s", listen)
	}
	return listen
}
