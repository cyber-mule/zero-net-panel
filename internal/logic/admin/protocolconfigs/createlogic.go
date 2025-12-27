package protocolconfigs

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles protocol config creation.
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

// Create creates a protocol config.
func (l *CreateLogic) Create(req *types.AdminCreateProtocolConfigRequest) (*types.ProtocolConfigSummary, error) {
	protocol, ok := normalizeProtocol(req.Protocol)
	if !ok {
		return nil, repository.ErrInvalidArgument
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}

	cfg := repository.ProtocolConfig{
		Name:        strings.TrimSpace(req.Name),
		Protocol:    protocol,
		Status:      status,
		Tags:        append([]string(nil), req.Tags...),
		Description: strings.TrimSpace(req.Description),
		Profile:     req.Profile,
	}

	created, err := l.svcCtx.Repositories.ProtocolConfig.Create(l.ctx, cfg)
	if err != nil {
		return nil, err
	}

	summary := mapProtocolConfigSummary(created)
	return &summary, nil
}
