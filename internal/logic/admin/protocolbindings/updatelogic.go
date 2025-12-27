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
	if req.ProtocolConfigID != nil && *req.ProtocolConfigID > 0 {
		if _, err := l.svcCtx.Repositories.ProtocolConfig.Get(l.ctx, *req.ProtocolConfigID); err != nil {
			return nil, err
		}
		input.ProtocolConfigID = req.ProtocolConfigID
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
	if req.Status != nil {
		status := strings.TrimSpace(*req.Status)
		input.Status = &status
	}
	if req.KernelID != nil {
		kernelID := strings.TrimSpace(*req.KernelID)
		input.KernelID = &kernelID
	}
	if req.SyncStatus != nil {
		syncStatus := strings.TrimSpace(*req.SyncStatus)
		input.SyncStatus = &syncStatus
	}
	if req.HealthStatus != nil {
		healthStatus := strings.TrimSpace(*req.HealthStatus)
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
