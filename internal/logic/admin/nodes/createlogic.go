package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/nodecfg"
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
	endpoint := strings.TrimSpace(req.ControlEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("%w: node control endpoint required", repository.ErrInvalidArgument)
	}
	accessKey := strings.TrimSpace(req.ControlAccessKey)
	if accessKey == "" {
		accessKey = strings.TrimSpace(req.AK)
	}
	secretKey := strings.TrimSpace(req.ControlSecretKey)
	if secretKey == "" {
		secretKey = strings.TrimSpace(req.SK)
	}
	statusSyncEnabled := true
	if req.StatusSyncEnabled != nil {
		statusSyncEnabled = *req.StatusSyncEnabled
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
	if req.CapacityMbps < 0 {
		return nil, repository.ErrInvalidArgument
	}

	kernelDefaultProtocol := strings.TrimSpace(req.KernelDefaultProtocol)
	if kernelDefaultProtocol == "" {
		kernelDefaultProtocol = nodecfg.DefaultKernelProtocol
	}
	if kernelDefaultProtocol == "" {
		return nil, repository.ErrInvalidArgument
	}
	kernelHTTPTimeoutSeconds := nodecfg.DefaultKernelHTTPTimeoutSeconds
	if req.KernelHTTPTimeoutSeconds != nil {
		kernelHTTPTimeoutSeconds = *req.KernelHTTPTimeoutSeconds
	}
	if kernelHTTPTimeoutSeconds <= 0 {
		return nil, repository.ErrInvalidArgument
	}
	kernelStatusPollIntervalSeconds := nodecfg.DefaultKernelStatusPollIntervalSeconds
	if req.KernelStatusPollIntervalSeconds != nil {
		kernelStatusPollIntervalSeconds = *req.KernelStatusPollIntervalSeconds
	}
	if kernelStatusPollIntervalSeconds < 0 {
		return nil, repository.ErrInvalidArgument
	}
	kernelStatusPollBackoffEnabled := nodecfg.DefaultKernelStatusPollBackoffEnabled
	if req.KernelStatusPollBackoffEnabled != nil {
		kernelStatusPollBackoffEnabled = *req.KernelStatusPollBackoffEnabled
	}
	kernelStatusPollBackoffMaxIntervalSeconds := nodecfg.DefaultKernelStatusPollBackoffMaxIntervalSeconds
	if req.KernelStatusPollBackoffMaxIntervalSeconds != nil {
		kernelStatusPollBackoffMaxIntervalSeconds = *req.KernelStatusPollBackoffMaxIntervalSeconds
	}
	if kernelStatusPollBackoffMaxIntervalSeconds < 0 {
		return nil, repository.ErrInvalidArgument
	}
	kernelStatusPollBackoffMultiplier := float64(nodecfg.DefaultKernelStatusPollBackoffMultiplier)
	if req.KernelStatusPollBackoffMultiplier != nil {
		kernelStatusPollBackoffMultiplier = *req.KernelStatusPollBackoffMultiplier
	}
	if kernelStatusPollBackoffMultiplier <= 1 {
		return nil, repository.ErrInvalidArgument
	}
	kernelStatusPollBackoffJitter := nodecfg.DefaultKernelStatusPollBackoffJitter
	if req.KernelStatusPollBackoffJitter != nil {
		kernelStatusPollBackoffJitter = *req.KernelStatusPollBackoffJitter
	}
	if kernelStatusPollBackoffJitter < 0 || kernelStatusPollBackoffJitter > 1 {
		return nil, repository.ErrInvalidArgument
	}
	kernelOfflineProbeMaxIntervalSeconds := nodecfg.DefaultKernelOfflineProbeMaxIntervalSeconds
	if req.KernelOfflineProbeMaxIntervalSeconds != nil {
		kernelOfflineProbeMaxIntervalSeconds = *req.KernelOfflineProbeMaxIntervalSeconds
	}
	if kernelOfflineProbeMaxIntervalSeconds < 0 {
		return nil, repository.ErrInvalidArgument
	}
	if kernelStatusPollBackoffEnabled && kernelStatusPollIntervalSeconds <= 0 {
		return nil, repository.ErrInvalidArgument
	}
	if kernelStatusPollBackoffEnabled && kernelStatusPollBackoffMaxIntervalSeconds > 0 &&
		kernelStatusPollBackoffMaxIntervalSeconds < kernelStatusPollIntervalSeconds {
		return nil, repository.ErrInvalidArgument
	}

	now := time.Now().UTC()
	node := repository.Node{
		Name:                            name,
		Region:                          strings.TrimSpace(req.Region),
		Country:                         strings.TrimSpace(req.Country),
		ISP:                             strings.TrimSpace(req.ISP),
		Status:                          status,
		Tags:                            tags,
		CapacityMbps:                    req.CapacityMbps,
		Description:                     strings.TrimSpace(req.Description),
		AccessAddress:                   strings.TrimSpace(req.AccessAddress),
		ControlEndpoint:                 endpoint,
		ControlAccessKey:                accessKey,
		ControlSecretKey:                secretKey,
		ControlToken:                    strings.TrimSpace(req.ControlToken),
		KernelDefaultProtocol:           kernelDefaultProtocol,
		KernelHTTPTimeoutSeconds:        kernelHTTPTimeoutSeconds,
		KernelStatusPollIntervalSeconds: kernelStatusPollIntervalSeconds,
		KernelStatusPollBackoffEnabled:  kernelStatusPollBackoffEnabled,
		KernelStatusPollBackoffMaxIntervalSeconds: kernelStatusPollBackoffMaxIntervalSeconds,
		KernelStatusPollBackoffMultiplier:         kernelStatusPollBackoffMultiplier,
		KernelStatusPollBackoffJitter:             kernelStatusPollBackoffJitter,
		KernelOfflineProbeMaxIntervalSeconds:      kernelOfflineProbeMaxIntervalSeconds,
		StatusSyncEnabled:                         statusSyncEnabled,
		CreatedAt:                                 now,
		UpdatedAt:                                 now,
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
				"name":   created.Name,
				"status": created.Status,
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
