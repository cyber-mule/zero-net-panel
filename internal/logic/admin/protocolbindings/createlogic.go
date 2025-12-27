package protocolbindings

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles protocol binding creation.
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

// Create creates a new binding.
func (l *CreateLogic) Create(req *types.AdminCreateProtocolBindingRequest) (*types.ProtocolBindingSummary, error) {
	role, ok := normalizeRole(req.Role)
	if !ok {
		return nil, repository.ErrInvalidArgument
	}
	if _, err := l.svcCtx.Repositories.Node.Get(l.ctx, req.NodeID); err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.Repositories.ProtocolConfig.Get(l.ctx, req.ProtocolConfigID); err != nil {
		return nil, err
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}
	binding := repository.ProtocolBinding{
		Name:             strings.TrimSpace(req.Name),
		NodeID:           req.NodeID,
		ProtocolConfigID: req.ProtocolConfigID,
		Role:             role,
		Listen:           strings.TrimSpace(req.Listen),
		Connect:          strings.TrimSpace(req.Connect),
		Status:           status,
		KernelID:         strings.TrimSpace(req.KernelID),
		Tags:             append([]string(nil), req.Tags...),
		Description:      strings.TrimSpace(req.Description),
		Metadata:         req.Metadata,
	}

	created, err := l.svcCtx.Repositories.ProtocolBinding.Create(l.ctx, binding)
	if err != nil {
		return nil, err
	}

	summary := mapProtocolBindingSummary(created)
	return &summary, nil
}
