package nodes

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// DisableLogic handles node disabling.
type DisableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDisableLogic constructs DisableLogic.
func NewDisableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableLogic {
	return &DisableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Disable marks a node as disabled.
func (l *DisableLogic) Disable(req *types.AdminDisableNodeRequest) (*types.AdminNodeResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}

	statusCode := status.NodeStatusDisabled
	input := repository.UpdateNodeInput{
		Status: &statusCode,
	}

	var updated repository.Node
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		node, err := txRepos.Node.Update(l.ctx, req.NodeID, input)
		if err != nil {
			return err
		}
		updated = node

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.node.disable",
			ResourceType: "node",
			ResourceID:   fmt.Sprintf("%d", updated.ID),
			Metadata: map[string]any{
				"status": statusCode,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminNodeResponse{
		Node: mapNodeSummary(updated),
	}, nil
}
