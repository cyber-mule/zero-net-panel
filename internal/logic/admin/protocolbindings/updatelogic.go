package protocolbindings

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic handles protocol binding updates.
type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLogic constructs UpdateLogic.
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update updates an existing binding.
func (l *UpdateLogic) Update(req *types.AdminUpdateProtocolBindingRequest) (*types.ProtocolBindingSummary, error) {
	var input repository.UpdateProtocolBindingInput

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		input.Name = &name
	}
	if req.NodeID != nil && *req.NodeID > 0 {
		if _, err := l.svcCtx.Repositories.Node.Get(l.ctx, *req.NodeID); err != nil {
			return nil, err
		}
		input.NodeID = req.NodeID
	}
	if req.Protocol != nil {
		protocol := strings.ToLower(strings.TrimSpace(*req.Protocol))
		if protocol == "" {
			return nil, repository.ErrInvalidArgument
		}
		input.Protocol = &protocol
	}
	if req.Role != nil {
		role, ok := normalizeRole(*req.Role)
		if !ok {
			return nil, repository.ErrInvalidArgument
		}
		input.Role = &role
	}
	if req.Listen != nil {
		listen := strings.TrimSpace(*req.Listen)
		input.Listen = &listen
	}
	if req.Connect != nil {
		connect := strings.TrimSpace(*req.Connect)
		input.Connect = &connect
	}
	if req.AccessPort != nil {
		if *req.AccessPort < 0 {
			return nil, repository.ErrInvalidArgument
		}
		input.AccessPort = req.AccessPort
	}
	if req.Status != nil {
		statusCode, err := normalizeBindingStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		input.Status = &statusCode
	}
	if req.KernelID != nil {
		kernelID := strings.TrimSpace(*req.KernelID)
		if kernelID == "" {
			return nil, repository.ErrInvalidArgument
		}
		input.KernelID = &kernelID
	}
	if req.SyncStatus != nil {
		syncStatus, err := normalizeBindingSyncStatus(*req.SyncStatus)
		if err != nil {
			return nil, err
		}
		input.SyncStatus = &syncStatus
	}
	if req.HealthStatus != nil {
		healthStatus, err := normalizeBindingHealthStatus(*req.HealthStatus)
		if err != nil {
			return nil, err
		}
		input.HealthStatus = &healthStatus
	}
	if req.LastSyncedAt != nil {
		ts := time.Unix(*req.LastSyncedAt, 0).UTC()
		input.LastSyncedAt = &ts
	}
	if req.LastHeartbeatAt != nil {
		ts := time.Unix(*req.LastHeartbeatAt, 0).UTC()
		input.LastHeartbeatAt = &ts
	}
	if req.LastSyncError != nil {
		errMsg := strings.TrimSpace(*req.LastSyncError)
		input.LastSyncError = &errMsg
	}
	if req.Tags != nil {
		tags := append([]string(nil), req.Tags...)
		input.Tags = &tags
	}
	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)
		input.Description = &description
	}
	if req.Profile != nil {
		profile := cloneBindingProfile(req.Profile)
		input.Profile = &profile
	}
	if req.Metadata != nil {
		metadata := req.Metadata
		input.Metadata = &metadata
	}

	updated, err := l.svcCtx.Repositories.ProtocolBinding.Update(l.ctx, req.BindingID, input)
	if err != nil {
		return nil, err
	}

	summary := mapProtocolBindingSummary(updated)
	return &summary, nil
}
