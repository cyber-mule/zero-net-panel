package protocolentries

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles protocol entry creation.
type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic constructs CreateLogic.
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create creates a protocol entry.
func (l *CreateLogic) Create(req *types.AdminCreateProtocolEntryRequest) (*types.ProtocolEntrySummary, error) {
	if req.BindingID == 0 {
		return nil, repository.ErrInvalidArgument
	}
	binding, err := l.svcCtx.Repositories.ProtocolBinding.Get(l.ctx, req.BindingID)
	if err != nil {
		return nil, err
	}

	protocol := strings.ToLower(strings.TrimSpace(req.Protocol))
	if protocol == "" {
		protocol = strings.ToLower(strings.TrimSpace(binding.Protocol))
	}
	if protocol == "" || !strings.EqualFold(protocol, binding.Protocol) {
		return nil, repository.ErrInvalidArgument
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = strings.TrimSpace(binding.Name)
	}
	if name == "" {
		return nil, repository.ErrInvalidArgument
	}

	entryAddress := strings.TrimSpace(req.EntryAddress)
	if entryAddress == "" || req.EntryPort <= 0 {
		return nil, repository.ErrInvalidArgument
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}

	profile := cloneEntryProfile(req.Profile)

	entry := repository.ProtocolEntry{
		Name:         name,
		BindingID:    binding.ID,
		Protocol:     protocol,
		Status:       status,
		EntryAddress: entryAddress,
		EntryPort:    req.EntryPort,
		Tags:         append([]string(nil), req.Tags...),
		Description:  strings.TrimSpace(req.Description),
		Profile:      profile,
	}

	created, err := l.svcCtx.Repositories.ProtocolEntry.Create(l.ctx, entry)
	if err != nil {
		return nil, err
	}

	summary := mapProtocolEntrySummary(created)
	return &summary, nil
}
