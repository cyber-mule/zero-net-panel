package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic handles node updates.
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

// Update patches node fields.
func (l *UpdateLogic) Update(req *types.AdminUpdateNodeRequest) (*types.AdminNodeResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}

	input := repository.UpdateNodeInput{}
	metadata := map[string]any{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, repository.ErrInvalidArgument
		}
		input.Name = &name
		metadata["name"] = name
	}
	if req.Region != nil {
		region := strings.TrimSpace(*req.Region)
		input.Region = &region
		metadata["region"] = region
	}
	if req.Country != nil {
		country := strings.TrimSpace(*req.Country)
		input.Country = &country
		metadata["country"] = country
	}
	if req.ISP != nil {
		isp := strings.TrimSpace(*req.ISP)
		input.ISP = &isp
		metadata["isp"] = isp
	}
	if req.Status != nil {
		status, err := normalizeNodeStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		input.Status = &status
		metadata["status"] = status
	}
	if req.Tags != nil {
		tags := normalizeTags(req.Tags)
		input.Tags = &tags
		metadata["tags"] = tags
	}
	if req.Protocols != nil {
		protocols := normalizeProtocols(req.Protocols)
		input.Protocols = &protocols
		metadata["protocols"] = protocols
	}
	if req.CapacityMbps != nil {
		if *req.CapacityMbps < 0 {
			return nil, repository.ErrInvalidArgument
		}
		input.CapacityMbps = req.CapacityMbps
		metadata["capacity_mbps"] = *req.CapacityMbps
	}
	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		input.Description = &desc
		metadata["description"] = desc
	}

	var updated repository.Node
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		node, err := txRepos.Node.Update(l.ctx, req.NodeID, input)
		if err != nil {
			return err
		}
		updated = node

		if len(metadata) == 0 {
			return nil
		}
		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.node.update",
			ResourceType: "node",
			ResourceID:   fmt.Sprintf("%d", updated.ID),
			Metadata:     metadata,
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminNodeResponse{
		Node: mapNodeSummary(updated),
	}, nil
}
