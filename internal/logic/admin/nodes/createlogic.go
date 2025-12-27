package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles node creation.
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

// Create provisions a new node.
func (l *CreateLogic) Create(req *types.AdminCreateNodeRequest) (*types.AdminNodeResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, repository.ErrInvalidArgument
	}

	status := "offline"
	if strings.TrimSpace(req.Status) != "" {
		normalized, err := normalizeNodeStatus(req.Status)
		if err != nil {
			return nil, err
		}
		status = normalized
	}

	tags := normalizeTags(req.Tags)
	protocols := normalizeProtocols(req.Protocols)
	if req.CapacityMbps < 0 {
		return nil, repository.ErrInvalidArgument
	}

	now := time.Now().UTC()
	node := repository.Node{
		Name:         name,
		Region:       strings.TrimSpace(req.Region),
		Country:      strings.TrimSpace(req.Country),
		ISP:          strings.TrimSpace(req.ISP),
		Status:       status,
		Tags:         tags,
		Protocols:    protocols,
		CapacityMbps: req.CapacityMbps,
		Description:  strings.TrimSpace(req.Description),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	var created repository.Node
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		createdNode, err := txRepos.Node.Create(l.ctx, node)
		if err != nil {
			return err
		}
		created = createdNode

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.node.create",
			ResourceType: "node",
			ResourceID:   fmt.Sprintf("%d", created.ID),
			Metadata: map[string]any{
				"name":      created.Name,
				"status":    created.Status,
				"protocols": created.Protocols,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminNodeResponse{
		Node: mapNodeSummary(created),
	}, nil
}
