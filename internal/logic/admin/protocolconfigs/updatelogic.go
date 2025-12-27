package protocolconfigs

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic handles protocol config updates.
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

// Update updates a protocol config.
func (l *UpdateLogic) Update(req *types.AdminUpdateProtocolConfigRequest) (*types.ProtocolConfigSummary, error) {
	var input repository.UpdateProtocolConfigInput

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		input.Name = &name
	}
	if req.Protocol != nil {
		proto, ok := normalizeProtocol(*req.Protocol)
		if !ok {
			return nil, repository.ErrInvalidArgument
		}
		input.Protocol = &proto
	}
	if req.Status != nil {
		status := strings.TrimSpace(*req.Status)
		input.Status = &status
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
		profile := req.Profile
		input.Profile = &profile
	}

	updated, err := l.svcCtx.Repositories.ProtocolConfig.Update(l.ctx, req.ConfigID, input)
	if err != nil {
		return nil, err
	}

	summary := mapProtocolConfigSummary(updated)
	return &summary, nil
}
