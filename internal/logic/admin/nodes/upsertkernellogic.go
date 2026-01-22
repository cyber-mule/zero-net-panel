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

// UpsertKernelLogic handles kernel endpoint configuration.
type UpsertKernelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpsertKernelLogic constructs UpsertKernelLogic.
func NewUpsertKernelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertKernelLogic {
	return &UpsertKernelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Upsert configures node kernel endpoint for a protocol.
func (l *UpsertKernelLogic) Upsert(req *types.AdminUpsertNodeKernelRequest) (*types.AdminNodeKernelUpsertResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}

	protocol := strings.TrimSpace(req.Protocol)
	endpoint := strings.TrimSpace(req.Endpoint)
	if protocol == "" || endpoint == "" {
		return nil, repository.ErrInvalidArgument
	}

	var revisionPtr *string
	if strings.TrimSpace(req.Revision) != "" {
		rev := strings.TrimSpace(req.Revision)
		revisionPtr = &rev
	}

	var statusPtr *int
	if req.Status != nil {
		statusCode, err := normalizeKernelStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		statusPtr = &statusCode
	}

	var lastSynced *time.Time
	if req.LastSyncedAt != nil {
		if *req.LastSyncedAt <= 0 {
			return nil, repository.ErrInvalidArgument
		}
		ts := time.Unix(*req.LastSyncedAt, 0).UTC()
		lastSynced = &ts
	}

	input := repository.UpsertNodeKernelInput{
		Protocol:     protocol,
		Endpoint:     endpoint,
		Revision:     revisionPtr,
		Status:       statusPtr,
		Config:       req.Config,
		LastSyncedAt: lastSynced,
	}

	var kernel repository.NodeKernel
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		stored, err := txRepos.Node.UpsertKernel(l.ctx, req.NodeID, input)
		if err != nil {
			return err
		}
		kernel = stored

		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.node.kernel.upsert",
			ResourceType: "node",
			ResourceID:   fmt.Sprintf("%d", req.NodeID),
			Metadata: map[string]any{
				"protocol": kernel.Protocol,
				"endpoint": kernel.Endpoint,
				"status":   kernel.Status,
			},
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminNodeKernelUpsertResponse{
		NodeID: req.NodeID,
		Kernel: mapKernelSummary(kernel),
	}, nil
}
