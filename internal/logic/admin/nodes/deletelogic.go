package nodes

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// DeleteLogic handles node deletion.
type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteLogic constructs DeleteLogic.
func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Delete removes a node and related bindings.
func (l *DeleteLogic) Delete(req *types.AdminDeleteNodeRequest) error {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return repository.ErrForbidden
	}

	node, err := l.svcCtx.Repositories.Node.Get(l.ctx, req.NodeID)
	if err != nil {
		return err
	}

	return l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		bindings, err := txRepos.ProtocolBinding.ListByNodeIDs(l.ctx, []uint64{node.ID})
		if err != nil {
			return err
		}

		bindingIDs := make([]uint64, 0, len(bindings))
		for _, binding := range bindings {
			if binding.ID == 0 {
				continue
			}
			bindingIDs = append(bindingIDs, binding.ID)
		}

		if len(bindingIDs) > 0 {
			if err := txRepos.PlanProtocolBinding.DeleteByBindingIDs(l.ctx, bindingIDs); err != nil {
				return err
			}
			for _, bindingID := range bindingIDs {
				if err := txRepos.ProtocolBinding.Delete(l.ctx, bindingID); err != nil {
					return err
				}
			}
		}

		if err := txRepos.Node.Delete(l.ctx, node.ID); err != nil {
			return err
		}

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.node.delete",
			ResourceType: "node",
			ResourceID:   fmt.Sprintf("%d", node.ID),
			Metadata: map[string]any{
				"node_name":   node.Name,
				"binding_ids": bindingIDs,
			},
		})
		return err
	})
}
