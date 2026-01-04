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

	protocol := strings.ToLower(strings.TrimSpace(req.Protocol))
	if protocol == "" {
		return nil, repository.ErrInvalidArgument
	}
	kernelID := strings.TrimSpace(req.KernelID)
	if kernelID == "" {
		return nil, repository.ErrInvalidArgument
	}
	if req.AccessPort < 0 {
		return nil, repository.ErrInvalidArgument
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}

	if req.Profile == nil {
		return nil, repository.ErrInvalidArgument
	}
	profile := req.Profile

	binding := repository.ProtocolBinding{
		Name:        strings.TrimSpace(req.Name),
		NodeID:      req.NodeID,
		Protocol:    protocol,
		Role:        role,
		Listen:      strings.TrimSpace(req.Listen),
		Connect:     strings.TrimSpace(req.Connect),
		AccessPort:  req.AccessPort,
		Status:      status,
		KernelID:    kernelID,
		Tags:        append([]string(nil), req.Tags...),
		Description: strings.TrimSpace(req.Description),
		Profile:     cloneBindingProfile(profile),
		Metadata:    req.Metadata,
	}

	created, err := l.svcCtx.Repositories.ProtocolBinding.Create(l.ctx, binding)
	if err != nil {
		return nil, err
	}

	summary := mapProtocolBindingSummary(created)
	return &summary, nil
}

func cloneBindingProfile(profile map[string]any) map[string]any {
	if profile == nil {
		return nil
	}
	cloned := make(map[string]any, len(profile))
	for key, value := range profile {
		cloned[key] = value
	}
	return cloned
}
