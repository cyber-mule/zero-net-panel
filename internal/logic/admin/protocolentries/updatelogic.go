package protocolentries

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic handles protocol entry updates.
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

// Update updates an existing entry.
func (l *UpdateLogic) Update(req *types.AdminUpdateProtocolEntryRequest) (*types.ProtocolEntrySummary, error) {
	var input repository.UpdateProtocolEntryInput

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		input.Name = &name
	}

	var binding repository.ProtocolBinding
	if req.BindingID != nil && *req.BindingID > 0 {
		bindingValue, err := l.svcCtx.Repositories.ProtocolBinding.Get(l.ctx, *req.BindingID)
		if err != nil {
			return nil, err
		}
		binding = bindingValue
		input.BindingID = req.BindingID
		if req.Protocol == nil {
			protocol := strings.ToLower(strings.TrimSpace(binding.Protocol))
			input.Protocol = &protocol
		}
	}

	if req.Protocol != nil {
		if binding.ID == 0 {
			return nil, repository.ErrInvalidArgument
		}
		protocol := strings.ToLower(strings.TrimSpace(*req.Protocol))
		if protocol == "" || !strings.EqualFold(protocol, binding.Protocol) {
			return nil, repository.ErrInvalidArgument
		}
		input.Protocol = &protocol
	}

	if req.Status != nil {
		status := strings.TrimSpace(*req.Status)
		input.Status = &status
	}
	if req.EntryAddress != nil {
		entryAddress := strings.TrimSpace(*req.EntryAddress)
		if entryAddress == "" {
			return nil, repository.ErrInvalidArgument
		}
		input.EntryAddress = &entryAddress
	}
	if req.EntryPort != nil {
		if *req.EntryPort <= 0 {
			return nil, repository.ErrInvalidArgument
		}
		input.EntryPort = req.EntryPort
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
		profile := cloneEntryProfile(req.Profile)
		input.Profile = &profile
	}

	updated, err := l.svcCtx.Repositories.ProtocolEntry.Update(l.ctx, req.EntryID, input)
	if err != nil {
		return nil, err
	}

	summary := mapProtocolEntrySummary(updated)
	return &summary, nil
}
