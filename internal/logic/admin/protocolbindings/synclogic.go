package protocolbindings

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

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
}

// NewSyncLogic constructs SyncLogic.
func NewSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncLogic {
	return &SyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
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
		Tags:        mergeTags(binding.ProtocolConfig.Tags, binding.Tags),
		Description: firstNonEmpty(binding.Description, binding.ProtocolConfig.Description),
		Profile:     cloneBindingProfile(binding.Profile),
	}
	if len(profile.Profile) == 0 {
		profile.Profile = cloneBindingProfile(binding.ProtocolConfig.Profile)
	}
	if len(profile.Profile) == 0 {
		profile.Profile = map[string]any{}
	}

	if profile.ID == "" {
		result.Message = "kernel_id is required"
		_, _ = l.updateSyncState(binding, result.Status, result.Message)
		return result
	}

	req := kernel.ProtocolUpsertRequest{
		Listen:  binding.Listen,
		Connect: binding.Connect,
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

func mergeTags(base []string, extra []string) []string {
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
	for _, tag := range base {
		appendTag(tag)
	}
	for _, tag := range extra {
		appendTag(tag)
	}
	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
